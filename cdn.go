package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)

type Embed struct {
	Type         string `json:"type"`
	AuthorName   string `json:"author_name"`
	ProviderName string `json:"provider_name"`
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

func homePage(w http.ResponseWriter, _ *http.Request) {
	fmt.Println("Endpoint Hit: homePage")
	_, err := fmt.Fprintf(w, "Welcome to the HomePage!")
	if err != nil {
		log.Fatal(err)
	}
}

func getImage(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["id"]

	fmt.Println("Endpoint Hit: getImage")

	s3Endpoint := getConfigVar("SPACES_ENDPOINT")
	cdnEndpoint := getConfigVar("CDN_ENDPOINT")
	imageURL := fmt.Sprintf("%s/%s", s3Endpoint, key)
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

	s3Endpoint := getConfigVar("SPACES_ENDPOINT")
	imageURL := fmt.Sprintf("%s/%s", s3Endpoint, key)

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

//func uploadImage(res *http.Request)  {
//	s3Endpoint := getConfigVar("SPACES_ENDPOINT")
//	s3AccessKey := getConfigVar("SPACES_ACCESS_KEY")
//	s3SecretKey := getConfigVar("SPACES_SECRET_KEY")
//
//
//}

func handleRequests() {
	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.HandleFunc("/", homePage)
	myRouter.HandleFunc("/{id}", getImage)
	myRouter.HandleFunc("/oembed/{id}", getOGEmbed)
	log.Fatal(http.ListenAndServe(":3000", myRouter))
}

func main() {
	fmt.Println("CDN v1.0 - Copyright (c) 2021 Anthony")
	handleRequests()
}
