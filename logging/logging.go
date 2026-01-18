package logging

import (
	"fmt"
	"os"
)

func Log(message string) {
	if os.Getenv("LOGGING") == "true" {
		fmt.Println(message)
	}
}
