package swimmy

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"strconv"
	"strings"
	"text/template"
	"time"

	"github.com/asaskevich/govalidator"
	"github.com/microcosm-cc/bluemonday"
	"golang.org/x/net/html"
)

//Sanitize sanitize html or txt with blueMonday
func Sanitize(htmlContent string, policy ...*bluemonday.Policy) string {
	if len(policy) < 1 {
		return DefaultPageDataBuilder.PreSanitizePolicy.Sanitize(htmlContent)
	}
	return policy[0].Sanitize(htmlContent)
}

//ParseTime parse time string
func ParseTime(timeStr string) (*time.Time, string, error) {
	formatStrings := []string{
		"2006-01-02T15:04:05-07:00",
		"2006-01-02T15:04:05Z",
		"20060102T150405-0700",
		"20060102T150405Z"}

	format := ""
	var resultTime *time.Time
	resultTime = nil
	var err error
	err = nil
	for _, fStr := range formatStrings {
		t, err := time.Parse(fStr, timeStr)
		if err == nil {
			format = fStr
			resultTime = &t
			break
		}
	}
	if resultTime == nil {
		err = fmt.Errorf("error: cannot parse format")
	}
	return resultTime, format, err
}

//TakeMarkedUpText is take marked-up text between begin tag and end tag
func TakeMarkedUpText(ct *html.Tokenizer, tagName []byte) string {
	depth := 0
	taking := true
	var sb *strings.Builder

	for taking {
		tt := ct.Next()
		switch tt {
		case html.StartTagToken:
			tName, _ := ct.TagName()
			if bytes.Equal(tName, tagName) {
				depth++
			}
			WriteCurrentString(ct, tt, sb)

		case html.EndTagToken:
			tName, _ := ct.TagName()
			if bytes.Equal(tName, tagName) {
				depth--
				if depth < 1 {
					taking = false
				} else {
					WriteCurrentString(ct, tt, sb)
				}
			} else {
				WriteCurrentString(ct, tt, sb)
			}
		case html.ErrorToken:
			taking = false
		default:
			WriteCurrentString(ct, tt, sb)
		}
	}
	return sb.String()
}

//WriteCurrentString write string of now tag or text to strings.Builder
func WriteCurrentString(tokenizer *html.Tokenizer, tokenType html.TokenType, sb *strings.Builder) {
	switch tokenType {
	case html.ErrorToken:
	case html.TextToken:
		sb.WriteString(EscapeBytes(tokenizer.Text()))
	case html.StartTagToken:
		sb.WriteString("<")
		writeCurrentTagString(tokenizer, sb)
		sb.WriteString(">")
	case html.EndTagToken:
		sb.WriteString("</")
		tName, _ := tokenizer.TagName()
		sb.Write(tName)
		sb.WriteString(">")
	case html.SelfClosingTagToken:
		sb.WriteString("<")
		writeCurrentTagString(tokenizer, sb)
		sb.WriteString("/>")
	case html.CommentToken:
		sb.WriteString("<!--")
		sb.WriteString(EscapeBytes(tokenizer.Text()))
		sb.WriteString("-->")
	case html.DoctypeToken:
		sb.WriteString("<!DOCTYPE ")
		sb.WriteString(EscapeBytes(tokenizer.Text()))
		sb.WriteString(">")
	default:
		sb.WriteString("Invalid token:< ")
		sb.WriteString(strconv.Itoa(int(tokenType)))
		sb.WriteString(">")
	}
}

func writeCurrentTagString(tokenizer *html.Tokenizer, sb *strings.Builder) {
	tName, hasAttr := tokenizer.TagName()
	if hasAttr {
		sb.Write(tName)

		moreAttr := true
		var key, val []byte
		for moreAttr {
			key, val, moreAttr = tokenizer.TagAttr()

			sb.WriteByte(' ')
			sb.WriteString(EscapeBytes(key))
			sb.WriteString(`="`)
			sb.WriteString(EscapeBytes(val))
			sb.WriteByte('"')
		}
	} else {
		sb.Write(tName)
	}
}

//EscapeBytes html escape byte array
func EscapeBytes(str []byte) string {
	return html.EscapeString(string(str))
}

//CommentifyString commentify string inspired by cobra's commentifyString
func CommentifyString(input string) string {
	nlcode := "\n"
	replacer := strings.NewReplacer("\r\n", nlcode, "\r", nlcode, "\n", nlcode)
	inputNLd := replacer.Replace(input)

	lines := strings.Split(inputNLd, "\n")
	var sb strings.Builder
	sb.Grow(len(input) + len(lines)*len("\n"))
	c := "//"
	for _, l := range lines {
		if strings.HasPrefix(l, c) {
			sb.WriteString(l)
			sb.WriteString(nlcode)
		} else {
			sb.WriteString(c)
			if l != "" {
				sb.WriteString(l)
			}
			sb.WriteString(nlcode)
		}
	}

	return strings.TrimSuffix(sb.String(), nlcode)
}

//ExecLicenseTextTemp exec template using templateStr and data
func ExecLicenseTextTemp(templateStr string, data interface{}) (string, error) {
	fm := template.FuncMap{"comment": CommentifyString}
	t, err := template.New("").Funcs(fm).Parse(templateStr)

	if err != nil {
		return "", err
	}

	bb := bytes.NewBuffer(make([]byte, 0))

	err = t.Execute(bb, data)

	return bb.String(), err
}

//ReadList read list with newline-delimited.
func ReadList(listPath string) ([]string, error) {
	f, err := os.Open(listPath)
	defer f.Close()
	if err != nil {
		return nil, err
	}

	scanner := bufio.NewScanner(f)

	list := make([]string, 0)
	for scanner.Scan() {
		list = append(list, scanner.Text())
	}
	return list, scanner.Err()
}

//IsFilePath validate whether val is filepath or not and confirm that it exist and it is not directory.
func IsFilePath(val string) (bool, error) {
	i, _ := govalidator.IsFilePath(val)
	if i {

	}

}
