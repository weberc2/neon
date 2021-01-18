package build

import (
	"io/ioutil"
	"os"
	"path/filepath"
)

type FilepathIter struct {
	filenames   []string
	directories []string
	err         error
	next        string
}

func (fpi *FilepathIter) Next() bool {
	// if we have more files in our queue, load in the next one
	if len(fpi.filenames) > 0 {
		fpi.next = fpi.filenames[0]
		fpi.filenames = fpi.filenames[1:]
		return true
	}

	// if we have more dictionaries on our stack, load in the next one
	if len(fpi.directories) > 0 {
		var files []os.FileInfo
		files, fpi.err = ioutil.ReadDir(fpi.directories[0])
		if fpi.err != nil {
			return false
		}
		for _, file := range files {
			path := filepath.Join(fpi.directories[0], file.Name())
			if file.IsDir() {
				fpi.directories = append(fpi.directories, path)
				continue
			}
			fpi.filenames = append(fpi.filenames, path)
		}
		fpi.directories = fpi.directories[1:]
		return fpi.Next()
	}

	// if we have no more files nor dictionaries, return false
	return false
}

func (fpi *FilepathIter) Filepath() string {
	return fpi.next
}

func (fpi *FilepathIter) Err() error {
	return fpi.err
}
