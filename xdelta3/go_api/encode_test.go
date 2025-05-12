// Copyright (c) 2024-2025 Antmicro

package xdelta

import (
	"path/filepath"
	"testing"
)

func TestXd3EncodeNormal(t *testing.T) {
	source := openFile(t, filepath.Join(TEST_DATA_DIR, "source_normal.bin"))
	defer source.Close()
	target := openFile(t, filepath.Join(TEST_DATA_DIR, "target_modified.bin"))
	defer target.Close()
	delta := createFile(t, "delta_output.bin")
	defer delta.Close()

	if err := Xd3Encode(source, target, delta); err != nil {
		t.Fatalf("Xd3Encode failed: %v", err)
	}

	info, err := delta.Stat()
	if err != nil {
		t.Fatalf("Delta file not created: %v", err)
	}
	if info.Size() == 0 {
		t.Fatal("Delta file is empty")
	}
}

func TestXd3EncodeEmptySource(t *testing.T) {
	source := openFile(t, filepath.Join(TEST_DATA_DIR, "source_empty.bin"))
	defer source.Close()
	target := openFile(t, filepath.Join(TEST_DATA_DIR, "target_normal.bin"))
	defer target.Close()
	delta := createFile(t, "delta_output_empty_source.bin")
	defer delta.Close()

	if err := Xd3Encode(source, target, delta); err != nil {
		t.Fatalf("Xd3Encode failed: %v", err)
	}

	info, err := delta.Stat()
	if err != nil {
		t.Fatalf("Delta file not created: %v", err)
	}
	if info.Size() == 0 {
		t.Fatal("Delta file is empty")
	}
}

func TestXd3EncodeEmptyTarget(t *testing.T) {
	source := openFile(t, filepath.Join(TEST_DATA_DIR, "source_normal.bin"))
	defer source.Close()
	target := openFile(t, filepath.Join(TEST_DATA_DIR, "target_empty.bin"))
	defer target.Close()
	delta := createFile(t, "delta_output_empty_target.bin")
	defer delta.Close()

	if err := Xd3Encode(source, target, delta); err != nil {
		t.Fatalf("Failed: %v", err)
	}
	assertFileValid(t, delta, "Empty target encode")
}
