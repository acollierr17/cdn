package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/gofiber/fiber/v2"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

func homePageRoute(ctx *fiber.Ctx) error {
	return ctx.Redirect("https://acollier.dev", fiber.StatusMovedPermanently)
}

func getImageRoute(ctx *fiber.Ctx) error {
	key := ctx.Params("id")

	fmt.Println("Endpoint Hit: getImage")

	queries := new(ImageResponseQuery)

	if queryErr := ctx.QueryParser(queries); queryErr != nil {
		return fiber.NewError(fiber.StatusInternalServerError, queryErr.Error())
	}

	imageURL := fmt.Sprintf("%s/%s", cdnConfig.SpacesConfig.SpacesUrl, key)
	oembedURL := fmt.Sprintf("%s/oembed/%s", cdnConfig.CdnEndpoint, key)
	res, err := http.Head(imageURL)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	if res.StatusCode != http.StatusOK {
		return fiber.NewError(res.StatusCode, "An error occurred redirecting to the image")
	}

	if queries.Download == "true" {
		ctx.Set("Content-Type", res.Header.Values("content-type")[0])
		ctx.Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%v"`, key))
		ctx.Set("Content-Length", fmt.Sprintf("%v", res.ContentLength))

		resp, _ := http.Get(imageURL)

		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}

		return ctx.Send(body)
	}


	if ctx.Get("User-Agent") == "Mozilla/5.0 (compatible; Discordbot/2.0; +https://discordapp.com)" {
		ctx.Type("html")
		return ctx.Send([]byte(fmt.Sprintf(
			`<!DOCTYPE html>
			<html>
				<head>
					<meta name="theme-color" content="#dd9323">
					<meta property="og:title" content="%v">
					<meta content="%v" property="og:image">
					<link type="application/json+oembed" href="%v" />
				</head>
			</html>`,
			key, imageURL, oembedURL)),
		)
	} else {
		return ctx.Redirect(imageURL, fiber.StatusMovedPermanently)
	}
}

func getOGEmbedRoute(ctx *fiber.Ctx) error {
	key := ctx.Params("id")

	fmt.Println("Endpoint Hit: getOGEmbed")

	imageURL := fmt.Sprintf("%s/%s", cdnConfig.SpacesConfig.SpacesUrl, key)

	res, headErr := http.Head(imageURL)
	if headErr != nil {
		return fiber.NewError(fiber.StatusInternalServerError, headErr.Error())
	}

	if res.StatusCode != http.StatusOK {
		return fiber.NewError(res.StatusCode, "An error occurred redirecting to the image.")
	}

	if ctx.Get("User-Agent") == "Mozilla/5.0 (compatible; Discordbot/2.0; +https://discordapp.com)" {
		ctx.Type("json", "utf-8")
		contentLengthHeader := res.Header.Values("content-length")[0]
		contentTypeHeader := res.Header.Values("content-type")[0]
		contentLength, parseErr := strconv.ParseInt(contentLengthHeader, 0, 64)
		if parseErr != nil {
			return fiber.NewError(fiber.StatusInternalServerError, parseErr.Error())
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
			return fiber.NewError(fiber.StatusInternalServerError, jsonErr.Error())
		}

		return ctx.Send(b)
	} else {
		return ctx.Redirect(imageURL, fiber.StatusMovedPermanently)
	}
}

func uploadImageRoute(ctx *fiber.Ctx) error {
	fmt.Println("Endpoint Hit: uploadImage")

	fileHeader, err := ctx.FormFile("image")
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Could not get the uploaded file!")
	}

	s, sessionErr := session.NewSession(s3Config)
	if sessionErr != nil {
		return fiber.NewError(fiber.StatusInternalServerError, sessionErr.Error())
	}

	fileName, s3UploadErr := uploadFileToS3(s, fileHeader)
	if s3UploadErr != nil {
		return fiber.NewError(fiber.StatusInternalServerError, s3UploadErr.Error())
	}

	jsonObj := ImageResponse{
		Url:     fmt.Sprintf("%v/%v", cdnConfig.CdnEndpoint, fileName),
		Success: true,
	}
	b, jsonErr := json.Marshal(jsonObj)
	if jsonErr != nil {
		return fiber.NewError(fiber.StatusInternalServerError, jsonErr.Error())
	}

	return ctx.Send(b)
}

func deleteImageRoute(ctx *fiber.Ctx) error {
	key := ctx.Params("id")

	fmt.Println("Endpoint Hit: deleteImage")

	imageURL := fmt.Sprintf("%s/%s", cdnConfig.SpacesConfig.SpacesUrl, key)
	res, err := http.Get(imageURL)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	if res.StatusCode != http.StatusOK {
		return fiber.NewError(fiber.StatusInternalServerError, "An error occurred redirecting to the image.")
	}

	s, sessionErr := session.NewSession(s3Config)
	if sessionErr != nil {
		return fiber.NewError(fiber.StatusInternalServerError, sessionErr.Error())
	}

	fileDeleted, s3DeleteErr := deleteFileFromS3(s, key)
	if s3DeleteErr != nil {
		return fiber.NewError(fiber.StatusInternalServerError, s3DeleteErr.Error())
	}

	jsonObj := ImageDeletedRespone{
		ImageName: key,
		Deleted:   fileDeleted,
	}
	b, jsonErr := json.Marshal(jsonObj)
	if jsonErr != nil {
		return fiber.NewError(fiber.StatusInternalServerError, jsonErr.Error())
	}

	return ctx.Send(b)
}

func generateAccessTokenRoute(ctx *fiber.Ctx) error {
	t := new(AuthTokenRequest)

	fmt.Println("Endpoint Hit: token")

	if err := ctx.BodyParser(t); err != nil {
		return fiber.NewError(fiber.StatusUnprocessableEntity, err.Error())
	}

	token := t.Token
	firebaseCtx := context.Background()

	verified, verifiedErr := firebaseAuth.VerifyIDToken(firebaseCtx, token)
	if verifiedErr != nil {
		return fiber.NewError(fiber.StatusUnauthorized, verifiedErr.Error())
	}

	verifiedUID := verified.UID
	accessToken := generateToken()

	_, updateErr := firebaseFirestore.Collection("tokens").Doc(verifiedUID).Set(firebaseCtx, map[string]interface{}{
		"uid": verifiedUID,
		"token": accessToken,
	})
	if updateErr != nil {
		return fiber.NewError(fiber.StatusInternalServerError, updateErr.Error())
	}

	return ctx.JSON(fiber.Map{
		"success": true,
	})
}

func getImagesRoute(ctx *fiber.Ctx) error  {
	fmt.Println("Endpoint Hit: getImages")

	s, sessionErr := session.NewSession(s3Config)
	if sessionErr != nil {
		return fiber.NewError(fiber.StatusInternalServerError, sessionErr.Error())
	}

	objects, objectsErr := getFilesFromS3(s)
	if objectsErr != nil {
		return fiber.NewError(fiber.StatusInternalServerError, objectsErr.Error())
	}

	data := ImagesResults{
		Images: objects,
		Length: len(objects),
	}

	return ctx.JSON(data)
}
