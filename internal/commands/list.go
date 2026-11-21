package commands

import (
	"strconv"

	"github.com/codecrafters-io/redis-starter-go/internal/resp"
)

func (h *Handler) rpush(args []string) (string, error) {
	if err := assertArgs(args, RPUSH, 3); err != nil {
		return "", err
	}

	length, err := h.store.RPush(args[1], args[2:]...)
	if err != nil {
		return "", err
	}

	return resp.Int(length), nil
}

func (h *Handler) lpush(args []string) (string, error) {
	if err := assertArgs(args, RPUSH, 3); err != nil {
		return "", err
	}

	length, err := h.store.LPush(args[1], args[2:]...)
	if err != nil {
		return "", err
	}

	return resp.Int(length), nil
}

func (h *Handler) lrange(args []string) (string, error) {
	if err := assertExactArgs(args, LRANGE, 4); err != nil {
		return "", err
	}

	start, err := strconv.Atoi(args[2])
	if err != nil {
		return "", errNotInteger
	}

	stop, err := strconv.Atoi(args[3])
	if err != nil {
		return "", errNotInteger
	}

	values, err := h.store.LRange(args[1], start, stop)
	if err != nil {
		return "", err
	}

	return resp.Array(values), nil
}

func (h *Handler) llen(args []string) (string, error) {
	if err := assertExactArgs(args, RPUSH, 2); err != nil {
		return "", err
	}

	length, err := h.store.LLen(args[1])
	if err != nil {
		return "", err
	}

	return resp.Int(length), nil
}
