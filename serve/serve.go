package serve

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/jaschaephraim/lrserver"

	"gopkg.in/fsnotify.v1"

	"bitbucket.org/weberc2/neon/build_"
	"bitbucket.org/weberc2/neon/config"
)

func serve(liveReload *lrserver.Server, dir, port string) error {
	server := http.FileServer(http.Dir(dir))
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.Method, r.URL)
		server.ServeHTTP(w, r)
	})
	fmt.Println("Listening on port", port)
	return http.ListenAndServe(":"+port, nil)
}

func addWatchDirs(w *fsnotify.Watcher, inputDir, outputDir string) error {
	outputDirInfo, err := os.Stat(outputDir)
	if err != nil {
		return err
	}

	directories := []string{inputDir}
	w.Add(inputDir)
	for len(directories) > 0 {
		// pop off the current directory
		currentDirectory := directories[0]
		directories = directories[1:]

		files, err := ioutil.ReadDir(currentDirectory)
		if err != nil {
			return err
		}

		for _, file := range files {
			// if the file is a directory, add it to the watcher and push it on
			// the stack to be scanned for subdirectories (unless it's the
			// output directory--we want to skip him)
			if file.IsDir() {
				// ignore the output directory
				if os.SameFile(file, outputDirInfo) {
					continue
				}
				directory := filepath.Join(currentDirectory, file.Name())
				w.Add(directory)
				directories = append(directories, directory)
			}
		}
	}
	return nil
}

func watch(conf config.Config, liveReload *lrserver.Server) error {
	w, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}
	defer w.Close()

	if err := addWatchDirs(
		w,
		conf.InputDirectory,
		conf.OutputDirectory,
	); err != nil {
		return err
	}

	// poll every 2 seconds to see if we need to rebuild
	ignoreEventsDelay := 2000 * time.Millisecond
	timer := time.NewTimer(ignoreEventsDelay)
	ignoringEvents := false
	for {
		select {
		case event, ok := <-w.Events:
			if !ok {
				return nil
			}
			log.Println("File changed:", event.Name)
			if !ignoringEvents {
				ignoringEvents = true
				log.Println("Rebuilding...")
				if err := build.BuildSite(conf); err != nil {
					return err
				}
				liveReload.Reload(event.Name)
				log.Println("Rebuilding finished")
				timer.Reset(ignoreEventsDelay)
			}
		case err := <-w.Errors:
			log.Println("FSNOTIFY ERR:", err)
		case <-timer.C:
			ignoringEvents = false
		}
	}
}

func parallel(fs ...func() error) error {
	errs := make(chan error)
	for _, f := range fs {
		go func(f func() error) {
			errs <- f()
		}(f)
	}

	for range fs {
		if err := <-errs; err != nil {
			return err
		}
	}

	return nil
}

func Serve(port uint16) error {
	conf, err := config.Load()
	if err != nil {
		return err
	}

	conf.SiteRoot = "http://localhost:" + strconv.Itoa(int(port))

	if err := build.BuildSite(conf); err != nil {
		return err
	}

	liveReload := lrserver.New("live-reload-server", port)

	return parallel(
		func() error {
			return watch(conf, liveReload)
		},
		func() error {
			return serve(
				liveReload,
				conf.OutputDirectory,
				strconv.Itoa(int(port)),
			)
		},
	)
}
