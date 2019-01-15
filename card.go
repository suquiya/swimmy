package swimmy

import (
	"html/template"
	"strings"
)

//CardBuilder build card string from pagedata
type CardBuilder struct {
	Template   *template.Template
	ClassNames map[string]string
}

//NewCardBuilder create a empty instance of CardBuilder and return it
func NewCardBuilder() *CardBuilder {
	return &CardBuilder{nil, nil}
}

//DefSetCardBuilder return a CardBuilder which set swimmy's default template and Classes.
func DefSetCardBuilder() *CardBuilder {
	return &CardBuilder{DefaultTemplate(), DefaultClasses()}
}

//BuildCard build card as string
func (cb *CardBuilder) BuildCard(pd *PageData) string {
	m := make(map[string]interface{})

	m["ClassNames"] = cb.ClassNames
	m["PageData"] = pd

	var sb strings.Builder
	cb.Template.Execute(&sb, m)

	return sb.String()
}

//DefaultClasses return default classNames in card as map
func DefaultClasses() map[string]string {
	cns := make(map[string]string)

	cns["CardDiv"] = "swimmy-card"
	cns["SiteInfo"] = "sc-info"
	cns["PageInfo"] = "sc-contents"
	cns["PageImageWrapper"] = "sc-image-wrapper"
	cns["PageImage"] = "sc-image"
	cns["PageTitle"] = "sc-title"
	cns["PageURL"] = "sc-url"

	return cns
}

//DefaultTemplate return swimmy's default template
func DefaultTemplate() *template.Template {
	str := `
	<div class="{{.ClassNames.CardDiv}}" id="swimmy-{{.PageData.ID}}"><a href="{{.PageData.URL}}">
	<div class="{{.ClassNames.SiteInfo}}">{{ .PageData.OGP.SiteName }}</div>
	<div class="{{.ClassNames.PageInfo}}">
	<div class="{{.ClassNames.PageImageWrapper}}"><img class="{{.ClassNames.PageImage}}" /></div>
	<a href="{{.PageData.URL}}" class="{{.ClassNames.PageTitle}}">{{.Title}}</a>
	<a href="{{.URL}}" class="{{.ClassNames.PageURL}}">{{.URL}}</a>
	</div>
	</a></div>
	`

	tmpl, err := template.New("cardTemplate").Parse(str)
	if err != nil {
		panic(err)
	}

	return tmpl
}
