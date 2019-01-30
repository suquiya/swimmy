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

//FetchError is error struct for fetch
type FetchError struct {
	ErrorType int
	s         string
}

//Error is implement of fetcherror for error interface
func (fe *FetchError) Error() string {
	return fe.s
}

//NewFetchError create fetch error
func NewFetchError(t int, s string) *FetchError {
	return &FetchError{t, s}
}

const (
	//IsNotURLError represents error occur from input is not URL
	IsNotURLError = iota
	//BadEncodeError represents url contents is not encoded with encode methods that swimmy can handle
	BadEncodeError
	//InvalidContentTypeError represents content is not html or text
	InvalidContentTypeError
	//StatusError is StatusError
	StatusError
)

//Fetch fetch url contents with f's HTTPClient. If you don't set your custom http client, Fetch use DefaultClient of net/http package.
func (cf *ContentFetcher) Fetch(url string) (string, string, []byte, error) {

	if !govalidator.IsRequestURL(url) {
		return "", "", nil, NewFetchError(IsNotURLError, "Input is not url")
	}

	var res *http.Response
	var err error
	if cf.HTTPClient != nil {
		res, err = cf.HTTPClient.Get(url)
	} else {
		res, err = http.DefaultClient.Get(url)
	}
	defer res.Body.Close()
	if err != nil {
		return "", "", nil, err
	}
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
			return url, "bad Encode", cbyte, NewFetchError(BadEncodeError, fmt.Sprintf("Bad Encode: %s", name))
		}

		s := fmt.Sprintf("Invalid Content-Type: %s", contentType)
		return "", "", nil, NewFetchError(InvalidContentTypeError, s)
	}
	return "", "StatusError", []byte(res.Status), NewFetchError(StatusError, "StatusError")
}
