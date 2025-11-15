package logger

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

var (
	initialized bool
	mu          sync.Mutex
	file        *os.File
)

func checkDebug() bool {
	_, exists := os.LookupEnv("DEBUG")
	return exists
}

func dateString() string {
	return time.Now().Format("02-01-2006_15:04")
}

func Init(dirPath string) error {
	mu.Lock()
	defer mu.Unlock()

	if initialized {
		return nil
	}

	var err error

	fullPath := filepath.Join(dirPath, dateString()+".log")

	file, err = os.OpenFile(fullPath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return fmt.Errorf("error opening log file: %w", err)
	}

	if checkDebug() {
		if _, err := file.WriteString("\n=====================================================\n==== APP in DEBUG mode ==============================\n"); err != nil {
			return fmt.Errorf("error writing to log file: %w", err)
		}
	} else {
		if _, err := file.WriteString("\n=====================================================\n==== App in NORMAL mode =============================\n"); err != nil {
			return fmt.Errorf("error writing to log file: %w", err)
		}
	}

	initialized = true
	return nil
}

func Log(msg string) {
	mu.Lock()
	defer mu.Unlock()

	if !initialized || file == nil {
		return
	}

	lines := strings.Split(msg, "\n")
	if len(lines) > 1 {
		for i, line := range lines {
			if i == 0 {
				continue
			}
			lines[i] = strings.Repeat("=", 2) + "> " + line
		}
	}

	msg = strings.Join(lines, "\n")

	logLine := fmt.Sprintf("[%s] %s\n", dateString(), msg)
	if _, err := file.WriteString(logLine); err != nil {
		fmt.Println("error writing to log file")
	}

	fmt.Println(logLine)
}

func Debug(msg string) {
	if !checkDebug() {
		return
	}

	Log(fmt.Sprintf("[DEBUG] %s", msg))
}

func Err(err error) {
	Log(fmt.Sprintf("[ERROR] %s", err.Error()))
}

func Critical(err error) {
	Log(fmt.Sprintf("[ERROR] %s\n Exit (1)", err.Error()))
	log.Fatal(err)
}

func Close() error {
	mu.Lock()
	defer mu.Unlock()

	if file != nil {
		err := file.Close()
		file = nil
		return err
	}

	return nil
}
