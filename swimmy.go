/*
Package swimmy is a package that fetch and process URL Info for embedding external site information as card or outputting as JSON.
First: swimmy fetch url information (html document and text document).
Second: swimmy sanitize html contents and parse it in order to get the information of webpage.*/
package swimmy

import (
	"fmt"
	"time"

	"github.com/microcosm-cc/bluemonday"
	"github.com/spf13/cobra/cobra/cmd"
)

//DefaultContentFetcher is swimmy's defaultContentFetcher
var DefaultContentFetcher *ContentFetcher

//DefaultPageDataBuilder is swimmy's default PageDataBuilder
var DefaultPageDataBuilder *PageDataBuilder

//DefaultCardBuilder is swimmy's default CardDataBuilder
var DefaultCardBuilder *CardBuilder

//IDCount count of PageData's ID
var IDCount int

//Init is initialize function. If you want to use default variables, use this function.
func Init() {
	DefaultContentFetcher = NewContentFetcher(nil)
	DefaultPageDataBuilder = NewPageDataBuilder(CPolicy(), TPolicy())
	DefaultCardBuilder = DefSetCardBuilder()
	IDCount = 0
}

//ShowLicense show license of swimmy
func ShowLicense() string {
	swimmyLicense := cmd.Licenses["bsd"]
	data := make(map[string]string)
	data["copyright"] = "Copyright (c) " + time.Now().Format("2006") + "suquiya"
	ltxt, err := ExecLicenseTextTemp(swimmyLicense.Text, data)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	return ltxt
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
