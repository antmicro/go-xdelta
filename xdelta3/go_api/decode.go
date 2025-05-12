// Copyright (c) 2024-2025 Antmicro

package xdelta

/*
#include <xdelta3/xdelta3-api.h>
#include <stdlib.h>
#cgo LDFLAGS: -lxdelta3 -llzma
*/
import "C"

import (
	"fmt"
	"bytes"
	"io"
	"runtime/cgo"
)

// Performs delta decoding using xdelta3, reconstructing the target data from a base and a delta and
// writing it to the output via the provided writer. It leverages the C xd3_decode function, passing
// Go io.Reader and io.Writer via CGO handles.
//
// Parameters:
//   - base: io.Reader for the source (base) data.
//   - delta: io.Reader for the delta data.
//   - target: io.Writer for the resulting target data.
//
// Returns nil on success, otherwise an error.
func Xd3Decode(base io.Reader, delta io.Reader, target io.Writer) error {
	// Check if delta is empty
	if seeker, ok := delta.(io.Seeker); ok {
		// Seekable reader: peek one byte and restore position
		buf := make([]byte, 1)
		n, err := delta.Read(buf)
		if err != nil && err != io.EOF {
			return fmt.Errorf("failed to check delta: %v", err)
		}
		if n == 0 && err == io.EOF {
			return ErrDeltaIsEmpty
		}
		// Restore position
		if _, err := seeker.Seek(0, io.SeekStart); err != nil {
			return fmt.Errorf("failed to restore delta position: %v", err)
		}
	} else {
		// Non-seekable reader: buffer one byte and chain back
		buf := make([]byte, 1)
		n, err := delta.Read(buf)
		if err != nil && err != io.EOF {
			return fmt.Errorf("failed to check delta: %v", err)
		}
		if n == 0 && err == io.EOF {
			return ErrDeltaIsEmpty
		}
		// Chain buffered data back with original reader
		delta = io.MultiReader(bytes.NewReader(buf[:n]), delta)
	}

	// Create CGO handles for io.Reader and io.Writer
	baseHandle := cgo.NewHandle(base)
	deltaHandle := cgo.NewHandle(delta)
	targetHandle := cgo.NewHandle(target)
	defer baseHandle.Delete()
	defer deltaHandle.Delete()
	defer targetHandle.Delete()

	// Call xd3_decode and convert C.int to Go int
	ret := C.xd3_go_decode(C.uintptr_t(baseHandle), C.uintptr_t(deltaHandle), C.uintptr_t(targetHandle))

	if int(ret) != 0 {
		return fmt.Errorf("xDelta decoding failed with return code: %d", ret)
	}
	return nil
}
