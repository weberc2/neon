package build

import (
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"strings"

	"bitbucket.org/weberc2/neon/config"
)

func postsDirectory(conf config.Config) string {
	return filepath.Join(conf.InputDirectory, "posts")
}

func themeDirectory(conf config.Config, innerPaths ...string) string {
	return filepath.Join(append([]string{
		conf.InputDirectory,
		"themes",
		conf.Theme,
	}, innerPaths...)...)
}

func copyResources(conf config.Config) error {
	outputDirectory := filepath.Join(conf.OutputDirectory, "resources")
	inputDirectory := themeDirectory(conf, "resources")
	if err := copyDirectory(outputDirectory, inputDirectory); err != nil {
		if !os.IsNotExist(err) {
			return err
		}
	}
	return nil
}

func copyStaticDir(conf config.Config) error {
	outputDirectory := filepath.Join(conf.OutputDirectory, "static")
	inputDirectory := filepath.Join(conf.InputDirectory, "static")
	if err := copyDirectory(outputDirectory, inputDirectory); err != nil {
		if !os.IsNotExist(err) {
			return err
		}
	}
	return nil
}

func BuildSite(conf config.Config) error {
	funcs := template.FuncMap{
		"snippet":  snippet,
		"markdown": mkmdfunc(conf.CodeHighlightTheme),
		"html_":    func(b []byte) template.HTML { return template.HTML(b) },
		"rellink": func(urlpath string) string {
			root := strings.TrimRight(conf.SiteRoot, "/")
			urlpath = strings.TrimLeft(urlpath, "/")
			return root + "/" + urlpath
		},
	}
	templatesDirectory := themeDirectory(conf, "templates")
	postTemplate, indexTemplate, err := loadTemplates(templatesDirectory, funcs)
	if err != nil {
		return buildErrors.New(
			"template_loading",
			"Failed to load template",
			err,
		)
	}

	// index posts
	posts, err := indexPosts(
		postsDirectory(conf),
		conf.PostOutputDirectory,
		conf.OutputExtension,
	)
	if err != nil {
		return buildErrors.New("post_indexing", "Failed to index posts", err)
	}

	// build post pages
	pb := PageBuilder{
		OutputDirectory: conf.OutputDirectory,
		ThemeConfig:     conf.ThemeConfig,
		Renderer:        PageRenderer{T: postTemplate},
	}
	if err := pb.BuildPages(posts); err != nil {
		return buildErrors.New(
			"page_building",
			"Failed to build post pages",
			err,
		)
	}

	// build index pages
	pb.Renderer.T = indexTemplate
	indexPages := IndexPages{
		Posts:    posts,
		PageSize: conf.IndexPageSize,
		IDFunc: func(i int) string {
			if i < 1 {
				return "index.html"
			}
			return fmt.Sprintf("%d.html", i)
		},
	}
	if err := pb.BuildPages(indexPages); err != nil {
		return buildErrors.New(
			"page_building",
			"Failed to build index pages",
			err,
		)
	}

	// copy over the resources directory
	if err := copyResources(conf); err != nil {
		return buildErrors.New(
			"page_building",
			"Failed to copy theme resources",
			err,
		)
	}

	// copy over the static assets directory
	if err := copyStaticDir(conf); err != nil {
		return buildErrors.New(
			"page_building",
			"Failed to copy the /static directory",
		)
	}

	return nil
}
