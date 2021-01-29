package build

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/Depado/bfchroma"
	"github.com/russross/blackfriday/v2"
)

type markdownFunc func(
	input []byte,
	footnoteURLPrefix string,
	headingLevelOffset int,
) []byte

func mkmdfunc(
	codeHighlightTheme string,
	linkPrefix string,
	outputExtension string,
) markdownFunc {
	return func(
		input []byte,
		footnoteURLPrefix string,
		headingLevelOffset int,
	) []byte {
		return markdown(
			input,
			footnoteURLPrefix,
			codeHighlightTheme,
			linkPrefix,
			outputExtension,
			headingLevelOffset,
		)
	}
}

func markdown(
	input []byte,
	footnoteURLPrefix string,
	codeHighlightTheme string,
	linkPrefix string,
	outputExtension string,
	headingLevelOffset int,
) []byte {
	return blackfriday.Run(
		input,
		blackfriday.WithExtensions(
			blackfriday.CommonExtensions|
				blackfriday.Footnotes|
				blackfriday.Tables|
				blackfriday.AutoHeadingIDs,
		),
		blackfriday.WithRenderer(
			&renderer{
				linkPrefix:      linkPrefix,
				outputExtension: outputExtension,
				Renderer: bfchroma.NewRenderer(
					bfchroma.Extend(blackfriday.NewHTMLRenderer(
						blackfriday.HTMLRendererParameters{
							FootnoteAnchorPrefix: footnoteURLPrefix,
							Flags:                blackfriday.CommonHTMLFlags,
							HeadingLevelOffset:   headingLevelOffset,
						},
					)),
					bfchroma.WithoutAutodetect(),
					bfchroma.Style(codeHighlightTheme),
				),
			},
		),
	)
}

// This renderer handles relative links to posts. First of all, it allows the
// author to link to source files so the source markdown doesn't need to know
// (and thus be tightly coupled to) configuration details like site root, post
// output directory, and post output extension. More importantly, it prevents
// links in snippet texts from failing to resolve on index pages and feed
// descriptions. E.g., if a post 'foo.md' links to './bar.md' and the post
// output directory is /posts/, then the snippet for the 'foo' post will link
// to ./bar.html instead of ./posts/bar.html. To work around this, we make all
// links beginning with './' into absolute links (we also want to take care to
// handle cases where the site root is a nested path, like
// 'example.com/blog/').
type renderer struct {
	blackfriday.Renderer
	linkPrefix      string
	outputExtension string
}

func (r *renderer) RenderNode(
	w io.Writer,
	node *blackfriday.Node,
	entering bool,
) blackfriday.WalkStatus {
	const prefix = "./"
	n := *node
	if bytes.HasPrefix(n.LinkData.Destination, []byte(prefix)) {
		n.LinkData.Destination = []byte(fmt.Sprintf(
			"%s/%s",
			r.linkPrefix,
			bytes.ReplaceAll(
				n.LinkData.Destination[len(prefix):],
				[]byte(postInputExtension),
				[]byte(r.outputExtension),
			),
		))
	}
	return r.Renderer.RenderNode(w, &n, entering)
}

func snippet(input []byte) []byte {
	idx := bytes.Index(input, []byte("<!-- more -->"))
	if idx < 0 {
		return nil
	}
	return input[:idx]
}

func copyFile(dst, src string) error {
	srcf, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcf.Close()

	dstf, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer dstf.Close()

	_, err = io.Copy(dstf, srcf)
	return err
}

func copyDirectory(dst, src string) error {
	files, err := ioutil.ReadDir(src)
	if err != nil {
		return err
	}

	if err := os.Mkdir(dst, 0777); err != nil {
		if !os.IsExist(err) {
			return err
		}
	}

	for _, file := range files {
		frompath := filepath.Join(src, file.Name())
		topath := filepath.Join(dst, file.Name())
		if file.IsDir() {
			if err := copyDirectory(topath, frompath); err != nil {
				return err
			}
		} else {
			if err := copyFile(topath, frompath); err != nil {
				return err
			}
		}
	}

	return nil
}
