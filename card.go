/*
swimmy is a package for fetch URL Info for embedding external site information as card or outputting as JSON.

First: swimmy fetch url information (html document and text document).
Second: swimmy sanitize html contents and parse it in order to get the information of webpage.
*/

package swimmy

//Card is data for struct and output card or json
type Card struct {
	LinkPage PageData
}
