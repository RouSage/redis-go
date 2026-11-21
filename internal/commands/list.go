package commands

import "github.com/codecrafters-io/redis-starter-go/internal/resp"

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
