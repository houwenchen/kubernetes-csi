package helper

import (
	"os"
	"strings"

	"github.com/google/uuid"
)

// GetCaseInsensitiveMap coercs the map's keys to lower case, which only works
// when unicode char is in ASCII subset. May overwrite key-value pairs on
// different permutations of key case as in Key and key. DON'T force values to the
// lower case unconditionally, because values for keys such as mountpoint or
// keylocation are case-sensitive.
// Note that although keys such as 'comPREssion' are accepted and processed,
// even if they are technically invalid, updates to rectify such typing will be
// prohibited as a forbidden update.
func GetCaseInsensitiveMap(dict *map[string]string) map[string]string {
	insensitiveDict := map[string]string{}

	for k, v := range *dict {
		insensitiveDict[strings.ToLower(k)] = v
	}
	return insensitiveDict
}

// GetInsensitiveParameter handles special case ofGetCaseInsensitiveMap looking up one
// key-value pair only
func GetInsensitiveParameter(dict *map[string]string, key string) string {
	insensitiveDict := GetCaseInsensitiveMap(dict)
	return insensitiveDict[strings.ToLower(key)]
}

func exists(path string) (os.FileInfo, bool) {
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return nil, false
	}
	return info, true
}

// FileExists checks if a file exists and is not a directory
func FileExists(filepath string) bool {
	info, present := exists(filepath)
	return present && info.Mode().IsRegular()
}

// DirExists checks if a directory exists
func DirExists(path string) bool {
	info, present := exists(path)
	return present && info.IsDir()
}

// IsValidUUID validates whether a string is a valid UUID
func IsValidUUID(u string) bool {
	_, err := uuid.Parse(u)
	return err == nil
}
