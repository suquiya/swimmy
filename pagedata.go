package swimmy

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/url"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/asaskevich/govalidator"
	"github.com/microcosm-cc/bluemonday"
	"golang.org/x/net/html"
	"golang.org/x/net/html/charset"
)

//PageData is a struct for storage data(information) of web page specified with url in order to create embed card or json data
type PageData struct {
	URL           string             `json:"URL"`
	ID            int                `json:"ID"`
	CannonicalURL string             `json:"CannonicalURL"`
	ContentType   string             `json:"ContentType"`
	Title         string             `json:"Title"`
	Description   string             `json:"Description"`
	FaviconURL    []string           `json:"FaviconURL"`
	OGP           *OpenGraphProtocol `json:"OGP"`
}

//OpenGraphProtocol is strage for open graph protocol. OpenGraphProtocol in swimmy is only for creating data for embedding in website, so it does not storage video and music.
type OpenGraphProtocol struct {
	URL          string            `json:"URL"`
	SiteName     string            `json:"SiteName"`
	Title        string            `json:"Title"`
	Description  string            `json:"Description"`
	Locale       string            `json:"Locale"`
	Type         string            `json:"Type"`
	OgImage      *ImageData        `json:"OgImage"`
	TwitterImage *ImageData        `json:"TwitterImage"`
	TwitterID    string            `json:"TwitterID"`
	UpdatedTime  *time.Time        `json:"UpdatedTime"`
	OtherAttrs   map[string]string `json:"OtherAttrs"`
	OtherInfo    map[string]string `json:"OtherInfo"`
}

//Set set meta values to ogp fields. contentAttr is assumed after sanitizing.
func (ogp *OpenGraphProtocol) Set(nameAttr, contentAttr string) {
	if strings.HasPrefix(nameAttr, "og:") {
		val := strings.TrimLeft(nameAttr, "og:")
		catched := true
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
			catched = false
		}

		if !catched {
			if strings.HasPrefix(nameAttr, "og:image") {
				attr := strings.TrimLeft(nameAttr, "og:image")
				catched = true
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
				default:
					catched = false
				}
			}
		}

		if !catched {
			ogp.OtherAttrs[nameAttr] = contentAttr
		}
	} else if strings.HasPrefix(nameAttr, "twitter:") {
		tname := strings.TrimLeft(nameAttr, "twitter:")
		switch tname {
		case "image":
			if govalidator.IsRequestURL(contentAttr) {
				ogp.TwitterImage.URL = contentAttr
			}
		case "site":
			if strings.HasPrefix(contentAttr, "@") {
				ogp.TwitterID = html.EscapeString(contentAttr)
			}
		default:
			ogp.OtherAttrs[nameAttr] = contentAttr
		}
	} else {
		ogp.OtherAttrs[nameAttr] = contentAttr
	}
}

//ImageData storage properties of image
type ImageData struct {
	URL        string `json:"URL"`
	SecureURL  string `json:"SecureURL"`
	FormatType string `json:"FormatType"`
	AltText    string `json:"AltText"`
	Width      int    `json:"Width"`
	Height     int    `json:"Height"`
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
	npd := &PageData{url, IDCount, ctype, "", "", "", make([]string, 0, 1), NewOGP()}
	IDCount++
	return npd
}

//NewOGP return new instance of OGP
func NewOGP() *OpenGraphProtocol {
	return &OpenGraphProtocol{"", "", "", "", "", "", NewImageData(), NewImageData(), "", nil, make(map[string]string), make(map[string]string)}
}

//PageDataBuilder is processer for creating pagedata
type PageDataBuilder struct {
	PreSanitizePolicy        *bluemonday.Policy
	TagContentSanitizePolicy *bluemonday.Policy
}

//TagContentSanitize sanitize content of tag
func (p *PageDataBuilder) TagContentSanitize(str string) string {
	return p.TagContentSanitizePolicy.Sanitize(str)
}

//NewPageDataBuilder generate New instance of PageDataBuilder
func NewPageDataBuilder(PrePolicy, tagContentPolicy *bluemonday.Policy) *PageDataBuilder {
	return &PageDataBuilder{PrePolicy, tagContentPolicy}
}

//Sanitize sanitize html content with p's sanitize policy.
func (p *PageDataBuilder) Sanitize(htmlContent string) string {
	return Sanitize(htmlContent, p.PreSanitizePolicy)
}

//BuildPageData build pagedata on base pagedata
func BuildPageData(url string, ctype string, htmlContent string) *PageData {
	return DefaultPageDataBuilder.BuildPageData(url, ctype, htmlContent)
}

//ErrorPageData return pagedata if get err in fetch.
func ErrorPageData(url, ctype string, content []byte, err error) *PageData {
	pd := NewPageData(url, ctype)

	if err.Error() == "input is not URL" {

		pd.URL = ""
		pd.CannonicalURL = ""
		d := html.EscapeString(DefaultPageDataBuilder.TagContentSanitize(url))
		pd.Title = "Link value is not URL: " + d
		pd.Description = "リンクとして指定された値" + d + "は正しくなかったためこのカードは白紙の状態で提示されます"
		return pd
	}
	if err.Error() == "statusError" {
		pd.Title = html.EscapeString(string(content))
		pd.Description = fmt.Sprintf("%s - %s", pd.Title, url)

		return pd
	}

	if strings.HasPrefix(err.Error(), "Invalid Content-Type: ") {

		pd.Description = ""
		pd.Title = url

		return pd
	}

	return pd
}

