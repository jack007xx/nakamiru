package nakamiru-kun

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"text/template"

	"github.com/TylerBrock/colorjson"
)

type req struct {
	To          string
	Method      string
	URL         string
	ContentType string
	Query       string
	Body        string
}

var logTemplateString = `
Nakamiru-kun found here!!

[Reqest]
Method :      {{.Method}}
To     :      {{.To}}
URL    :      {{.URL}}

Content-Type :    {{.ContentType}}

Query :
{{.Query}}

Body :
{{.Body}}
`

func Nakamiru(aHandler http.Handler) http.Handler {
	return http.HandlerFunc(func(aWriter http.ResponseWriter, aRequest *http.Request) {
		tBody, tError := ioutil.ReadAll(aRequest.Body)
		if tError != nil {
			print("[Nakamiru_kun panic!!]")
			panic(tError)
		}
		var tJSONBody map[string]interface{}
		json.Unmarshal(tBody, &tJSONBody)
		tBodyFormatter := colorjson.NewFormatter()
		tBodyFormatter.Indent = 2

		tPrettyBody, tError := tBodyFormatter.Marshal(tJSONBody)

		tReq := &req{
			Method:      aRequest.Method,
			To:          aRequest.Host + " [" + aRequest.RemoteAddr + "]",
			URL:         aRequest.URL.Path,
			ContentType: aRequest.Header.Get("Content-Type"),
			Body:        string(tPrettyBody),
			Query:       aRequest.URL.RawQuery,
		}

		tTemplate, tError := template.New("tReq").Parse(logTemplateString)
		if tError != nil {
			print("[Nakamiru_kun panic!!]")
			panic(tError)
		}

		tBuf := new(bytes.Buffer)
		tError = tTemplate.Execute(tBuf, tReq)
		if tError != nil {
			print("[Nakamiru_kun panic!!]")
			panic(tError)
		}

		log.Printf(tBuf.String())
		aHandler.ServeHTTP(aWriter, aRequest)
	})
}
