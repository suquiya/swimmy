package swimmy

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"
	"time"

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
