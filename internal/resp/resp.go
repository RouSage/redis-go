package resp

import "fmt"

const (
	crlf string = "\r\n"
)

func String(value string) string {
	return fmt.Sprintf("+%s%s", value, crlf)
}
