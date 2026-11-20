package resp

import "fmt"

func SimpleString(value string) string {
	return fmt.Sprintf("+%s%s", value, sep)
}

func BulkString(value string) string {
	// Null bulk string
	if len(value) == 0 {
		return fmt.Sprintf("$-1%s", sep)
	}

	return fmt.Sprintf("$%d%s%s%s", len(value), sep, value, sep)
}
