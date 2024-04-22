package format

import "fmt"

func ErrorTypeAndMessage(packageName string, error error) string {
	return fmt.Sprintf("[%v]\n\terror type: %T\n\terror message: %v\n", packageName, error, error)
}
