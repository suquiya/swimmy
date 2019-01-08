package swimmy

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/asaskevich/govalidator"
	"golang.org/x/net/html"
)

//PageData is a struct for storage data(information) of web page specified with url in order to create embed card or json data
type PageData struct {
	URL           string
	CannonicalURL string
	ContentType   string
	Title         string
	Description   string
	FaviconURL    string
	OGP           *OpenGraphProtocol
}

//OpenGraphProtocol is strage for open graph protocol. OpenGraphProtocol in swimmy is only for creating data for embedding in website, so it does not storage video and music.
type OpenGraphProtocol struct {
	URL          string
	SiteName     string
	Title        string
	Description  string
	Locale       string
	Type         string
	OgImage      *ImageData
	TwitterImage *ImageData
	UpdatedTime  *time.Time
	OtherAttrs   map[string]string
	OtherInfo    map[string]string
}

//Set set meta values to ogp fields
func (ogp *OpenGraphProtocol) Set(nameAttr, contentAttr string) {
	if strings.HasPrefix(nameAttr, "og:") {
		val := strings.TrimLeft(nameAttr, "og:")
		other := false
		switch val {
		case "locale":
			ogp.Locale = html.EscapeString(contentAttr)
		case "title":
			ogp.Title = html.EscapeString(contentAttr)
		case "type":
			ogp.Type = html.EscapeString(contentAttr)
		case "description":
			ogp.Type = html.EscapeString(contentAttr)
		case "url":
			if govalidator.IsRequestURL(contentAttr) {
				ogp.URL = contentAttr
			}
		case "site_name":
			ogp.SiteName = html.EscapeString(contentAttr)
		case "updated_time":
			timeString := html.EscapeString(contentAttr)
			t, f, err := ParseTime(timeString)
			if err == nil {
				ogp.UpdatedTime = t
				ogp.OtherInfo["TimeFormat"] = f
			} else {
				fmt.Println(err)
				ogp.UpdatedTime = nil
				ogp.OtherInfo["TimeFormat"] = "unknown"
				ogp.OtherInfo["UpdatedTimeString"] = "timeString"
			}
		default:
			other = true
		}

		if other {
			if strings.HasPrefix(nameAttr, "og:image") {
				attr := strings.TrimLeft(nameAttr, "og:image")
				switch attr {
				case "":
					if govalidator.IsRequestURL(contentAttr) {
						ogp.OgImage.URL = contentAttr
					}
				case "secure_url":
					if strings.HasPrefix(contentAttr, "https://") && govalidator.IsRequestURL(contentAttr) {
						ogp.OgImage.SecureURL = contentAttr
					}
				case "width":
					c := html.EscapeString(contentAttr)
					if govalidator.IsInt(c) {
						w, _ := strconv.Atoi(c)
						ogp.OgImage.Width = w
					}
				case "height":
					c := html.EscapeString(contentAttr)
					if govalidator.IsInt(c) {
						h, _ := strconv.Atoi(c)
						ogp.OgImage.Height = h
					}
				case "alt":
					ogp.OgImage.AltText = html.EscapeString(contentAttr)

				}
			}
		}
	}
}

//ImageData storage properties of image
type ImageData struct {
	URL        string
	SecureURL  string
	FormatType string
	AltText    string
	Width      int
	Height     int
}

//CreateImageData return new instance of ImageData
func CreateImageData(url, secureURL, formatType, alt string, width, height int) *ImageData {
	return &ImageData{url, secureURL, formatType, alt, width, height}
}

//NewImageData return new initialized(emply) instance of ImageData
func NewImageData() *ImageData {
	return &ImageData{"", "", "", "", -1, -1}
}

//NewPageData return new instance of PageData
func NewPageData(url string, ctype string) *PageData {
	return &PageData{url, ctype, "", "", "", "", NewOGP()}
}

//NewOGP return new instance of OGP
func NewOGP() *OpenGraphProtocol {
	return &OpenGraphProtocol{"", "", "", "", "", "", NewImageData(), nil, nil, make(map[string]string), make(map[string]string)}
}
