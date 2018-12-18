package swimmy

import (
	"github.com/microcosm-cc/bluemonday"
)

//Sanitize sanitize html or txt with blueMonday
func Sanitize(htmlContent string, policy *bluemonday.Policy) string {
	return policy.Sanitize(htmlContent)
}

//Parser is url content parser
type Parser struct {
	SanitizePolicy *bluemonday.Policy
}

//NewParser generate NewParser
func (p *Parser) NewParser(custompolicy ...*bluemonday.Policy) *Parser {
	if len(custompolicy) < 1 {
		return &Parser{CPolicy()}
	}
	return &Parser{custompolicy[0]}
}

//Sanitize sanitize html content with p's sanitize policy.
func (p *Parser) Sanitize(htmlContent string) string {
	return Sanitize(htmlContent, p.SanitizePolicy)
}

/*
Parse parse html content.
Before parsing, Parse sanitize html content with its SanitizePolicy.
*/
func (p *Parser) Parse(htmlContent string) {
	sanitizedContent := Sanitize(htmlContent, p.SanitizePolicy)

}

//CPolicy return default policy of swimmy
func CPolicy() *bluemonday.Policy {
	cp := bluemonday.NewPolicy()

	cp.AllowElements("head")
	cp.AllowElements("body")
	cp.AllowElements("title")
	cp.AllowAttrs("name", "content", "property").OnElements("meta")

	return cp
}
