// Copyright (c) 2024-2025 Antmicro

package xdelta

import (
	"path/filepath"
	"testing"
)

func TestXd3DecodeNormal(t *testing.T) {
	source := openFile(t, filepath.Join(TEST_DATA_DIR, "source_normal.bin"))
	defer source.Close()
	delta := openFile(t, filepath.Join(TEST_DATA_DIR, "delta_normal.bin"))
	defer delta.Close()
	decoded := createFile(t, "decoded_normal.bin")
	defer decoded.Close()

	if err := Xd3Decode(source, delta, decoded); err != nil {
		t.Fatalf("Failed: %v", err)
	}
	assertFileValid(t, decoded, "Normal decode")

	if !filesAreEqual(t, filepath.Join(TEST_DATA_DIR, "target_modified.bin"), decoded.Name()) {
		t.Fatal("Decoded file doesn't match target file")
	}
}

func TestXd3DecodeEmptyDelta(t *testing.T) {
	source := openFile(t, filepath.Join(TEST_DATA_DIR, "source_normal.bin"))
	defer source.Close()
	delta := openFile(t, filepath.Join(TEST_DATA_DIR, "delta_empty.bin"))
	defer delta.Close()
	decoded := createFile(t, "decoded_empty_delta.bin")
	defer decoded.Close()

	if err := Xd3Decode(source, delta, decoded); err == nil {
		t.Fatal("Xd3Decode should fail for empty delta")
	}
}

func TestXd3DecodeInvalidBase(t *testing.T) {
	source := openFile(t, filepath.Join(TEST_DATA_DIR, "source_empty.bin"))
	defer source.Close()
	delta := openFile(t, filepath.Join(TEST_DATA_DIR, "delta_normal.bin"))
	defer delta.Close()
	decoded := createFile(t, "decoded_invalid_base.bin")
	defer decoded.Close()

	if err := Xd3Decode(source, delta, decoded); err == nil {
		t.Fatal("Xd3Decode should fail for invalid base")
	}
}
