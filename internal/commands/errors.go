package commands

import (
	"errors"
	"fmt"
	"strings"
)

var (
	errSyntax        = errors.New("ERR syntax error")
	errInvalidExpire = errors.New("ERR invalid expire time in 'set' command")
)

func errWrongArgs(cmd string) error {
	return fmt.Errorf("ERR wrong number of arguments for '%s' command", strings.ToLower(cmd))
}

func errUnknownCommand(cmd string) error {
	return fmt.Errorf("ERR unknown command '%s'", cmd)
}
