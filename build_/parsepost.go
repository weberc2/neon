package build

import (
	"bufio"
	"io"

	"github.com/ghodss/yaml"
)

func parsePost(r io.Reader) (Post, error) {
	scanner := bufio.NewScanner(r)
	var post Post

	// scan in first line
	scanner.Scan()

	// if the first line is the metadata delim, grab all the bytes until the
	// next metadata delim
	if scanner.Text() == "---" {
		metadataBytes := []byte{}
		for {
			if !scanner.Scan() {
				// TODO: Unexpected error or EOF
			}
			if scanner.Text() == "---" {
				break
			}
			metadataBytes = append(metadataBytes, scanner.Bytes()...)
			metadataBytes = append(metadataBytes, '\n')
		}
		if err := yaml.Unmarshal(metadataBytes, &post); err != nil {
			return post, err
		}
	} else { // otherwise regard the first line as body
		post.Body = scanner.Bytes()
	}

	for scanner.Scan() {
		post.Body = append(post.Body, scanner.Bytes()...)
		post.Body = append(post.Body, '\n')
	}

	if err := scanner.Err(); err != nil {
		return post, err
	}

	return post, nil
}
