package build

import (
	"os"
	"path/filepath"
	"sort"
	"strings"
)

type Post struct {
	ID    string
	Date  Time
	Title string
	Body  []byte
}

type ByDate []Post

func (posts ByDate) Swap(i, j int) {
	temp := posts[i]
	posts[i] = posts[j]
	posts[j] = temp
}

func (posts ByDate) Less(i, j int) bool {
	return posts[i].Date.After(posts[j].Date)
}

func (posts ByDate) Len() int {
	return len(posts)
}

func (posts ByDate) At(i int) (string, interface{}) {
	return posts[i].ID, posts[i]
}

func parseFile(path string) (Post, error) {
	file, err := os.Open(path)
	if err != nil {
		return Post{}, err
	}
	defer file.Close()

	return parsePost(file)
}

const postInputExtension = ".md"

func indexPosts(dir, outputDir, outputExtension string) (ByDate, error) {
	fpi := &FilepathIter{directories: []string{dir}}

	posts := ByDate{}
	for fpi.Next() {
		path := fpi.Filepath()
		if !strings.HasSuffix(path, postInputExtension) {
			continue
		}
		relpath, err := filepath.Rel(dir, path)
		if err != nil {
			// we shouldn't get an error since all paths returned by the
			// FilePathIter should be valid relative paths
			panic(err)
		}

		post, err := parseFile(path)
		if err != nil {
			return nil, err
		}

		idx := strings.LastIndex(relpath, ".")
		post.ID = filepath.Join(
			outputDir,
			relpath[:idx]+outputExtension,
		)

		posts = append(posts, post)
	}
	if err := fpi.Err(); err != nil {
		return nil, err
	}

	sort.Sort(posts)
	return posts, nil
}
