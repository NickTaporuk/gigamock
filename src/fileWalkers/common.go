package fileWalkers

import (
	"strings"
)

const (
	inMemoryKeySeparator = "|"
)

// PrepareImMemoryStoreKey
func PrepareImMemoryStoreKey(path, method string) string {
	return path + inMemoryKeySeparator + strings.ToUpper(method)
}
