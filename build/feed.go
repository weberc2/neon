package build

import (
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/gorilla/feeds"
	"github.com/weberc2/neon/config"
)

type markdownRenderer func(input []byte, footnoteURLPrefix string) []byte

func buildFeed(
	conf config.Config,
	posts ByDate,
	renderMarkdown markdownRenderer,
) error {
	var now time.Time
	if len(posts) > 0 {
		now = time.Time(posts[0].Date)
	} else {
		now = time.Now()
	}

	feed := &feeds.Feed{
		Title:       conf.Feed.Title,
		Link:        &feeds.Link{Href: conf.SiteRoot},
		Description: conf.Feed.Description,
		Author:      &feeds.Author{Name: conf.Feed.Author},
		Created:     now,
	}
	for _, post := range posts {
		body := renderMarkdown(post.Body, post.ID)
		description := snippet(body)
		if len(description) < 1 {
			description = body
		}
		feed.Items = append(
			feed.Items,
			&feeds.Item{
				Title:       post.Title,
				Link:        &feeds.Link{Href: relLink(conf.SiteRoot, post.ID)},
				Author:      &feeds.Author{Name: conf.Feed.Author},
				Created:     time.Time(post.Date),
				Description: string(description),
			},
		)
	}

	file, err := os.Create(filepath.Join(conf.OutputDirectory, "feed.atom"))
	if err != nil {
		return err
	}
	defer func() {
		if err := file.Close(); err != nil {
			log.Printf("ERROR Failed to close file: %v", err)
		}
	}()

	return feed.WriteAtom(file)
}
