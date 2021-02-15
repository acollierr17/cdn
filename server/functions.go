package main

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"log"
	"math/rand"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"strings"
	"time"
)

func handleRequests() {
	app := fiber.New()

	app.Use(cors.New(cors.Config{
		AllowOrigins: "http://localhost:8080, https://acollier.dev, https://cdn.acollier.dev, https://acolliercdn.ngrok.io",
		AllowHeaders: "Origin, Content-Type, Accept, Access-Token, User-Agent",
	}))

	app.Static("/admin", "dist")

	api := app.Group("/api").Use(accessTokenMiddleware)
	api.Post("/upload", uploadImageRoute)
	api.Delete("/delete/:id", deleteImageRoute)
	api.Post("/token", generateAccessTokenRoute)

	app.Get("/:id", getImageRoute)
	app.Get("/oembed/:id", getOGEmbedRoute)
	app.Get("/", homePageRoute)

	admin := app.Group("/admin")
	admin.Get("/admin/*", func(ctx *fiber.Ctx) error {
		return ctx.SendFile("./dist/index.html")
	})

	log.Fatal(app.Listen(":3000"))
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

func uploadFileToS3(s *session.Session, fileHeader *multipart.FileHeader) (string, error) {
	size := fileHeader.Size
	buffer := make([]byte, size)
	file, headerOpenErr := fileHeader.Open()
	if headerOpenErr != nil {
		return "", headerOpenErr
	}
	_, fileReadErr := file.Read(buffer)
	if fileReadErr != nil {
		return "", fileReadErr
	}

	file.Close()

	rand.Seed(time.Now().UnixNano())
	tempFileName := randSeq(7) + filepath.Ext(fileHeader.Filename)

	object := s3.PutObjectInput{
		Bucket:               aws.String(cdnConfig.SpacesConfig.SpacesName),
		Key:                  aws.String(tempFileName),
		ACL:                  aws.String("public-read"),
		Body:                 strings.NewReader(string(buffer)),
		ContentLength:        aws.Int64(size),
		ContentType:          aws.String(http.DetectContentType(buffer)),
		ServerSideEncryption: aws.String("AES256"),
	}

	_, putObjErr := s3.New(s).PutObject(&object)
	if putObjErr != nil {
		return "", putObjErr
	}

	return tempFileName, nil
}

func deleteFileFromS3(s *session.Session, key string) (bool, error) {
	object := &s3.DeleteObjectInput{
		Bucket: aws.String(cdnConfig.SpacesConfig.SpacesName),
		Key:    aws.String(key),
	}

	_, err := s3.New(s).DeleteObject(object)
	if err != nil {
		return false, err
	}

	return true, err
}

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randSeq(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func generateToken() string {
	b := make([]byte, 32)
	rand.Read(b)
	return fmt.Sprintf("%x", b)
}
