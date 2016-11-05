package linein

import (
	"errors"
	"net/http"
	"net/url"
	"strings"
)

type Linein struct {
	ResponseType string
	ClientId     string
	ClientSecret string
	GrantType    string
	Client       *http.Client
}

var (
	urlWeblogin    = "https://access.line.me/dialog/oauth/weblogin"
	urlAccessToken = "https://api.line.me/v1/oauth/accessToken"
)

func NewLinein(cId string, cSecret string) *Linein {
	l := &Linein{
		ResponseType: "code",
		ClientId:     cId,
		ClientSecret: cSecret,
		GrantType:    "authorization_code",
	}
	return l
}

func (l *Linein) GetWebLoginURL(redirectUri string, state string) (string, url.Values, error) {
	if redirectUri == "" {
		return "", nil, errors.New("redirect_uri is required")
	}
	values := url.Values{
		"response_type": {l.ResponseType},
		"client_id":     {l.ClientId},
		"redirect_uri":  {redirectUri},
		"state":         {state},
	}
	return urlWeblogin, values, nil
}

func (l *Linein) GetAccessTokenURL(code string, redirectUri string) (string, url.Values, error) {
	values := url.Values{
		"grant_type":    {l.GrantType},
		"client_id":     {l.ClientId},
		"client_secret": {l.ClientSecret},
		"code":          {code},
		"redirect_uri":  {redirectUri},
	}
	return urlAccessToken, values, nil
}

func Get(client *http.Client, urlStr string, form url.Values) (*http.Response, error) {
	req, err := http.NewRequest("GET", urlStr, nil)
	if err != nil {
		return nil, err
	}
	if req.URL.RawQuery != "" {
		return nil, errors.New("oauth: url must not contain a query string")
	}
	req.URL.RawQuery = form.Encode()
	if client == nil {
		client = http.DefaultClient
	}
	return client.Do(req)
}

func do(client *http.Client, method string, urlStr string, form url.Values) (*http.Response, error) {
	req, err := http.NewRequest(method, urlStr, strings.NewReader(form.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	if client == nil {
		client = http.DefaultClient
	}
	return client.Do(req)
}

func Post(client *http.Client, urlStr string, form url.Values) (*http.Response, error) {
	return do(client, "POST", urlStr, form)
}
