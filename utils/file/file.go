package file

import (
  "io/ioutil"
  "os"
  "strings"
  "bufio"
)

/**
 check the directory or file exists or not.
 */
func PathIsExist(path string) bool {
  if path == "" {
    return false
  }
  _, err := os.Stat(path)
  return err == nil || os.IsExist(err)
}

func ReadFileToBytes(filePath string) ([]byte, error) {
  return ioutil.ReadFile(filePath)
}

func ReadFileToString(filePath string) (string, error) {
  bs, err := ioutil.ReadFile(filePath)
  if err != nil {
    return "", err
  }
  return string(bs), nil
}

func ReadFileToTrimedString(filePath string) (string, error) {
  bs, err := ioutil.ReadFile(filePath)
  if err != nil {
    return "", err
  }
  return strings.TrimSpace(string(bs)), nil
}

func ReadLine(r *bufio.Reader) ([]byte, error) {
  line, isPrefix, err := r.ReadLine()
  for isPrefix && err == nil {
    var bs []byte
    bs, isPrefix, err = r.ReadLine()
    line = append(line, bs...)
  }

  return line, err
}
