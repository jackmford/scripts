package main

import (
  "fmt"
  "os"
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

func get_access_token () (string, string) {
  return "", ""
}

func get_bookmarks(access_token string, access_token_secret string) []Bookmark {
  return nil
}

func main() {
  fmt.Println("Hi")

  var access_token string
  var access_token_secret string
  var bookmarks []Bookmark

  access_token, access_token_secret = get_access_token() 
  bookmarks = get_bookmarks(access_token, access_token_secret)

  for _, bookmark := range bookmarks {
    fmt.Println(bookmark.Title)
  }

  fmt.Println(CONSUMER_KEY)
}
