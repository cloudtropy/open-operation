package g

import (
	"io/ioutil"
	"os"
	"strings"
)

// IsExist checks whether a file or directory exists.
// It returns false when the file or directory does not exist.
func PathIsExist(fp string) bool {
	_, err := os.Stat(fp)
	return err == nil || os.IsExist(err)
}

func FileToString(filePath string) (string, error) {
	b, err := ioutil.ReadFile(filePath)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func FileToTrimString(filePath string) (string, error) {
	str, err := FileToString(filePath)
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(str), nil
}
