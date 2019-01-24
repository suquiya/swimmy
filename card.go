package swimmy

import (
	"html/template"
	"io"
)

//CardBuilder build card string from pagedata
type CardBuilder struct {
	CardTemplate *template.Template
	ClassNames   map[string]string
}

//NewCardBuilder create a empty instance of CardBuilder and return it
func NewCardBuilder() *CardBuilder {
	return &CardBuilder{nil, nil}
}

//DefSetCardBuilder return a CardBuilder which set swimmy's default template and Classes.
func DefSetCardBuilder() *CardBuilder {
	return &CardBuilder{DefaultTemplate(), DefaultClasses()}
}

//Execute build card by execute html template.
func (cb *CardBuilder) Execute(pd *PageData, w io.Writer) error {
	d := make(map[string]interface{})

	d["ClassNames"] = cb.ClassNames
	d["PageData"] = pd

	err := cb.CardTemplate.Execute(w, d)

	return err
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
	cns["PageDescription"] = "sc-description"

	return cns
}

//DefaultTemplate return swimmy's default template
func DefaultTemplate() *template.Template {

	str := `<div class="{{.ClassNames.CardDiv}}" id="swimmy-{{.PageData.ID}}">
	<a href="{{.PageData.URL}}">
	<div class="{{.ClassNames.SiteInfo}}">{{ .PageData.OGP.SiteName }}</div>
	<div class="{{.ClassNames.PageInfo}}">
	<div class="{{.ClassNames.PageImageWrapper}}">
	<img class="{{.ClassNames.PageImage}}" src="{{.PageData.OGP.OgImage.URL}}" />
	</div>
	<a href="{{.PageData.URL}}" class="{{.ClassNames.PageTitle}}">{{.PageData.Title}}</a>
	<a href="{{.PageData.URL}}" class="{{.ClassNames.PageURL}}">{{.PageData.URL}}</a>
	<div class="{{.ClassNames.PageDescription}}">
	{{.PageData.Description}}
	</div>
	</div>
	</a>
	</div>`

	//str := "test\r\n"
	tmpl, err := template.New("DefaultCard").Parse(str)
	if err != nil {
		panic(err)
	}

	return tmpl
}
