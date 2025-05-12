// Copyright (c) 2024-2025 Antmicro

package xdelta

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"
)

const TEST_DATA_DIR = "../test_data"

func openFile(t *testing.T, path string) *os.File {
	f, err := os.Open(path)
	if err != nil {
		t.Fatalf("Cannot open %s: %v", path, err)
	}
	return f
}

func createFile(t *testing.T, name string) *os.File {
	path := filepath.Join(t.TempDir(), name)
	f, err := os.Create(path)
	if err != nil {
		t.Fatalf("Cannot create %s: %v", path, err)
	}
	return f
}

// assertFileValid checks if a file exists and is non-empty
func assertFileValid(t *testing.T, f *os.File, context string) {
	info, err := f.Stat()
	if err != nil {
		t.Fatalf("%s: File not created: %v", context, err)
	}
	if info.Size() == 0 {
		t.Fatalf("%s: File is empty", context)
	}
}

func filesAreEqual(t *testing.T, file1, file2 string) bool {
	data1, err := os.ReadFile(file1)
	if err != nil {
		t.Fatalf("Cannot read %s: %v", file1, err)
	}
	data2, err := os.ReadFile(file2)
	if err != nil {
		t.Fatalf("Cannot read %s: %v", file2, err)
	}
	return bytes.Equal(data1, data2)
}
