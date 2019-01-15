/*Package swimmy is a package that fetch and process URL Info for embedding external site information as card or outputting as JSON.
First: swimmy fetch url information (html document and text document).
Second: swimmy sanitize html contents and parse it in order to get the information of webpage.*/
package swimmy

import "github.com/microcosm-cc/bluemonday"

//DefaultContentFetcher is swimmy's defaultContentFetcher
var DefaultContentFetcher *ContentFetcher

//DefaultPageDataBuilder is swimmy's default PageDataBuilder
var DefaultPageDataBuilder *PageDataBuilder

//DefaultCardBuilder is swimmy's default CardDataBuilder
var DefaultCardBuilder *CardBuilder

//IDCount count of PageData's ID
var IDCount int

func init() {
	DefaultContentFetcher = NewContentFetcher(nil)
	DefaultPageDataBuilder = NewPageDataBuilder(CPolicy(), TPolicy())
	DefaultCardBuilder = DefSetCardBuilder()
	IDCount = 0
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
