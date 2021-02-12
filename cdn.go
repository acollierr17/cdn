package main

import (
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"log"
	"math/rand"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

type Embed struct {
	Type         string `json:"type"`
	AuthorName   string `json:"author_name"`
	ProviderName string `json:"provider_name"`
}

type ImageResponse struct {
	Url     string `json:"url"`
	Success bool   `json:"success"`
}

var key = getConfigVar("SPACES_ACCESS_KEY")
var secret = getConfigVar("SPACES_SECRET_KEY")
var endpoint = getConfigVar("SPACES_ENDPOINT")
var spaceURL = getConfigVar("SPACES_URL")
var spaceName = getConfigVar("SPACES_NAME")
var region = getConfigVar("SPACES_REGION")
var cdnEndpoint = getConfigVar("CDN_ENDPOINT")

var s3Config = &aws.Config{
	Credentials: credentials.NewStaticCredentials(key, secret, ""),
	Endpoint:    aws.String(endpoint),
	Region:      aws.String(region),
}

func getConfigVar(key string) string {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	return os.Getenv(key)
}

func errorHandler(w http.ResponseWriter, _ *http.Request, statusCode int) {
	w.WriteHeader(statusCode)
	_, err := fmt.Fprintf(w, "There was an error finding the image...")
	if err != nil {
		log.Fatal(err)
	}
}

func roundFloat64(num float64) string {
	return fmt.Sprintf("%.2f", num)
}

func getFileSize(size int64) string {
	if size < 1024 {
		return fmt.Sprintf("%vB", size)
	} else if size < 1048576 {
		num := float64(size / 1024)
		return fmt.Sprintf("%vKiB", roundFloat64(num))
	} else if size < 1073741824 {
		num := float64(size / 1048576)
		return fmt.Sprintf("%vMiB", roundFloat64(num))
	} else {
		num := float64(size / 1073741824)
		return fmt.Sprintf("%vGiB", roundFloat64(num))
	}
}

func homePage(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "https://acollier.dev", http.StatusPermanentRedirect)
}

func getImage(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["id"]

	fmt.Println("Endpoint Hit: getImage")

	imageURL := fmt.Sprintf("%s/%s", spaceURL, key)
	oembedURL := fmt.Sprintf("%s/oembed/%s", cdnEndpoint, key)
	res, err := http.Head(imageURL)
	if err != nil {
		log.Fatal(err)
	}

	if res.StatusCode == http.StatusOK {
		if r.Header.Get("User-Agent") == "Mozilla/5.0 (compatible; Discordbot/2.0; +https://discordapp.com)" {
			_, err := fmt.Fprintf(
				w,
				`<!DOCTYPE html>
					<html>
						<head>
							<meta name="theme-color" content="#dd9323">
							<meta property="og:title" content="%v">
							<meta content="%s" property="og:image">
							<link type="application/json+oembed" href="%s" />
						</head>
					</html>
				`,
				key, imageURL, oembedURL)
			if err != nil {
				log.Fatal(err)
			}
		} else {
			http.Redirect(w, r, imageURL, http.StatusMovedPermanently)
		}
	} else {
		errorHandler(w, r, res.StatusCode)
	}
}

func getOGEmbed(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["id"]

	fmt.Println("Endpoint Hit: getOGEmbed")

	imageURL := fmt.Sprintf("%s/%s", spaceURL, key)

	res, err := http.Head(imageURL)
	if err != nil {
		log.Fatal(err)
	}

	if res.StatusCode == http.StatusOK {
		if r.Header.Get("User-Agent") == "Mozilla/5.0 (compatible; Discordbot/2.0; +https://discordapp.com)" {
			contentLengthHeader := res.Header.Values("content-length")[0]
			contentTypeHeader := res.Header.Values("content-type")[0]
			contentLength, err := strconv.ParseInt(contentLengthHeader, 0, 64)
			if err != nil {
				log.Fatal(err)
			}

			var objType string
			objAuthor := fmt.Sprintf("%v | %v", getFileSize(contentLength), contentTypeHeader)
			objProvider := res.Header.Values("last-modified")[0]

			if strings.HasPrefix(contentTypeHeader, "image") {
				objType = "photo"
			} else if strings.HasPrefix(contentTypeHeader, "video") {
				objType = "video"
			} else {
				objType = "link"
			}

			jsonObj := Embed{objType, objAuthor, objProvider}
			b, err := json.Marshal(jsonObj)
			if err != nil {
				log.Fatal(err)
			}

			_, err = w.Write(b)
			if err != nil {
				log.Fatal(err)
			}
		} else {
			http.Redirect(w, r, imageURL, http.StatusMovedPermanently)
		}
	} else {
		errorHandler(w, r, res.StatusCode)
	}
}

func uploadImage(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint Hit: uploadImage")
	if r.Header.Get("Access-Token") == "" {
		fmt.Fprintf(w, "Please provide an access token!")
		return
	}

	if r.Header.Get("Access-Token") != getConfigVar("ACCESS_TOKEN") {
		fmt.Fprintf(w, "Invalid access token provided!")
		return
	}

	file, fileHeader, err := r.FormFile("image")
	if err != nil {
		log.Println(err)
		fmt.Fprintf(w, "Could not get uploaded file")
		return
	}

	defer file.Close()

	s, err := session.NewSession(s3Config)
	if err != nil {
		log.Fatal(err)
	}

	fileName, err := uploadFileToS3(s, file, fileHeader)
	if err != nil {
		log.Fatal(err)
	}

	jsonObj := ImageResponse{fmt.Sprintf("%v/%v", cdnEndpoint, fileName), true}
	b, err := json.Marshal(jsonObj)
	if err != nil {
		log.Fatal(err)
	}

	_, err = w.Write(b)
	if err != nil {
		log.Fatal(err)
	}
}

func uploadFileToS3(s *session.Session, file multipart.File, fileHeader *multipart.FileHeader) (string, error) {
	size := fileHeader.Size
	buffer := make([]byte, size)
	_, err := file.Read(buffer)
	if err != nil {
		log.Fatal(err)
	}

	rand.Seed(time.Now().UnixNano())
	tempFileName := randSeq(7) + filepath.Ext(fileHeader.Filename)

	s, err = session.NewSession(s3Config)
	if err != nil {
		log.Fatal(err)
	}

	object := s3.PutObjectInput{
		Bucket:               aws.String(spaceName),
		Key:                  aws.String(tempFileName),
		ACL:                  aws.String("public-read"),
		Body:                 strings.NewReader(string(buffer)),
		ContentLength:        aws.Int64(size),
		ContentType:          aws.String(http.DetectContentType(buffer)),
		ServerSideEncryption: aws.String("AES256"),
	}

	_, err = s3.New(s).PutObject(&object)
	if err != nil {
		return "", err
	}

	return tempFileName, err
}

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randSeq(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func handleRequests() {
	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.HandleFunc("/{id}", getImage).Methods("GET")
	myRouter.HandleFunc("/oembed/{id}", getOGEmbed).Methods("GET")
	myRouter.HandleFunc("/upload", uploadImage).Methods("POST")
	myRouter.HandleFunc("/", homePage).Methods("GET")
	log.Fatal(http.ListenAndServe(":3000", myRouter))
}

func main() {
	fmt.Println("CDN v1.0 - Copyright (c) 2021 Anthony Collier")
	handleRequests()
}
