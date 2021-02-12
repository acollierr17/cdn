package main

import (
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/gorilla/mux"
	"log"
	"math/rand"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"strings"
	"time"
)

func errorHandler(w http.ResponseWriter, _ *http.Request, statusCode int, err string) {
	w.WriteHeader(statusCode)

	jsonObj := ImageResponse{
		Url:     fmt.Sprintf("An error occurred: %v", err),
		Success: false,
	}

	b, jsonError := json.Marshal(jsonObj)
	if jsonError != nil {
		log.Fatal(jsonError)
	}

	_, writeError := w.Write(b)
	if writeError != nil {
		log.Fatal(jsonError)
	}
}

func handleRequests() {
	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.HandleFunc("/{id}", getImageRoute).Methods("GET")
	myRouter.HandleFunc("/oembed/{id}", getOGEmbedRoute).Methods("GET")
	myRouter.HandleFunc("/upload", uploadImageRoute).Methods("POST")
	myRouter.HandleFunc("/", homePageRoute).Methods("GET")
	log.Fatal(http.ListenAndServe(":3000", myRouter))
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

func uploadFileToS3(s *session.Session, file multipart.File, fileHeader *multipart.FileHeader) (string, error) {
	size := fileHeader.Size
	buffer := make([]byte, size)
	_, err := file.Read(buffer)
	if err != nil {
		return "", err
	}

	rand.Seed(time.Now().UnixNano())
	tempFileName := randSeq(7) + filepath.Ext(fileHeader.Filename)

	s, err = session.NewSession(s3Config)
	if err != nil {
		return "", err
	}

	object := s3.PutObjectInput{
		Bucket:               aws.String(cdnConfig.SpacesConfig.SpacesName),
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
