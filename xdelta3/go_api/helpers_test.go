// Copyright (c) 2024-2025 Antmicro

package xdelta

import (
	"bytes"
	"io"
	"testing"
)

func TestWrapXd3ReadFromGoStream(t *testing.T) {
	t.Run("Valid", func(t *testing.T) {
		reader := bytes.NewReader([]byte("hello"))
		if err := WrapXd3ReadFromGoStream(reader, 5); err != nil {
			t.Fatalf("Failed: %v", err)
		}
	})

	t.Run("Invalid", func(t *testing.T) {
		if err := WrapXd3ReadFromGoStream(nil, 5); err == nil {
			t.Fatal("Should fail for nil reader")
		}
	})
}

func TestWrapXd3WriteToGoStream(t *testing.T) {
	t.Run("Valid", func(t *testing.T) {
		writer := &bytes.Buffer{}
		if err := WrapXd3WriteToGoStream(writer, []byte("hello")); err != nil {
			t.Fatalf("Failed: %v", err)
		}
		if got := writer.String(); got != "hello" {
			t.Fatalf("Expected 'hello', got '%s'", got)
		}
	})

	t.Run("Invalid", func(t *testing.T) {
		if err := WrapXd3WriteToGoStream(nil, []byte("hello")); err == nil {
			t.Fatal("Should fail for nil writer")
		}
	})
}

func TestWrapXd3SeekGoStream(t *testing.T) {
	t.Run("Valid", func(t *testing.T) {
		reader := bytes.NewReader([]byte("hello"))
		if err := WrapXd3SeekGoStream(reader, 2, io.SeekStart); err != nil {
			t.Fatalf("Failed: %v", err)
		}
		if pos, _ := reader.Seek(0, io.SeekCurrent); pos != 2 {
			t.Fatalf("Expected position 2, got %d", pos)
		}
	})

	t.Run("Invalid", func(t *testing.T) {
		if err := WrapXd3SeekGoStream(nil, 0, io.SeekStart); err == nil {
			t.Fatal("Should fail for nil seeker")
		}
	})
}
