package fetch

import (
	"context"
	"errors"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/chromedp/cdproto/dom"
	"github.com/chromedp/chromedp"
)

func Fetch(uri string) (string, error) {
	// Craft the request for the page.
	req, _ := http.NewRequest("GET", uri, nil)
	req.Header.Set(
		"Accept",
		"text/html;q=0.9,*/*;q=0.8",
	)
	req.Header.Set(
		"User-Agent",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/15.0 Safari/605.1.15",
	)
	req.Header.Set("Accept-Language", "en-us")
	req.Header.Set("Connection", "keep-alive")

	// Make the request.
	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// Make sure content type is HTML
	cp := resp.Header.Get("Content-Type")
	if !strings.Contains(cp, "text/html") {
		return "", errors.New("URL is not a HTML document")
	}

	html, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(html), nil
}

func FetchWithChrome(uri string) (string, error) {
	ctx, cancel := chromedp.NewContext(
		context.Background(),
		// chromedp.WithLogf(log.Printf),
	)
	defer cancel()

	var html string
	err := chromedp.Run(ctx,
		chromedp.Navigate(uri),
		chromedp.ActionFunc(func(ctx context.Context) error {
			node, err := dom.GetDocument().Do(ctx)
			if err != nil {
				return err
			}
			html, err = dom.GetOuterHTML().WithNodeID(node.NodeID).Do(ctx)
			return err
		}),
	)
	if err != nil {
		return "", err
	}

	return html, nil
}
