package scrapper

import (
	"fmt"
	"log"
	"net/url"
	"strings"
	"time"

	md "github.com/JohannesKaufmann/html-to-markdown"

	"github.com/go-shiori/go-readability"
)

type Article struct {
	Url      string `json:"url"`
	Title    string `json:"title"`
	Byline   string `json:"byline"`
	Length   int    `json:"length"`
	Excerpt  string `json:"excerpt"`
	SiteName string `json:"siteName"`
	Image    string `json:"image"`
	Favicon  string `json:"favicon"`
	Content  string `json:"content"`
	Markdown string `json:"markdown"`
	Fetched  string `json:"fetched"`
}

func Scrape(html string, uri string) (Article, error) {
	u, err := url.ParseRequestURI(uri)
	if err != nil {
		log.Fatal(err)
	}

	// Parse with readability library.
	parser := readability.NewParser()
	article, err := parser.Parse(strings.NewReader(html), u)
	if err != nil {
		log.Fatalf("failed to parse %s, %v\n", uri, err)
	}

	// Convert to markdown.
	converter := md.NewConverter("", true, nil)
	byline := ""
	if article.Byline != "" {
		byline = fmt.Sprintf("<strong>%s</strong>", article.Byline)
	}
	markdown, err := converter.ConvertString(
		fmt.Sprintf("<body><h1>%s</h1>%s%s</body>", article.Title, byline, article.Content),
	)

	if err != nil {
		log.Fatal(err)
	}

	return Article{
		Url:      uri,
		Title:    article.Title,
		Byline:   article.Byline,
		Length:   article.Length,
		Excerpt:  article.Excerpt,
		SiteName: article.SiteName,
		Image:    article.Image,
		Favicon:  article.Favicon,
		Content:  article.Content,
		Markdown: markdown,
		Fetched:  time.Now().Format(time.RFC3339),
	}, nil
}
