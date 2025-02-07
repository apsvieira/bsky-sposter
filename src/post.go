package sposter

import (
	"bytes"
	"fmt"
	"strings"
	"text/template"

	"github.com/PuerkitoBio/goquery"
	"github.com/mmcdole/gofeed"
)

// GetFirstSentence returns the first sentence of the content of a post.
func GetFirstSentence(rawContent string, maxLen int) (string, error) {
	if maxLen < 1 {
		return "", fmt.Errorf("maxLen must be greater than 0")
	}

	// Parse the content HTML
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(rawContent))
	if err != nil {
		return "", fmt.Errorf("error parsing content: %v", err)
	}

	// Get the first sentence from the first paragraph element
	firstParagraph := doc.Find("p").First().Text()
	firstSentence := strings.Split(firstParagraph, ".")[0]

	if len(firstSentence) == 0 {
		return "", fmt.Errorf("no content found")
	}
	if len(firstSentence) < maxLen {
		return firstSentence, nil
	}

	return firstSentence[:maxLen-3] + "...", nil
}

type Post struct {
	Title         string
	Link          string
	PublishedDate string
	FirstSentence string
	LongIntro     string
}

func NewPostFromFeedItem(item *gofeed.Item) (*Post, error) {
	firstSentence, err := GetFirstSentence(item.Content, 200)
	if err != nil {
		return nil, fmt.Errorf("error getting first sentence: %v", err)
	}

	return &Post{
		Title:         item.Title,
		Link:          item.Link,
		PublishedDate: item.PublishedParsed.Format("2006-01-02"),
		FirstSentence: firstSentence,
		LongIntro:     item.Description,
	}, nil
}

func (p *Post) BskyPost() (string, error) {
	tmpl, err := template.New("bsky-post").Parse(bskyPostTemplate)
	if err != nil {
		return "", fmt.Errorf("error creating template: %v", err)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, p); err != nil {
		return "", fmt.Errorf("error executing template: %v", err)
	}

	return buf.String(), nil
}
