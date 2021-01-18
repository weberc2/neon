package build

import (
	"bytes"
	"html/template"
	"io"
)

type PageRenderer struct {
	T      *template.Template
	Buffer bytes.Buffer
}

func (pr PageRenderer) RenderPage(w io.Writer, v interface{}) error {
	if err := pr.T.ExecuteTemplate(&pr.Buffer, "base", v); err != nil {
		return err
	}

	_, err := io.Copy(w, &pr.Buffer)
	return err
}
