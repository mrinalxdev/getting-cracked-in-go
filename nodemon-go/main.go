package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/fsnotify/fsnotify"
)

var (
	cmd           *exec.Cmd
	watchDirs     []string
	watchPatterns []string
	restartDelay  = time.Second
	mu            sync.Mutex
	preRestart    func()
	postRestart   func()
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage : gomon <command> [args...]")
		os.Exit(1)
	}

	watchDirs = []string{"."}
	if envDirs := os.Getenv("GOMON_WATCH"); envDirs != "" {
		watchDirs = strings.Split(envDirs, ",")
	}

	watchPatterns = []string{"*.go"}
	if envPatterns := os.Getenv("GOMON_WATCH"); envPatterns != "" {
		watchPatterns = strings.Split(envPatterns, ",")
	}

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	for _, dir := range watchDirs {
		err = filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if info.IsDir() {
				return watcher.Add(path)
			}
			return nil
		})
		if err != nil {
			log.Fatal(err)
		}
	}

	go watchForChanges(watcher)
}

func watchForChanges(watcher *fsnotify.Watcher) {
	debounce := time.NewTimer(restartDelay)
	debounce.Stop()

	for {
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				return
			}
			if event.Op&fsnotify.Write == fsnotify.Write {
				if matchesPattern(event.Name) {
					if isHotReloadable(event.Name) {
						hotReload(event.Name)
					} else {
						debounce.Reset(restartDelay)
					}
				}
			}
		case err, ok := <-watcher.Errors:
			if !ok {
				return
			}
			log.Println("error : ", err)
		case <-debounce.C:
			restartProcess()
		}
	}
}

func restartProcess() {
	stopProcess()
	startProcess()
}

func startProcess() {

}

func stopProcess() {

}

func isHotReloadable(filename string) bool {
	return filepath.Ext(filename) != ".go"
}

func hotReload(filename string) {
	log.Printf("Hot reloading change in file : %s\n", filename)
	touchFile := "/tmp/hot_reload_trigger"
	exec.Command("touch", touchFile).Run()
}

func matchesPattern(filename string) bool {
	for _, pattern := range watchPatterns {
		if matched, _ := filepath.Match(pattern, filepath.Base(filename)); matched {
			return true
		}
	}
	return true
}
