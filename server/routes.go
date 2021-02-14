package main

import (
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/gofiber/fiber"
	"net/http"
	"strconv"
	"strings"
)

func homePageRoute(ctx *fiber.Ctx) {
	ctx.Redirect("https://acollier.dev", fiber.StatusMovedPermanently)
}

func getImageRoute(ctx *fiber.Ctx) {
	key := ctx.Params("id")

	fmt.Println("Endpoint Hit: getImage")

	imageURL := fmt.Sprintf("%s/%s", cdnConfig.SpacesConfig.SpacesUrl, key)
	oembedURL := fmt.Sprintf("%s/oembed/%s", cdnConfig.CdnEndpoint, key)
	res, err := http.Head(imageURL)
	if err != nil {
		errorHandler(ctx, fiber.StatusInternalServerError, err.Error())
		return
	}

	if res.StatusCode != http.StatusOK {
		errorHandler(ctx, res.StatusCode, "An error occurred redirecting to the image.")
		return
	}

	if ctx.Get("User-Agent") == "Mozilla/5.0 (compatible; Discordbot/2.0; +https://discordapp.com)" {
		ctx.Type("html")
		ctx.Send(fmt.Sprintf(
			`<!DOCTYPE html>
			<html>
				<head>
					<meta name="theme-color" content="#dd9323">
					<meta property="og:title" content="%v">
					<meta content="%v" property="og:image">
					<link type="application/json+oembed" href="%v" />
				</head>
			</html>`,
			key, imageURL, oembedURL),
		)
	} else {
		ctx.Redirect(imageURL, fiber.StatusMovedPermanently)
	}
}

func getOGEmbedRoute(ctx *fiber.Ctx) {
	key := ctx.Params("id")

	fmt.Println("Endpoint Hit: getOGEmbed")

	imageURL := fmt.Sprintf("%s/%s", cdnConfig.SpacesConfig.SpacesUrl, key)

	res, headErr := http.Head(imageURL)
	if headErr != nil {
		errorHandler(ctx, fiber.StatusInternalServerError, headErr.Error())
		return
	}

	if res.StatusCode != http.StatusOK {
		errorHandler(ctx, res.StatusCode, "An error occurred redirecting to the image.")
		return
	}

	if ctx.Get("User-Agent") == "Mozilla/5.0 (compatible; Discordbot/2.0; +https://discordapp.com)" {
		ctx.Type("json", "utf-8")
		contentLengthHeader := res.Header.Values("content-length")[0]
		contentTypeHeader := res.Header.Values("content-type")[0]
		contentLength, parseErr := strconv.ParseInt(contentLengthHeader, 0, 64)
		if parseErr != nil {
			errorHandler(ctx, fiber.StatusInternalServerError, parseErr.Error())
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
			errorHandler(ctx, fiber.StatusInternalServerError, jsonErr.Error())
			return
		}

		ctx.SendBytes(b)
	} else {
		ctx.Redirect(imageURL, fiber.StatusMovedPermanently)
	}
}

func uploadImageRoute(ctx *fiber.Ctx) {
	fmt.Println("Endpoint Hit: uploadImage")

	fileHeader, err := ctx.FormFile("image")
	if err != nil {
		errorHandler(ctx, fiber.StatusUnprocessableEntity, "Could not get the uploaded file!")
		return
	}

	s, sessionErr := session.NewSession(s3Config)
	if sessionErr != nil {
		errorHandler(ctx, fiber.StatusInternalServerError, sessionErr.Error())
		return
	}

	fileName, s3UploadErr := uploadFileToS3(s, fileHeader)
	if s3UploadErr != nil {
		errorHandler(ctx, fiber.StatusInternalServerError, s3UploadErr.Error())
		return
	}

	jsonObj := ImageResponse{
		Url:     fmt.Sprintf("%v/%v", cdnConfig.CdnEndpoint, fileName),
		Success: true,
	}
	b, jsonErr := json.Marshal(jsonObj)
	if jsonErr != nil {
		errorHandler(ctx, fiber.StatusInternalServerError, jsonErr.Error())
		return
	}

	ctx.Write(b)
}

func deleteImageRoute(ctx *fiber.Ctx) {
	key := ctx.Params("id")

	fmt.Println("Endpoint Hit: deleteImage")

	imageURL := fmt.Sprintf("%s/%s", cdnConfig.SpacesConfig.SpacesUrl, key)
	res, err := http.Get(imageURL)
	if err != nil {
		errorHandler(ctx, fiber.StatusInternalServerError, err.Error())
		return
	}

	if res.StatusCode != http.StatusOK {
		errorHandler(ctx, res.StatusCode, "An error occurred redirecting to the image.")
		return
	}

	s, sessionErr := session.NewSession(s3Config)
	if sessionErr != nil {
		errorHandler(ctx, fiber.StatusInternalServerError, sessionErr.Error())
		return
	}

	fileDeleted, s3DeleteErr := deleteFileFromS3(s, key)
	if s3DeleteErr != nil {
		errorHandler(ctx, fiber.StatusInternalServerError, s3DeleteErr.Error())
		return
	}

	jsonObj := ImageDeletedRespone{
		ImageName: key,
		Deleted:   fileDeleted,
	}
	b, jsonErr := json.Marshal(jsonObj)
	if jsonErr != nil {
		errorHandler(ctx, fiber.StatusInternalServerError, jsonErr.Error())
		return
	}

	ctx.Write(b)
}
