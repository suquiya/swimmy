/*Package swimmy is a package that fetch and process URL Info for embedding external site information as card or outputting as JSON.
First: swimmy fetch url information (html document and text document).
Second: swimmy sanitize html contents and parse it in order to get the information of webpage.*/
package swimmy

//DefaultContentFetcher is swimmy's defaultContentFetcher
var DefaultContentFetcher *ContentFetcher

func init() {
	DefaultContentFetcher = NewContentFetcher(nil)
	DefaultPageDataBuilder = NewPageDataBuilder(CPolicy(), TPolicy())
}
