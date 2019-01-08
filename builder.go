package swimmy

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"

	"golang.org/x/net/html/charset"

	"github.com/asaskevich/govalidator"
	"github.com/microcosm-cc/bluemonday"
	"golang.org/x/net/html"
)

//Sanitize sanitize html or txt with blueMonday
func Sanitize(htmlContent string, policy ...*bluemonday.Policy) string {
	if len(policy) < 1 {
		return DefaultPageDataBuilder.PreSanitizePolicy.Sanitize(htmlContent)
	}
	return policy[0].Sanitize(htmlContent)
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

//DefaultPageDataBuilder is swimmy's default pagedatabuilder
var DefaultPageDataBuilder *PageDataBuilder

//NewPageDataBuilder generate New instance of PageDataBuilder
func NewPageDataBuilder(PrePolicy, tagContentPolicy *bluemonday.Policy) *PageDataBuilder {
	return &PageDataBuilder{PrePolicy, tagContentPolicy}
}

//Sanitize sanitize html content with p's sanitize policy.
func (p *PageDataBuilder) Sanitize(htmlContent string) string {
	return Sanitize(htmlContent, p.PreSanitizePolicy)
}

/*
BuildPageData parse html content, retrieve tag info and fill PageData.
Before parsing, Parse sanitize html content with its SanitizePolicy.
*/
func (p *PageDataBuilder) BuildPageData(pd *PageData, htmlContent string) *PageData {

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
								nameAttr = string(key)
							case bytes.Equal(key, cstrb):
								contentAttr = string(val)
							}
						}

						p.TagContentSanitize(nameAttr)
						p.TagContentSanitize(contentAttr)
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
						case strings.HasPrefix(nameAttr, "og:"):
							pd.OGP.Set(nameAttr, contentAttr)

						}

					}
				case "title":
					pd.Title = TakeMarkedUpText(cTokenizer, tnByte)
				}
			}
		}
	} else {
		fmt.Println(WhyCannotTokenize)
	}

	return pd

}

//ParseTime parse time string
func ParseTime(timeStr string) (*time.Time, string, error) {
	formatStrings := []string{
		"2006-01-02T15:04:05-07:00",
		"2006-01-02T15:04:05Z",
		"20060102T150405-0700",
		"20060102T150405Z"}

	format := ""
	var resultTime *time.Time
	resultTime = nil
	var err error
	err = nil
	for _, fStr := range formatStrings {
		t, err := time.Parse(fStr, timeStr)
		if err == nil {
			format = fStr
			resultTime = &t
			break
		}
	}
	if resultTime == nil {
		err = fmt.Errorf("error: cannot parse format")
	}
	return resultTime, format, err
}

//TakeMarkedUpText is take marked-up text between begin tag and end tag
func TakeMarkedUpText(ct *html.Tokenizer, tagName []byte) string {
	depth := 0
	taking := true
	var sb *strings.Builder

	for taking {
		tt := ct.Next()
		switch tt {
		case html.StartTagToken:
			tName, _ := ct.TagName()
			if bytes.Equal(tName, tagName) {
				depth++
			}
			WriteCurrentString(ct, tt, sb)

		case html.EndTagToken:
			tName, _ := ct.TagName()
			if bytes.Equal(tName, tagName) {
				depth--
				if depth < 1 {
					taking = false
				} else {
					WriteCurrentString(ct, tt, sb)
				}
			} else {
				WriteCurrentString(ct, tt, sb)
			}
		case html.ErrorToken:
			taking = false
		default:
			WriteCurrentString(ct, tt, sb)
		}
	}
	return sb.String()
}

//WriteCurrentString write string of now tag or text to strings.Builder
func WriteCurrentString(tokenizer *html.Tokenizer, tokenType html.TokenType, sb *strings.Builder) {
	switch tokenType {
	case html.ErrorToken:
	case html.TextToken:
		sb.WriteString(EscapeBytes(tokenizer.Text()))
	case html.StartTagToken:
		sb.WriteString("<")
		writeCurrentTagString(tokenizer, sb)
		sb.WriteString(">")
	case html.EndTagToken:
		sb.WriteString("</")
		tName, _ := tokenizer.TagName()
		sb.Write(tName)
		sb.WriteString(">")
	case html.SelfClosingTagToken:
		sb.WriteString("<")
		writeCurrentTagString(tokenizer, sb)
		sb.WriteString("/>")
	case html.CommentToken:
		sb.WriteString("<!--")
		sb.WriteString(EscapeBytes(tokenizer.Text()))
		sb.WriteString("-->")
	case html.DoctypeToken:
		sb.WriteString("<!DOCTYPE ")
		sb.WriteString(EscapeBytes(tokenizer.Text()))
		sb.WriteString(">")
	default:
		sb.WriteString("Invalid token:< ")
		sb.WriteString(strconv.Itoa(int(tokenType)))
		sb.WriteString(">")
	}
}

func writeCurrentTagString(tokenizer *html.Tokenizer, sb *strings.Builder) {
	tName, hasAttr := tokenizer.TagName()
	if hasAttr {
		sb.Write(tName)

		moreAttr := true
		var key, val []byte
		for moreAttr {
			key, val, moreAttr = tokenizer.TagAttr()

			sb.WriteByte(' ')
			sb.WriteString(EscapeBytes(key))
			sb.WriteString(`="`)
			sb.WriteString(EscapeBytes(val))
			sb.WriteByte('"')
		}
	} else {
		sb.Write(tName)
	}
}

//EscapeBytes html escape byte array
func EscapeBytes(str []byte) string {
	return html.EscapeString(string(str))
}

//CPolicy return default policy of swimmy
func CPolicy() *bluemonday.Policy {
	cp := bluemonday.NewPolicy()

	cp.AllowElements("head")
	cp.AllowElements("body")
	cp.AllowElements("title")
	cp.AllowAttrs("name", "content", "property").OnElements("meta")
	cp.AllowAttrs("rel", "href").OnElements("link")

	return cp
}

//TPolicy return default tag policy of swimmy
func TPolicy() *bluemonday.Policy {
	tp := bluemonday.NewPolicy()

	return tp
}
