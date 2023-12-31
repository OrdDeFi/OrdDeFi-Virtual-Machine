package file_utils

import (
	"io"
	"os"
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
