package main

import (
	"cloud.google.com/go/firestore"
	"context"
	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/joho/godotenv"
	"google.golang.org/api/option"
	"log"
	"os"
)

var cdnConfig *Config
var s3Config *aws.Config
var firebaseApp *firebase.App
var firebaseAuth *auth.Client
var firebaseFirestore *firestore.Client

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	cdnConfig = &Config{
		SpacesConfig: SpacesConfig{
			SpacesAccessKey: 	os.Getenv("SPACES_ACCESS_KEY"),
			SpacesSecretKey: 	os.Getenv("SPACES_SECRET_KEY"),
			SpacesEndpoint:  	os.Getenv("SPACES_ENDPOINT"),
			SpacesUrl:       	os.Getenv("SPACES_URL"),
			SpacesCdn: 			os.Getenv("SPACES_CDN_URL"),
			SpacesName:      	os.Getenv("SPACES_NAME"),
			SpacesRegion:    	os.Getenv("SPACES_REGION"),
		},
		CdnEndpoint: os.Getenv("CDN_ENDPOINT"),
		AccessToken: os.Getenv("ACCESS_TOKEN"),
	}

	s3Config = &aws.Config{
		Credentials: credentials.NewStaticCredentials(cdnConfig.SpacesConfig.SpacesAccessKey, cdnConfig.SpacesConfig.SpacesSecretKey, ""),
		Endpoint:    aws.String(cdnConfig.SpacesConfig.SpacesEndpoint),
		Region:      aws.String(cdnConfig.SpacesConfig.SpacesRegion),
	}

	opt := option.WithCredentialsFile("./service-account.json")
	ctx := context.Background()

	app, newAppErr := firebase.NewApp(ctx, nil, opt)
	if err != newAppErr {
		log.Fatal(newAppErr)
		return
	}

	firebaseApp = app

	fbAuth, newAuthErr := app.Auth(ctx)
	if err != newAuthErr {
		log.Fatal(newAuthErr)
	}

	firebaseAuth = fbAuth

	fbFirestore, newFirestoreErr := app.Firestore(ctx)
	if newFirestoreErr != nil {
		log.Fatal(newFirestoreErr)
	}

	firebaseFirestore = fbFirestore
}

func main() {
	fmt.Println("CDN v1.0 - Copyright (c) 2021 Anthony Collier")

	handleRequests()
}
