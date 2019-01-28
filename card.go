package swimmy

import (
	"bytes"
	"fmt"
	"html/template"
	"io"

	min "github.com/tdewolff/minify/v2"
	"github.com/tdewolff/minify/v2/html"
)

//CardBuilder build card string from pagedata
type CardBuilder struct {
	CardTemplate *template.Template
	ClassNames   map[string]string
}

//NewCardBuilder create a empty instance of CardBuilder and return it
func NewCardBuilder(cardtemplate *template.Template, classnames map[string]string) *CardBuilder {
	return &CardBuilder{cardtemplate, classnames}
}

//WriteCardHTML write card html tag.
func (cb *CardBuilder) WriteCardHTML(pd *PageData, w io.Writer, minify bool) {
	if minify {
		var bb bytes.Buffer
		cb.Execute(pd, &bb)
		m := min.New()
		m.AddFunc("text/html", html.Minify)
		b, err := m.Bytes("text/html", bb.Bytes())
		if err != nil {
			panic(err)
		}
		w.Write(b)
	} else {
		err := cb.Execute(pd, w)
		if err != nil {
			fmt.Println(err)
		}
	}
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

	str := `<div class="{{.ClassNames.CardDiv}}" id="swimmy-{{.PageData.ID}}"><a href="{{.PageData.URL}}"><div class="{{.ClassNames.SiteInfo}}">{{ .PageData.OGP.SiteName }}</div><div class="{{.ClassNames.PageInfo}}"><div class="{{.ClassNames.PageImageWrapper}}"><img class="{{.ClassNames.PageImage}}" src="{{.PageData.OGP.OgImage.URL}}" /></div><a href="{{.PageData.URL}}" class="{{.ClassNames.PageTitle}}">{{.PageData.Title}}</a><a href="{{.PageData.URL}}" class="{{.ClassNames.PageURL}}">{{.PageData.URL}}</a><div class="{{.ClassNames.PageDescription}}">{{.PageData.Description}}</div></div></a></div>`

	//str := "test\r\n"
	tmpl, err := template.New("DefaultCard").Parse(str)
	if err != nil {
		panic(err)
	}

	return tmpl
}
