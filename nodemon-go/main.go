package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"sync"
	"time"
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
}
