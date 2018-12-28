package swimmy

import "time"

//Card is data for struct and output card or json
type Card struct {
	LinkPage PageData
}

//PageData is a struct for storage data(information) of web page specified with url in order to create embed card or json data
type PageData struct {
	URL         string
	ContentType string
	Title       string
	Description string
	FaviconURL  string
	Image       *ImageData
	Ogp         *OGP
}

//OGP is open graph protocol data
type OGP struct {
	URL          string
	SiteName     string
	Title        string
	Description  string
	OgImage      *ImageData
	TwitterImage *ImageData
	UpdatedTime  *time.Time
	OtherAttrs   map[string]string
}

//ImageData storage properties of image
type ImageData struct {
	URL        string
	SecureURL  string
	FormatType string
	width      int
	height     int
}

//NewImageData return new instance of ImageData
func NewImageData(url, secureURL, formatType string, width, height int) *ImageData {
	return &ImageData{url, secureURL, formatType, width, height}
}

//NewEmptyImageData return new initialized(emply) instance of ImageData
func NewEmptyImageData() *ImageData {
	return &ImageData{"", "", "", -1, -1}
}

//NewPageData return new instance of PageData
func NewPageData(url string, ctype string) *PageData {
	return &PageData{url, ctype, "", "", "", nil, nil}
}

//NewOGP return new instance of OGP
func NewOGP() *OGP {
	return &OGP{"", "", "", "", nil, nil, nil, nil}
}
