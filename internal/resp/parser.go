package resp

import (
	"bytes"
	"errors"
	"fmt"
	"strconv"
)

const (
	sep string = "\r\n"
)

func Parse(data []byte) ([]string, error) {
	if len(data) == 0 {
		return nil, errors.New("invalid command")
	}

	switch data[0] {
	case '*':
		return parseArray(data)
	case '$':
		bulkStr, _, err := parseBulkString(data)
		return []string{bulkStr}, err
	default:
		return nil, fmt.Errorf("unknown RESP type: %c", data[0])
	}
}

func parseArray(data []byte) ([]string, error) {
	line, nextIdx, err := readNext(data[1:])
	if err != nil {
		return nil, err
	}

	num, err := strconv.Atoi(string(line))
	if err != nil {
		return nil, err
	}

	// 1 to account for data[1:] (first byte)
	currPos := 1 + nextIdx
	elems := make([]string, num)
	for i := range num {
		if len(data) <= currPos {
			return nil, errors.New("incomplete array")
		}

		element, n, err := parseBulkString(data[currPos:])
		if err != nil {
			return nil, err
		}

		currPos += n
		elems[i] = element
	}

	return elems, nil
}

func parseBulkString(data []byte) (string, int, error) {
	line, nextIdx, err := readNext(data[1:])
	if err != nil {
		return "", 0, err
	}

	length, err := strconv.Atoi(string(line))
	if err != nil {
		return "", 0, err
	}

	start := 1 + nextIdx
	end := start + length

	if len(data) < end+2 {
		return "", 0, errors.New("incomplete bulk string")
	}

	return string(data[start:end]), end + 2, nil
}

func readNext(data []byte) ([]byte, int, error) {
	idx := bytes.Index(data, []byte(sep))
	if idx == -1 {
		return nil, idx, errors.New("terminator not found")
	}

	return data[:idx], idx + 2, nil
}
