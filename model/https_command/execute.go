package https_command

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

var (
	contentType = map[BodyType]string{
		Text:               "text/plain",
		Json:               "application/json",
		Html:               "text/html",
		Xml:                "application/xml",
		FormData:           "multipart/form-data",
		XWwwFormUrlencoded: "application/x-www-form-urlencoded",
	}
)

func (hc *HttpsCommand) Execute() (result string) {
	var body io.Reader
	fmt.Println(hc)
	if hc.Body != nil {
		switch hc.BodyType {
		case Text, Json, Html, Xml:
			body = bytes.NewBuffer(hc.Body)
		case FormData, XWwwFormUrlencoded:
			jsonMap := make(map[string]string)
			if e := json.Unmarshal(hc.Body, &jsonMap); e != nil {
				httpsCommandLog.Error().Printf("id: %d body unmarshal failed", hc.CommandID)
			}
			param := url.Values{}
			for key, val := range jsonMap {
				param.Add(key, val)
			}
			body = strings.NewReader(param.Encode())
		}
	}
	req, _ := http.NewRequest(string(hc.Method), hc.Url, body)
	header := make(map[string]string)
	if hc.Header != nil {
		if e := json.Unmarshal(hc.Header, &header); e != nil {
			httpsCommandLog.Error().Printf("id: %d header unmarshal failed", hc.CommandID)
		}
	}
	for key, val := range header {
		req.Header.Add(key, val)
	}
	req.Header.Set("Content-Type", contentType[hc.BodyType])
	client := &http.Client{}
	var resp *http.Response
	if resp1, e := client.Do(req); e != nil {
		result = ""
		httpsCommandLog.Error().Printf("id: %d request failed", hc.CommandID)
		return
	} else {
		resp = resp1
	}
	var respBody []byte
	if respBody1, e := io.ReadAll(resp.Body); e != nil {
		result = ""
		httpsCommandLog.Error().Printf("id: %d request body failed", hc.CommandID)
		return
	} else {
		respBody = respBody1
	}
	result = string(respBody)
	defer func() {
		if e := resp.Body.Close(); e != nil {
			httpsCommandLog.Error().Println("Response body closed failed")
		}
	}()
	httpsCommandLog.Info().Printf("id: %d request status: %v\nrequest result: %s", hc.CommandID, resp.Status, result)
	return
}
