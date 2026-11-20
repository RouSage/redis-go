package resp

import "fmt"

func SimpleString(value string) string {
	return fmt.Sprintf("+%s%s", value, sep)
}

func BulkString(value string) string {
	return fmt.Sprintf("$%d%s%s%s", len(value), sep, value, sep)
}
