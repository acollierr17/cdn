package main

import "time"

type Embed struct {
	Type         string `json:"type"`
	AuthorName   string `json:"author_name"`
	ProviderName string `json:"provider_name"`
}

type ImageResponse struct {
	Url     string `json:"url"`
	Success bool   `json:"success"`
}

type ImageResponseQuery struct {
	Download	string	`query:"download"`
}

type ImageDeletedRespone struct {
	ImageName string `json:"image_name"`
	Deleted   bool   `json:"deleted"`
}

type SpacesConfig struct {
	SpacesAccessKey string
	SpacesSecretKey string
	SpacesEndpoint  string
	SpacesUrl       string
	SpacesCdn	    string
	SpacesName      string
	SpacesRegion    string
}

type Config struct {
	SpacesConfig SpacesConfig
	CdnEndpoint  string
	AccessToken  string
}

type AuthTokenRequest struct {
	Token	string	`json:"token"`
}

type ImageResult struct {
	CdnUrl			    string		`json:"cdn_url"`
	SpacesUrl		    string		`json:"spaces_url"`
	SpacesCdn		    string		`json:"spaces_cdn"`
	FileName		    string		`json:"file_name"`
	LastModified		time.Time	`json:"last_modified"`
	Size			    int64		`json:"size"`
}

type ImagesResults struct {
	Images	[]*ImageResult	`json:"images"`
	Length	int		        `json:"length"`
}
