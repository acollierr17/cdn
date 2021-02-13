package main

type Embed struct {
	Type         string `json:"type"`
	AuthorName   string `json:"author_name"`
	ProviderName string `json:"provider_name"`
}

type ImageResponse struct {
	Url     string `json:"url"`
	Success bool   `json:"success"`
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
	SpacesName      string
	SpacesRegion    string
}

type Config struct {
	SpacesConfig SpacesConfig
	CdnEndpoint  string
	AccessToken  string
}

type SPAHandler struct {
	StaticPath string
	IndexPath  string
}
