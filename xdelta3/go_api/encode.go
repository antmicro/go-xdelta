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
	"io"
	"runtime/cgo"
)

// Performs delta encoding using xdelta3, creating a delta from a base to a target data and
// writing it to the output via the provided writer. It leverages the C xd3_encode function,
// passing Go io.Reader and io.Writer via CGO handles.
//
// Parameters:
//   - base: io.Reader for the source (base) data.
//   - target: io.Reader for the target data.
//   - delta: io.Writer for the resulting delta data.
//
// Returns nil on success, otherwise an error
func Xd3Encode(base io.Reader, target io.Reader, delta io.Writer) error {
	// Create CGO handles for io.Reader and io.Writer
	baseHandle := cgo.NewHandle(base)
	targetHandle := cgo.NewHandle(target)
	deltaHandle := cgo.NewHandle(delta)
	defer baseHandle.Delete()
	defer targetHandle.Delete()
	defer deltaHandle.Delete()

	// Call xd3_encode and convert C.int to Go int
	ret := C.xd3_go_encode(C.uintptr_t(baseHandle), C.uintptr_t(targetHandle), C.uintptr_t(deltaHandle))

	if int(ret) != 0 {
		return fmt.Errorf("xDelta encoding failed with return code: %d", ret)
	}
	return nil
}
