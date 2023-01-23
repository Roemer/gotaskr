package goext

import (
	"encoding/json"
	"errors"
	"os"
)

// WriteJsonToFile writes the given object into a file.
func WriteJsonToFile(object any, outputFilePath string, indented bool) error {
	var data []byte
	var err error
	if indented {
		data, err = json.MarshalIndent(object, "", "  ")
	} else {
		data, err = json.Marshal(object)
	}
	if err != nil {
		return err
	}
	if err := os.WriteFile(outputFilePath, data, os.ModePerm); err != nil {
		return err
	}
	return nil
}

// FileExists checks if a file exists (and it is not a directory).
func FileExists(filePath string) (bool, error) {
	info, err := os.Stat(filePath)
	if err == nil {
		return !info.IsDir(), nil
	}
	if errors.Is(err, os.ErrNotExist) {
		return false, nil
	}
	return false, err
}

// CopyFile is a simple file copy from source to destination.
func CopyFile(src string, dst string) (int64, error) {
	srcFile, err := os.Open(src)
	if err != nil {
		return -1, err
	}
	defer srcFile.Close()
	dstFile, err := os.Create(dst)
	if err != nil {
		return -1, err
	}
	defer dstFile.Close()
	return dstFile.ReadFrom(srcFile)
}
