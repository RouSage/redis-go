package commands

import (
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
		arg := strings.ToUpper(args[i])

		switch arg {
		case "EX", "PX":
			if options.ttl > 0 {
				return setOptions{}, errSyntax
			}
			if i+1 >= len(args) {
				return setOptions{}, errSyntax
			}

			unit := time.Second
			if arg == "PX" {
				unit = time.Millisecond
			}

			ttl, err := strconv.ParseInt(args[i+1], 10, 64)
			if err != nil || ttl <= 0 {
				return setOptions{}, errInvalidExpire
			}

			options.ttl = time.Duration(ttl) * unit
			i++ // skip the ms value

		default:
			return setOptions{}, errSyntax
		}
	}

	return options, nil
}
