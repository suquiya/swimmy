package swimmy

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

//HTTPClient keep httpClient if user set custom Client
var HTTPClient *http.Client

//Fetch fetch url contents with f's HTTPClient. If you don't set your custom http client, Fetch use DefaultClient of net/http package.
func Fetch(url string) (string, string, string, error) {

	var res *http.Response
	var err error
	if HTTPClient != nil {
		res, err = HTTPClient.Get(url)
	} else {
		res, err = http.DefaultClient.Get(url)
	}
	if err != nil {
		return "", "", "", err
	}
	defer res.Body.Close()
	contentType := res.Header.Get("Content-Type")
	if strings.HasPrefix(contentType, "text/html") || strings.HasPrefix(contentType, "text/plain") {
		cbyte, err := ioutil.ReadAll(res.Body)
		return url, contentType, string(cbyte), err
	}

	return "", "", "", fmt.Errorf("Invalid Content-Type: %s", contentType)
}
