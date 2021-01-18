package build

import (
	"html/template"
	"path/filepath"
)

func loadTemplates(
	dir string,
	funcs template.FuncMap,
) (post, index *template.Template, err error) {
	post, err = template.New("post").Funcs(funcs).ParseFiles(
		filepath.Join(dir, "base.html"),
		filepath.Join(dir, "post.html"),
	)
	if err != nil {
		return nil, nil, err
	}

	index, err = template.New("index").Funcs(funcs).ParseFiles(
		filepath.Join(dir, "base.html"),
		filepath.Join(dir, "index.html"),
	)
	if err != nil {
		return nil, nil, err
	}

	return post, index, nil
}
