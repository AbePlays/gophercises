package utils

import (
	"encoding/binary"
	"strconv"
)

func Atoi(s string) int {
	i, _ := strconv.Atoi(s)
	return i
}

func Btoi(b []byte) int {
	return int(binary.BigEndian.Uint64(b))
}

func Itob(i int) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(i))
	return b
}
