package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/dghubble/oauth1"
)

type Bookmark struct {
  Title string `json:"title"`
  Url   string `json:"url"`
  BookmarkId int `json:"bookmark_id"`
}

// LibraryData represents the structure of our saved library
type LibraryData struct {
	Bookmarks map[int]struct {
		Url   string `json:"url"`
		Title string `json:"title"`
	} `json:"bookmarks"`
}

// Retrieve Oauth token.
func get_access_token(consumer_key string, consumer_secret string, username string, password string) (string, string, error) {
	endpoint := "https://www.instapaper.com/api/1/oauth/access_token"

	// OAuth1 configuration
	config := oauth1.Config{
		ConsumerKey:    consumer_key,
		ConsumerSecret: consumer_secret,
		Endpoint: oauth1.Endpoint{
			RequestTokenURL: "",
			AuthorizeURL:    "",
			AccessTokenURL:  endpoint,
		},
	}

	token := oauth1.NewToken("", "") // No token initially
	client := config.Client(oauth1.NoContext, token)

	data := url.Values{}
	data.Set("x_auth_username", username)
	data.Set("x_auth_password", password)
	data.Set("x_auth_mode", "client_auth")

	req, err := http.NewRequest("POST", endpoint, strings.NewReader(data.Encode()))
	if err != nil {
		return "", "", err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

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

	values, err := url.ParseQuery(string(body))
	if err != nil {
		return "", "", err
	}

	return values.Get("oauth_token"), values.Get("oauth_token_secret"), nil
}

// Gets the archive folder, up to 500 articles.
// Currently checking with Instapaper how the limit works.
func get_bookmarks(access_token string, access_token_secret string) ([]Bookmark, error) {
  endpoint := "https://www.instapaper.com/api/1/bookmarks/list"
  
  config := oauth1.NewConfig(os.Getenv("INSTAPAPER_KEY"), os.Getenv("INSTAPAPER_SECRET"))
	token := oauth1.NewToken(access_token, access_token_secret)
	client := config.Client(oauth1.NoContext, token)

	params := url.Values{}
	params.Set("folder_id", "archive")
	params.Set("limit", "500")

	req, err := http.NewRequest("GET", endpoint+"?"+params.Encode(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var bookmarks []Bookmark
	if err := json.NewDecoder(resp.Body).Decode(&bookmarks); err != nil {
		return nil, fmt.Errorf("unable to parse JSON response: %v", err)
	}

	return bookmarks, nil
}

func initializeLibrary(libraryFilePath string) LibraryData {
	var library LibraryData
	
  // Check if file exists
	if _, err := os.Stat(libraryFilePath); err == nil {
		// File exists, read it
		data, err := os.ReadFile(libraryFilePath)
		if err != nil {
			panic(err)
		}
		
		if len(data) > 0 {
			if err := json.Unmarshal(data, &library); err != nil {
				// If error, try to initialize as empty
				library = LibraryData{
					Bookmarks: make(map[int]struct {
						Url   string `json:"url"`
						Title string `json:"title"`
					}),
				}
			}
		} else {
			// Empty file, initialize library
			library = LibraryData{
				Bookmarks: make(map[int]struct {
					Url   string `json:"url"`
					Title string `json:"title"`
				}),
			}
		}
	} else {
		// File doesn't exist, initialize library
		library = LibraryData{
			Bookmarks: make(map[int]struct {
				Url   string `json:"url"`
				Title string `json:"title"`
			}),
		}
	}
  return library
}

func main() {
	var consumer_key string = os.Getenv("INSTAPAPER_KEY")
	var consumer_secret string = os.Getenv("INSTAPAPER_SECRET")
	var username string = os.Getenv("INSTAPAPER_USERNAME")
	var password string = os.Getenv("INSTAPAPER_PASSWORD")

  if consumer_key == "" || consumer_secret == "" || username == "" || password == "" {
    fmt.Println("Missing required env vars.")
    os.Exit(1)
  }

	var access_token string
	var access_token_secret string
	var err error
	var bookmarks []Bookmark

	access_token, access_token_secret, err = get_access_token(consumer_key, consumer_secret, username, password)
	if err != nil {
		panic(err)
	}

	bookmarks, err = get_bookmarks(access_token, access_token_secret)
  if err != nil {
    fmt.Printf("Error retriving bookmarks %v\n", err)
    os.Exit(1)
  }

	// Read existing library file or create a new one
	libraryFilePath := "instapaper-library.json"
  library := initializeLibrary(libraryFilePath)

	// Open markdown file for appending
	markdownFile, err := os.OpenFile("instapaper-backup.md", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		panic(err)
	}
	defer markdownFile.Close()

  // TODO: Move this functionality into a function

	// Process new bookmarks
	newBookmarksCount := 0
	for _, bookmark := range bookmarks {
		// Check if bookmark already exists in library
		if _, exists := library.Bookmarks[bookmark.BookmarkId]; !exists {
			if bookmark.Title != "" {
				// Add to library
				library.Bookmarks[bookmark.BookmarkId] = struct {
					Url   string `json:"url"`
					Title string `json:"title"`
				}{
					Url:   bookmark.Url,
					Title: bookmark.Title,
				}
				
				// Write to markdown file
				fmt.Fprintf(markdownFile, "- %d [%s](%s)\n\n", bookmark.BookmarkId, bookmark.Title, bookmark.Url)
				newBookmarksCount++
			}
		}
	}

	// Write updated library back to file
	updatedLibraryData, err := json.MarshalIndent(library, "", "  ")
	if err != nil {
		panic(err)
	}
	
	if err := os.WriteFile(libraryFilePath, updatedLibraryData, 0644); err != nil {
		panic(err)
	}

	fmt.Printf("Added %d new bookmarks to library\n", newBookmarksCount)
	fmt.Printf("Total bookmarks in library: %d\n", len(library.Bookmarks))
}
