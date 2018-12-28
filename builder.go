package swimmy

import (
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

func getMarkedUpText(ct *html.Tokenizer, tagName string) {

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
