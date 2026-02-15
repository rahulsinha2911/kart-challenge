package utils

import "os"

func ShouldReadFromServer() bool {
	return os.Getenv("READ_FROM_SERVER") == "true"
}
