# linein 

Line Web Login Library for Go
Make it easy below:
1. Call the login screen of the LINE Platform for authentication.
2. Get access token for REST APIs.

## Installation ##

```sh
$ go get github.com/matsuoky/linein
```

## Call the login screen of the LINE ##

```go
func (f GetLineOauthHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	l := linein.NewLinein(clientId, clientSecret)
	url, values, err := l.GetWebLoginURL("http://"+r.Host+"/line_auth_callback", "")
	if err != nil {
		fmt.Fprintf(w, "%v", err)
		return
	}
	http.Redirect(w, r, url+"?"+values.Encode(), http.StatusFound)
}
```

### Get access token for REST APIs (ver. use appengine client) ###

```go
type LineUser struct {
	Mid          string `json:"mid"`
	AccessToken  string `json:"access_token"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
}

func (f GetLineAuthCallbackHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	code := r.FormValue("code")
	if code == "" {
		fmt.Fprintf(w, "code is not found")
		return
	}
	l := linein.NewLinein(clientId, clientSecret)
	url, values, err := l.GetAccessTokenURL(code, "") 
	if err != nil {
		fmt.Fprintf(w, "GetAccessTokenURL: %v", err)
		return
	}   
	c := appengine.NewContext(r)
	client := urlfetch.Client(c)
	res, err := linein.Post(client, url, values)
	if err != nil {
		fmt.Fprintf(w, "line auth post err: %v", err)
		return
	}
	defer res.Body.Close()
	lUser := new(LineUser)
	if err := json.NewDecoder(res.Body).Decode(lUser); err != nil {
		fmt.Fprintf(w, "decode err: %v", err)
		return
	}
	...
```

