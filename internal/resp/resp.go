package resp

import "fmt"

func SimpleString(value string) string {
	return fmt.Sprintf("+%s%s", value, sep)
}

func SimpleError(value string) string {
	return fmt.Sprintf("-%s%s", value, sep)
}

func BulkString(value string) string {
	return fmt.Sprintf("$%d%s%s%s", len(value), sep, value, sep)
}

func BulkError(value string) string {
	return fmt.Sprintf("!%d%s%s%s", len(value), sep, value, sep)
}

func NullBulkString() string {
	return fmt.Sprintf("$-1%s", sep)
}

func Int(value int) string {
	return fmt.Sprintf(":%d%s", value, sep)
}
