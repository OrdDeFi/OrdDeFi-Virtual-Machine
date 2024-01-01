package file_utils

import (
	"io"
	"os"
	"path/filepath"
)

func CopyFile(src, dest string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()
	destinationFile, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer destinationFile.Close()

	bufferSize := 4 * 1024 * 1024
	buffer := make([]byte, bufferSize)
	for {
		bytesRead, err := sourceFile.Read(buffer)
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}
		_, err = destinationFile.Write(buffer[:bytesRead])
		if err != nil {
			return err
		}
	}
	return nil
}

func CopyDir(srcPath, destPath string) error {
	err := os.RemoveAll(destPath)
	if err != nil {
		return err
	}
	err = os.MkdirAll(destPath, 0755)
	if err != nil {
		return err
	}
	entries, err := os.ReadDir(srcPath)
	if err != nil {
		return err
	}
	for _, entry := range entries {
		srcFile := filepath.Join(srcPath, entry.Name())
		destFile := filepath.Join(destPath, entry.Name())
		if entry.IsDir() {
			// recursive copy dir
			err := CopyDir(srcFile, destFile)
			if err != nil {
				return err
			}
		} else {
			err := CopyFile(srcFile, destFile)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