/*
BuildPageData parse html content, retrieve tag info and fill PageData.
Before parsing, Parse sanitize html content with its SanitizePolicy.
*/
func (p *PageDataBuilder) BuildPageData(url string, ctype string, htmlContent string) *PageData {

	pd := NewPageData(url, ctype)
	sanitizedContent := Sanitize(htmlContent, p.PreSanitizePolicy)
	canTokenize := true
	WhyCannotTokenize := ""
	if !utf8.ValidString(sanitizedContent) {
		sr := strings.NewReader(sanitizedContent)
		scByte, err := bufio.NewReader(sr).Peek(1024)
		if err != nil {
			panic(err)
		}
		e, name, _ := charset.DetermineEncoding(scByte, pd.ContentType)
		sr = strings.NewReader(sanitizedContent)
		if e != nil {
			r := e.NewDecoder().Reader(sr)
			scb, err := ioutil.ReadAll(r)
			if err != nil {
				panic(err)
			}
			sanitizedContent = string(scb)
			sanitizedContent = Sanitize(htmlContent, p.PreSanitizePolicy)
		} else {
			fmt.Printf("bad encode: %s", name)
			canTokenize = false
			WhyCannotTokenize = "cannot htmlContents tokenize because of content's charset encoding"
		}
	}

	if strings.HasPrefix(pd.ContentType, "text/plain") {
		if canTokenize {
			canTokenize = false
			WhyCannotTokenize = "cannot tokenize because of contentType is text"
		} else {
			WhyCannotTokenize = WhyCannotTokenize + "\r\ncannot tokenize because of contentType is text"
		}

		pd.Description = html.EscapeString(sanitizedContent)
	}

	if canTokenize {
		ContentReader := strings.NewReader(sanitizedContent)

		cTokenizer := html.NewTokenizer(ContentReader)

		parse := true

		metaNameEmptyCount := 0
		for parse {
			tt := cTokenizer.Next()

			parse = tt != html.ErrorToken

			if parse && tt != html.EndTagToken {
				tnByte, hasAttr := cTokenizer.TagName()
				tn := string(tnByte)
				switch tn {
				case "meta":
					if hasAttr {
						moreAttr := true
						var key, val []byte
						nameAttr := ""
						nstrb := []byte("name")
						contentAttr := ""
						cstrb := []byte("content")
						for moreAttr {
							key, val, moreAttr = cTokenizer.TagAttr()
							switch {
							case bytes.Equal(key, nstrb):
								nameAttr = string(val)
							case bytes.Equal(key, cstrb):
								contentAttr = string(val)
							}
						}

						nameAttr = p.TagContentSanitize(nameAttr)
						nameAttr = html.EscapeString(nameAttr)
						contentAttr = p.TagContentSanitize(contentAttr)
						switch {
						case nameAttr == "":
							if contentAttr != "" {
								metaNameEmptyCount++
								pd.OGP.OtherAttrs["empty"+strconv.Itoa(metaNameEmptyCount)] = contentAttr
							}
						case nameAttr == "description":
							pd.Description = html.EscapeString(contentAttr)
						case nameAttr == "cannonical":
							if govalidator.IsRequestURL(contentAttr) {
								pd.CannonicalURL = contentAttr
							}
						case strings.HasPrefix(nameAttr, "og:") || strings.HasPrefix(nameAttr, "twitter:"):
							pd.OGP.Set(nameAttr, contentAttr)
						default:
							pd.OGP.OtherAttrs[nameAttr] = contentAttr
						}
					}
				case "title":
					pd.Title = TakeMarkedUpText(cTokenizer, tnByte)
				case "link":
					if hasAttr {
						moreAttr := true
						var key, val []byte
						relAttr := ""
						rstrb := []byte("rel")
						hrefAttr := ""
						hstrb := []byte("href")

						for moreAttr {
							key, val, moreAttr = cTokenizer.TagAttr()
							switch {
							case bytes.Equal(key, rstrb):
								relAttr = string(val)
							case bytes.Equal(key, hstrb):
								hrefAttr = string(val)
							}
						}

						relAttr = p.TagContentSanitize(relAttr)
						relAttr = html.EscapeString(relAttr)
						hrefAttr = p.TagContentSanitize(hrefAttr)
						switch relAttr {
						case "cannonical":
							if govalidator.IsRequestURL(hrefAttr) {
								pd.CannonicalURL = hrefAttr
							}
						case "icon":
							if govalidator.IsRequestURL(hrefAttr) {
								pd.FaviconURL = append(pd.FaviconURL, hrefAttr)
							}
						}
					}
				case "body":
					if pd.Description == "" {
						pd.Description = html.EscapeString(TakeMarkedUpText(cTokenizer, tnByte))
					}
					parse = false
				}
			}
		}
	} else {
		fmt.Println(WhyCannotTokenize)
	}

	return pd

}

//ComplementBasicFields complement pagedata basic fields if some basic field is empty.
func (pd *PageData) ComplementBasicFields() {
	if pd.IsPlainText() {
		pd.Title = url.PathEscape(pd.URL)
	} else {
		if pd.CannonicalURL == "" {
			pd.CannonicalURL = pd.URL
		}
		if pd.Title == "" {
			pd.Title = pd.OGP.Title
		}
		if pd.Description == "" {
			pd.Description = pd.OGP.Description
		}
		if pd.OGP.OgImage.URL == "" {
			pd.OGP.OgImage = pd.OGP.TwitterImage
		}
	}
}

//ToJSON convert PageData to json.
func (pd *PageData) ToJSON() ([]byte, error) {
	j, err := json.Marshal(pd)
	if err != nil {
		return nil, err
	}
	return j, err
}

//IsPlainText return whether pagedata is text/plain or not.
func (pd *PageData) IsPlainText() bool {
	return strings.HasPrefix(pd.ContentType, "text/plain")
}

//IsPlainTextContentType return whether given contentType represents text/plain or not
func IsPlainTextContentType(ctype string) bool {
	return strings.HasPrefix(ctype, "text/plain")
}
