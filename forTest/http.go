package forTest

import (
	"bytes"
	"io"
	"net/http"
	"net/url"
	"strings"
)

func httpTest() {
	_, _ = http.NewRequest("GET", "http://api.themoviedb.org/3/tv/popular", nil)
	parm := url.Values{}
	parm.Add("client_id", "ayeshaj")
	parm.Add("response_type", "code")
	parm.Add("scope", "public_profile")
	parm.Add("redirect_uri", "http://ayeshaj:8080/playground")
	parm.Encode()
	strings.NewReader(parm.Encode())
	bytes.NewBuffer([]byte(`{a:aa}`))
	_, _ = io.ReadAll(bytes.NewBuffer([]byte(`{a:aa}`)))

}
