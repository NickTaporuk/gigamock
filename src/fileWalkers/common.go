package fileWalkers

import (
	"strings"
)

const (
	inMemoryKeySeparator = "|"
)

// PrepareInMemoryStoreKey
func PrepareInMemoryStoreKey(path, method string) string {
	return path + inMemoryKeySeparator + strings.ToUpper(method)
}
