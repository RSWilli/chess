package uci

import "strconv"

func parseInt(s string) (int, error) {
	i, err := strconv.ParseInt(s, 10, 64)

	return int(i), err
}
