package main

import (
  "fmt"
  "io/ioutil"
  "net/http"
  "net/url"
  "os"
  "strings"

  "github.com/dghubble/oauth1"
)

var CONSUMER_KEY string = os.Getenv("INSTAPAPER_KEY")
var CONSUMER_SECRET string = os.Getenv("INSTAPAPER_SECRET")
var USERNAME string = os.Getenv("INSTAPAPER_USERNAME")
var PASSWORD string = os.Getenv("INSTAPAPER_PASSWORD")

type Bookmark struct {
  Title string
  Url string
  Index int
}

func get_access_token() (string, string, error) {
	endpoint := "https://www.instapaper.com/api/1/oauth/access_token"

	// OAuth1 configuration
	config := oauth1.Config{
		ConsumerKey:    CONSUMER_KEY,
		ConsumerSecret: CONSUMER_SECRET,
		Endpoint: oauth1.Endpoint{
			RequestTokenURL: "",
			AuthorizeURL:    "",
			AccessTokenURL:  endpoint,
		},
	}

	token := oauth1.NewToken("", "") // No token initially
	client := config.Client(oauth1.NoContext, token)

	// Form data
	data := url.Values{}
	data.Set("x_auth_username", USERNAME)
	data.Set("x_auth_password", PASSWORD)
	data.Set("x_auth_mode", "client_auth")

	// Create request
	req, err := http.NewRequest("POST", endpoint, strings.NewReader(data.Encode()))
	if err != nil {
		return "", "", err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// Execute request
	resp, err := client.Do(req)
	if err != nil {
		return "", "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", "", fmt.Errorf("failed to obtain access token: %s", resp.Status)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", "", err
	}

	// Parse response
	values, err := url.ParseQuery(string(body))
	if err != nil {
		return "", "", err
	}

  fmt.Println(values.Get("oauth_token"))
	return values.Get("oauth_token"), values.Get("oauth_token_secret"), nil
}

func get_bookmarks(access_token string, access_token_secret string) []Bookmark {
  return nil
}

func main() {
  fmt.Println("Hi")

  var access_token string
  var access_token_secret string
  var err error
  var bookmarks []Bookmark

  access_token, access_token_secret, err = get_access_token() 
  if err != nil {
    panic(err)
  }
  bookmarks = get_bookmarks(access_token, access_token_secret)

  for _, bookmark := range bookmarks {
    fmt.Printf("%d: %s %s", bookmark.Index, bookmark.Title, bookmark.Url)
  }
}
