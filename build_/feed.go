package build

import (
	"os"
	"path/filepath"
	"time"

	"bitbucket.org/weberc2/neon/config"
	"github.com/gorilla/feeds"
)

func buildFeed(conf config.Config, posts ByDate) error {
	var now time.Time
	if len(posts) > 0 {
		now = time.Time(posts[0].Date)
	} else {
		now = time.Now()
	}

	feed := &feeds.Feed{
		Title:       conf.Title,
		Link:        &feeds.Link{Href: conf.SiteRoot},
		Description: conf.Description,
		Author:      &feeds.Author{Name: conf.Author},
		Created:     now,
	}
	for _, post := range posts {
		feed.Items = append(
			feed.Items,
			&feeds.Item{
				Title:       post.Title,
				Link:        &feeds.Link{Href: relLink(conf.SiteRoot, post.ID)},
				Author:      &feeds.Author{Name: conf.Author},
				Created:     time.Time(post.Date),
				Description: string(snippet(post.Body)),
			},
		)
	}

	file, err := os.Create(filepath.Join(
		conf.OutputDirectory,
		"feed.atom",
	))
	if err != nil {
		return err
	}
	defer file.Close()

	return feed.WriteAtom(file)
}