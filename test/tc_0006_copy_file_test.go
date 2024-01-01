package test

import (
	"OrdDeFi-Virtual-Machine/file_utils"
	"strings"
	"testing"
)

func TestCopyFile(t *testing.T) {
	err := file_utils.CopyFile("/Users/satoshi/large.file", "/Users/satoshi/large_copy.file")
	if err != nil {
		if strings.HasSuffix(err.Error(), "no such file or directory") == false {
			t.Errorf("Copy file error: %s", err.Error())
		}
	}
}

func TestCopyDir(t *testing.T) {
	err := file_utils.CopyDir("/Users/satoshi/dir", "/Users/satoshi/dir_copy")
	if err != nil {
		t.Errorf("Copy dir error: %s", err.Error())
	}
}
