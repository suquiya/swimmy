package swimmy

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

//SimpleFetch is download html object with Default Client
func SimpleFetch(url string) (*PageData, string, error) {
	return FetchWithGivenHTTPClient(url, http.DefaultClient)
}

//FetchWithGivenHTTPClient fetch url contents with custom http client. If you don't use your custom policy, use SimpleFetch(url string).
func FetchWithGivenHTTPClient(url string, hc *http.Client) (*PageData, string, error) {
	res, err := hc.Get(url)
	if err != nil {
		return nil, "", err
	}
	defer res.Body.Close()
	contentType := res.Header.Get("Content-Type")
	if contentType == "text/html" || contentType == "text/plain" {
		cbyte, err := ioutil.ReadAll(res.Body)
		return NewPageData(url, contentType), string(cbyte), err
	}

	return nil, "", fmt.Errorf("Invalid Content-Type: %s", contentType)
}
