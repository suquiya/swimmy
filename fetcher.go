package swimmy

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"unicode/utf8"

	"github.com/asaskevich/govalidator"

	"golang.org/x/net/html/charset"
)

//ContentFetcher fetch net content in Fetch(url string)
type ContentFetcher struct {
	HTTPClient *http.Client
}

//NewContentFetcher create new instance of ContentFetcher
func NewContentFetcher(HTTPClient *http.Client) *ContentFetcher {
	return &ContentFetcher{HTTPClient}
}

//Fetch with DefaultContentFetcher
func Fetch(url string) (string, string, []byte, error) {
	return DefaultContentFetcher.Fetch(url)
}

//Fetch fetch url contents with f's HTTPClient. If you don't set your custom http client, Fetch use DefaultClient of net/http package.
func (cf *ContentFetcher) Fetch(url string) (string, string, []byte, error) {

	if !govalidator.IsRequestURL(url) {
		return "", "", nil, fmt.Errorf("input is not URL")
	}

	var res *http.Response
	var err error
	if cf.HTTPClient != nil {
		res, err = cf.HTTPClient.Get(url)
	} else {
		res, err = http.DefaultClient.Get(url)
	}
	if err != nil {
		return "", "", nil, err
	}
	defer res.Body.Close()
	if res.StatusCode < 400 {
		contentType := res.Header.Get("Content-Type")
		if strings.HasPrefix(contentType, "text/html") || strings.HasPrefix(contentType, "text/plain") {
			cbyte, err := ioutil.ReadAll(res.Body)
			if utf8.Valid(cbyte) {
				return url, contentType, cbyte, err
			}
			byter := bytes.NewReader(cbyte)
			dbyte, err := bufio.NewReader(byter).Peek(1024)
			if err != nil {
				panic(err)
			}
			e, name, _ := charset.DetermineEncoding(dbyte, contentType)
			nr := bytes.NewReader(cbyte)
			if e != nil {
				r := e.NewDecoder().Reader(nr)
				scb, err := ioutil.ReadAll(r)
				if err != nil {
					panic(err)
				}
				return url, contentType, scb, err
			}
			fmt.Printf("Bad Encode: %s", name)
			return url, "bad Encode", cbyte, fmt.Errorf("Bad Encode: %s", name)
		}

		return "", "", nil, fmt.Errorf("Invalid Content-Type: %s", contentType)
	}
	return "", "StatusError", []byte(res.Status), fmt.Errorf("statusError")
}
