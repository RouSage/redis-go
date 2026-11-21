package commands

import (
	"errors"
	"fmt"
	"strings"
)

var (
	errSyntax        = errors.New("syntax error")
	errInvalidExpire = errors.New("invalid expire time in 'set' command")
)

func errWrongArgs(cmd string) error {
	return fmt.Errorf("wrong number of arguments for '%s' command", strings.ToLower(cmd))
}

func errUnknownCommand(cmd string) error {
	return fmt.Errorf("unknown command '%s'", cmd)
}
