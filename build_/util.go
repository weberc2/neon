package build

import (
	"bytes"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/Depado/bfchroma"
	"gopkg.in/russross/blackfriday.v2"
)

func mkmdfunc(codeHighlightTheme string) func([]byte, string) []byte {
	return func(input []byte, footnoteURLPrefix string) []byte {
		return markdown(input, footnoteURLPrefix, codeHighlightTheme)
	}
}

func markdown(
	input []byte,
	footnoteURLPrefix string,
	codeHighlightTheme string,
) []byte {
	return blackfriday.Run(
		input,
		blackfriday.WithExtensions(
			blackfriday.CommonExtensions|blackfriday.Footnotes,
		),
		blackfriday.WithRenderer(
			bfchroma.NewRenderer(
				bfchroma.Extend(blackfriday.NewHTMLRenderer(
					blackfriday.HTMLRendererParameters{
						FootnoteAnchorPrefix: footnoteURLPrefix,
						Flags:                blackfriday.CommonHTMLFlags,
					},
				)),
				bfchroma.WithoutAutodetect(),
				bfchroma.Style(codeHighlightTheme),
			),
		),
	)
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
