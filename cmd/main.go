package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/contrib/cors"
	"github.com/gin-gonic/gin"

	md "github.com/JohannesKaufmann/html-to-markdown"

	"github.com/go-shiori/go-readability"
	ginprometheus "github.com/zsais/go-gin-prometheus"

	"github.com/allokate-ai/scraper/internal/config"
	"github.com/allokate-ai/scraper/internal/logger"
	"github.com/allokate-ai/scraper/pkg/fetch"
	"github.com/allokate-ai/scraper/pkg/optional"
)

func Extract(html string, uri string) (Article, error) {
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
	markdown, err := converter.ConvertString(article.Content)
	if err != nil {
		log.Fatal(err)
	}

	return Article{
		URL:      uri,
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

type Article struct {
	URL      string `json:"url"`
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

func main() {
	l := logger.Get()
	defer l.Close()

	config, err := config.Get()
	if err != nil {
		l.Error(fmt.Sprint("invalid configuration:", err))
		log.Fatal(err)
	}

	fmt.Printf("Listening on port %d\n", config.Port)
	fmt.Println()

	router := gin.Default()

	p := ginprometheus.NewPrometheus("gin")
	p.Use(router)

	router.Use(cors.New(cors.Config{
		AllowedOrigins:   []string{"http://localhost:3000"},
		AllowedMethods:   []string{"PUT", "PATCH", "POST", "OPTIONS"},
		AllowedHeaders:   []string{"Origin", "Authorization"},
		ExposedHeaders:   []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	router.GET("/api/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	router.GET("/api/scrape", func(c *gin.Context) {
		var url *url.URL
		var useChrome bool

		errors := []string{}

		if p := c.Query("url"); p != "" {
			u, err := url.Parse(p)
			if err != nil {
				errors = append(errors, "url must be a valid URL")
			}

			url = u
		} else {
			errors = append(errors, "url is a required query parameter")
		}

		if p := c.Query("useChrome"); p != "" {
			boolValue, err := strconv.ParseBool(p)
			if err != nil {
				errors = append(errors, "useChrome must be a valid boolean value")
			}

			useChrome = boolValue
		}

		if len(errors) > 0 {
			c.JSON(http.StatusBadRequest, gin.H{
				"errors": errors,
			})
			return
		}

		fetcher := fetch.Fetch
		if useChrome {
			fetcher = fetch.FetchWithChrome
		}

		html, err := fetcher(url.String())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err,
			})
		}

		article, err := Extract(html, url.String())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err,
			})
		}

		c.JSON(http.StatusOK, article)
	})

	router.PUT("/api/scrape", func(c *gin.Context) {

		var body struct {
			Url      optional.Optional[string] `json:"url"`
			Document optional.Optional[string] `json:"document"`
		}

		errors := []string{}

		if err := c.ShouldBindJSON(&body); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"errors": "request must have valid JSON body",
			})
			return
		}

		if body.Url.Present() {
			_, err := url.Parse(body.Url.MustGet())
			if err != nil {
				errors = append(errors, "url must be a valid URL")
			}
		} else {
			errors = append(errors, "url is a required body parameter")
		}

		if !body.Document.Present() {
			errors = append(errors, "url is a required body parameter")
		}

		if len(errors) > 0 {
			c.JSON(http.StatusBadRequest, gin.H{
				"errors": errors,
			})
			return
		}

		article, err := Extract(body.Document.MustGet(), body.Url.MustGet())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err,
			})
		}

		c.JSON(http.StatusOK, article)
	})

	// Create a server and service incoming connections.
	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", config.Port),
		Handler: router,
	}
	go func() {
		server.ListenAndServe()
	}()

	// Wait for the signal to terminate.
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer func() {
		cancel()
	}()
	if err := server.Shutdown(ctx); err != nil {
		log.Fatal(err)
		l.Error("failed to start the server")
	}
	l.Info("received termination signal")
	log.Print("bye!")
}
