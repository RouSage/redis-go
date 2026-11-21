package commands

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/codecrafters-io/redis-starter-go/internal/resp"
)

type setOptions struct {
	ttl time.Duration
}

func (h *Handler) set(args []string) (string, error) {
	if err := assertArgs(args, SET, 3); err != nil {
		return "", err
	}

	options, err := parseSetOptions(args[3:])
	if err != nil {
		return "", err
	}

	h.store.Set(args[1], args[2], options.ttl)
	return resp.SimpleString("OK"), nil
}

func parseSetOptions(args []string) (setOptions, error) {
	var options setOptions

	for i := 0; i < len(args); i++ {
		switch strings.ToUpper(args[i]) {
		case "PX":
			if i+1 >= len(args) {
				return setOptions{}, errors.New("syntax error: PX requires a value")
			}

			ms, err := strconv.ParseInt(args[i+1], 10, 64)
			if err != nil || ms <= 0 {
				return setOptions{}, fmt.Errorf("invalid PX expire time in '%s' command: %q", SET, args[i+1])
			}

			options.ttl = time.Duration(ms) * time.Millisecond
			i++ // skip the ms value

		default:
			return setOptions{}, fmt.Errorf("syntax error: unexpected %q", args[i])
		}
	}

	return options, nil
}
