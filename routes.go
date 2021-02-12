package main

import (
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
	"strings"
)

func homePageRoute(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "https://acollier.dev", http.StatusPermanentRedirect)
}

func getImageRoute(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["id"]

	fmt.Println("Endpoint Hit: getImage")

	imageURL := fmt.Sprintf("%s/%s", cdnConfig.SpacesConfig.SpacesUrl, key)
	oembedURL := fmt.Sprintf("%s/oembed/%s", cdnConfig.CdnEndpoint, key)
	res, err := http.Head(imageURL)
	if err != nil {
		errorHandler(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	if res.StatusCode != http.StatusOK {
		errorHandler(w, r, res.StatusCode, "An error occurred redirecting to the image.")
		return
	}

	if r.Header.Get("User-Agent") == "Mozilla/5.0 (compatible; Discordbot/2.0; +https://discordapp.com)" {
		_, fmtErr := fmt.Fprintf(
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
		if fmtErr != nil {
			errorHandler(w, r, http.StatusInternalServerError, fmtErr.Error())
			return
		}
	} else {
		http.Redirect(w, r, imageURL, http.StatusMovedPermanently)
	}
}

func getOGEmbedRoute(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["id"]

	fmt.Println("Endpoint Hit: getOGEmbed")

	imageURL := fmt.Sprintf("%s/%s", cdnConfig.SpacesConfig.SpacesUrl, key)

	res, headErr := http.Head(imageURL)
	if headErr != nil {
		errorHandler(w, r, http.StatusInternalServerError, headErr.Error())
		return
	}

	if res.StatusCode != http.StatusOK {
		errorHandler(w, r, res.StatusCode, "An error occurred redirecting to the image.")
		return
	}

	if r.Header.Get("User-Agent") == "Mozilla/5.0 (compatible; Discordbot/2.0; +https://discordapp.com)" {
		contentLengthHeader := res.Header.Values("content-length")[0]
		contentTypeHeader := res.Header.Values("content-type")[0]
		contentLength, parseErr := strconv.ParseInt(contentLengthHeader, 0, 64)
		if parseErr != nil {
			errorHandler(w, r, http.StatusInternalServerError, parseErr.Error())
			return
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

		jsonObj := Embed{
			Type:         objType,
			AuthorName:   objAuthor,
			ProviderName: objProvider,
		}
		b, jsonErr := json.Marshal(jsonObj)
		if jsonErr != nil {
			errorHandler(w, r, http.StatusInternalServerError, jsonErr.Error())
			return
		}

		_, writeErr := w.Write(b)
		if writeErr != nil {
			errorHandler(w, r, http.StatusInternalServerError, writeErr.Error())
			return
		}
	} else {
		http.Redirect(w, r, imageURL, http.StatusMovedPermanently)
	}
}

func uploadImageRoute(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint Hit: uploadImage")

	if r.Header.Get("Access-Token") == "" {
		errorHandler(w, r, http.StatusUnauthorized, "Invalid access token provided!")
		return
	}

	if r.Header.Get("Access-Token") != cdnConfig.AccessToken {
		errorHandler(w, r, http.StatusUnauthorized, "Invalid access token provided!")
		return
	}

	file, fileHeader, err := r.FormFile("image")
	if err != nil {
		errorHandler(w, r, http.StatusUnprocessableEntity, "Could not get the uploaded file!")
		return
	}

	defer file.Close()

	s, sessionErr := session.NewSession(s3Config)
	if sessionErr != nil {
		errorHandler(w, r, http.StatusInternalServerError, sessionErr.Error())
		return
	}

	fileName, s3UploadErr := uploadFileToS3(s, file, fileHeader)
	if s3UploadErr != nil {
		errorHandler(w, r, http.StatusInternalServerError, s3UploadErr.Error())
		return
	}

	jsonObj := ImageResponse{
		Url:     fmt.Sprintf("%v/%v", cdnConfig.CdnEndpoint, fileName),
		Success: true,
	}
	b, jsonErr := json.Marshal(jsonObj)
	if jsonErr != nil {
		errorHandler(w, r, http.StatusInternalServerError, jsonErr.Error())
		return
	}

	_, writeErr := w.Write(b)
	if writeErr != nil {
		errorHandler(w, r, http.StatusInternalServerError, writeErr.Error())
		return
	}
}
