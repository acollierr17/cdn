package main

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/joho/godotenv"
	"log"
	"os"
)

var cdnConfig *Config
var s3Config *aws.Config

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	cdnConfig = &Config{
		SpacesConfig: SpacesConfig{
			SpacesAccessKey: os.Getenv("SPACES_ACCESS_KEY"),
			SpacesSecretKey: os.Getenv("SPACES_SECRET_KEY"),
			SpacesEndpoint:  os.Getenv("SPACES_ENDPOINT"),
			SpacesUrl:       os.Getenv("SPACES_URL"),
			SpacesName:      os.Getenv("SPACES_NAME"),
			SpacesRegion:    os.Getenv("SPACES_REGION"),
		},
		CdnEndpoint: os.Getenv("CDN_ENDPOINT"),
		AccessToken: os.Getenv("ACCESS_TOKEN"),
	}

	s3Config = &aws.Config{
		Credentials: credentials.NewStaticCredentials(cdnConfig.SpacesConfig.SpacesAccessKey, cdnConfig.SpacesConfig.SpacesSecretKey, ""),
		Endpoint:    aws.String(cdnConfig.SpacesConfig.SpacesEndpoint),
		Region:      aws.String(cdnConfig.SpacesConfig.SpacesRegion),
	}
}

func main() {
	fmt.Println("CDN v1.0 - Copyright (c) 2021 Anthony Collier")

	handleRequests()
}
