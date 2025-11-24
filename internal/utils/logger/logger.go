package logger

import (
	"fmt"
	"os"
	"strings"
	"sync"
	"time"
)

const logFile = "log.txt"

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
	return time.Now().Format("[02-01-2006_15:04]")
}

func Init() error {
	mu.Lock()
	defer mu.Unlock()

	if initialized {
		return nil
	}

	var err error

	file, err = os.OpenFile(logFile, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
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

	logLine := fmt.Sprintf("%s %s\n", dateString(), msg)
	if _, err := file.WriteString(logLine); err != nil {
		fmt.Println("error writing to log file")
	}
}

func DebugLog(msg string) {
	if !checkDebug() {
		return
	}

	Log(fmt.Sprintf("[DEBUG] %s", msg))
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
