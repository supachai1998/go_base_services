package logger

import (
	"fmt"
	"log"
	"os"
	"path"
	"time"

	"github.com/jtolds/gls"
)

const LOGPATH string = "logger/log/app"

var LogFile *os.File

func WriteLog() {
	folderName := fmt.Sprintf("./%s", LOGPATH)
	// Create the folder if it doesn't exist
	_, err := os.Stat(folderName)
	if os.IsNotExist(err) {
		os.MkdirAll(folderName, 0755)
	}
	var fileName string
	L().Info("Start write log at path: ", folderName)
	gls.Go(func() {
		// Wait for a signal to change the log file
		// Get the current date
		now := time.Now()
		// Format the date as "2006-01-02"
		date := now.Format("2006-01-02")

		fileName = date
		// Open a file for logging
		f, err := os.OpenFile(path.Join(folderName, fileName+".log"), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			f.Close()
			log.Fatalf("error opening file: %v", err)
		}
		LogFile = f

		// Log to the file
		if date != time.Now().Format("2006-01-02") {
			f.Close()
		}

		// Wait for 1 hour
		time.Sleep(1 * time.Hour)
	})
}
