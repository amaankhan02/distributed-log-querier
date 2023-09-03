package internal

import (
	"fmt"
	"os"
)

// Utility function that prints a formatted error message to stderr
// and exits the program with exit code 1
func exitOnError(err error) {
	message := "Error: " + err.Error() + "\n"
	fmt.Fprintln(os.Stderr, message)
	os.Exit(1)
}