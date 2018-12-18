package swimmy

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

//SimpleFetch is download html object with Default Client
func SimpleFetch(url string) (string, error){
	return FetchWithGivenHTTPClient(url, http.DefaultClient)
}

//FetchWithGivenHTTPClient fetch url contents with custom http client
func FetchWithGivenHTTPClient(url string, hc *http.Client) (string, error){
	res, err := hc.Get(url)
	if err != nil {
		return "",err
	}
	defer res.Body.Close()
	contentType := res.Header.Get("Content-Type")
	if contentType == "text/html" || contentType == "text/plain" {
		cbyte, err := ioutil.ReadAll(res.Body)
		return string(cbyte),err
	}

	return "", fmt.Errorf("Invalid Content-Type: %s", contentType)
}
