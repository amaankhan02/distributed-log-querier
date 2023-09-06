package test

import (
	"bufio"
	"log"
	"os"
)

var predefined_log_entries = []string{
	"2022-12-30 13:45:00 INFO: Application started",
	"2020-05-14 13:46:00 ERROR: Something went wrong",
	"2021-05-20 13:47:00 INFO: Request received",
	"2023-01-23 13:48:00 DEBUG: Processing...",
}

func generateRandomFiles() {

}
func main() {

	// initialize random number generator

	log_file_path := "test_logs\\"

	// create log file using path
	log_file, err := os.Create(log_file_path)

	if err != nil {
		log.Fatalf("Failed to create log file: %v", err)
	}

	// Create a buffered writer to write to the log file
	writer := bufio.NewWriter(log_file)

}
