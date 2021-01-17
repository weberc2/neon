package build

import (
	"os"
	"path/filepath"
)

type Page struct {
	Content interface{}
	NextID  string
	PrevID  string
}

type Slice interface {
	Len() int
	At(i int) (string, interface{})
}

type PageBuilder struct {
	ThemeConfig     interface{}
	Renderer        PageRenderer
	OutputDirectory string
}

func (pb PageBuilder) BuildPage(id string, page Page) error {
	path := filepath.Join(pb.OutputDirectory, id)
	if err := os.MkdirAll(filepath.Dir(path), 0777); err != nil {
		return err
	}

	file, err := os.Create(path)
	if err != nil {
		return err
	}

	return pb.Renderer.RenderPage(
		file,
		struct {
			ThemeConfig interface{}
			Page        Page
		}{
			ThemeConfig: pb.ThemeConfig,
			Page:        page,
		},
	)
}

func (pb PageBuilder) BuildPages(items Slice) error {
	if items.Len() < 1 {
		return nil
	}

	prevID, content := items.At(0)
	page := Page{Content: content}
	for i := 1; i < items.Len(); i++ {
		// prepare the page for writing
		id, content := items.At(i)
		page.NextID = id
		if err := pb.BuildPage(prevID, page); err != nil {
			return err
		}

		// update the page for the next loop
		page.Content = content
		page.PrevID = prevID
		prevID = id
		page.NextID = ""
	}

	return pb.BuildPage(prevID, page)
}
