package commands

import (
	"strings"
	"time"

	"github.com/codecrafters-io/redis-starter-go/internal/resp"
)

type Command = string

const (
	PING    Command = "PING"
	ECHO    Command = "ECHO"
	SET     Command = "SET"
	GET     Command = "GET"
	COMMAND Command = "COMMAND"
	RPUSH   Command = "RPUSH"
	LPUSH   Command = "LPUSH"
	LPOP    Command = "LPOP"
	LRANGE  Command = "LRANGE"
	LLEN    Command = "LLEN"
)

type Store interface {
	Set(key, value string, ttl time.Duration)
	Get(key string) (string, bool, error)
	RPush(key string, values ...string) (int, error)
	LPush(key string, values ...string) (int, error)
	LPop(key string) (string, bool, error)
	LRange(key string, start, stop int) ([]string, error)
	LLen(key string) (int, error)
}

type Handler struct {
	store Store
}

func New(store Store) *Handler {
	return &Handler{
		store: store,
	}
}

func (h *Handler) Handle(args []string) (string, error) {
	cmd := strings.ToUpper(args[0])

	switch cmd {
	case PING:
		return resp.SimpleString("PONG"), nil

	case ECHO:
		if err := assertArgs(args, cmd, 2); err != nil {
			return "", err
		}

		return resp.BulkString(args[1]), nil

	case SET:
		return h.set(args)

	case GET:
		if err := assertArgs(args, cmd, 2); err != nil {
			return "", err
		}

		value, ok, err := h.store.Get(args[1])
		if err != nil {
			return "", err
		}
		if !ok {
			return resp.NullBulkString(), nil
		}

		return resp.BulkString(value), nil

	case RPUSH:
		return h.rpush(args)

	case LPUSH:
		return h.lpush(args)

	case LPOP:
		return h.lpop(args)

	case LRANGE:
		return h.lrange(args)

	case LLEN:
		return h.llen(args)

	case COMMAND: // random return for COMMAND in order to use redis-cli
		return resp.SimpleString("OK"), nil

	default:
		return "", errUnknownCommand(args[0])
	}
}

func assertArgs(args []string, cmd string, n int) error {
	if len(args) < n {
		return errWrongArgs(cmd)
	}

	return nil
}

func assertExactArgs(args []string, cmd string, n int) error {
	if len(args) != n {
		return errWrongArgs(cmd)
	}

	return nil
}
