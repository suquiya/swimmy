package swimmy

import (
	"bytes"
	"strings"

	"github.com/microcosm-cc/bluemonday"
	"golang.org/x/net/html"
)

//Sanitize sanitize html or txt with blueMonday
func Sanitize(htmlContent string, policy *bluemonday.Policy) string {
	return policy.Sanitize(htmlContent)
}

//PageDataBuilder is processer for creating pagedata
type PageDataBuilder struct {
	SanitizePolicy *bluemonday.Policy
}

//NewPageDataBuilder generate New instance of PageDataBuilder
func (p *PageDataBuilder) NewPageDataBuilder(custompolicy ...*bluemonday.Policy) *PageDataBuilder {
	if len(custompolicy) < 1 {
		return &PageDataBuilder{CPolicy()}
	}
	return &PageDataBuilder{custompolicy[0]}
}

//Sanitize sanitize html content with p's sanitize policy.
func (p *PageDataBuilder) Sanitize(htmlContent string) string {
	return Sanitize(htmlContent, p.SanitizePolicy)
}

/*
BuildPageData parse html content, retrieve tag info and fill PageData.
Before parsing, Parse sanitize html content with its SanitizePolicy.
*/
func (p *PageDataBuilder) BuildPageData(pd *PageData, htmlContent string) *PageData {

	sanitizedContent := Sanitize(htmlContent, p.SanitizePolicy)

	ContentReader := strings.NewReader(sanitizedContent)

	cTokenizer := html.NewTokenizer(ContentReader)

	parse := true

	for parse {
		tt := cTokenizer.Next()

		parse = tt != html.ErrorToken

		if parse && tt != html.EndTagToken {
			tnByte, hasAttr := cTokenizer.TagName()
			tn := string(tnByte)
			switch tn {
			case "meta":
				if hasAttr {

				}
			case "title":
			}
		}
	}

	return pd

}

func takeMarkedUpText(ct *html.Tokenizer, tagName string) {
	depth := 0
	taking := true
	var sb *strings.Builder
	tagNameByte := []byte(tagName)

	for taking {
		tt := ct.Next()
		switch tt {
		case html.StartTagToken:
			tName, _ := ct.TagName()
			if bytes.Equal(tName, tagNameByte) {
				depth++
			}

		case html.EndTagToken:
			tName, _ := ct.TagName()
			if bytes.Equal(tName, tagNameByte) {
				depth--
				if depth < 1 {
					taking = false
				} else {

				}
			} else {

			}

		}
	}
}

//WriteCurrentString write string of now tag or text to strings.Builder
func WriteCurrentString(tokenizer *html.Tokenizer, tokenType html.TokenType, sb *strings.Builder) {
	switch tokenType {
	case html.ErrorToken:
		return
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
		sb.WriteString(strings.Itoa(tokenType))
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
