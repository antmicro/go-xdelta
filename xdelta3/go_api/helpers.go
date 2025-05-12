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
	"errors"
	"io"
	"runtime/cgo"
	"unsafe"
)

// Xdelta error codes (look at `xd3_rvalues` enum in xdelta3.h)
const (
	ErrXd3Internal         = C.int(-17710)
	ErrXd3InvalCgoHandle   = C.int(-17715)
)

var (
	ErrDeltaIsEmpty       = errors.New("Delta is empty")
	ErrBufferAllocFail    = errors.New("Failed to allocate buffer")
	ErrXd3ReadFail        = errors.New("Xd3ReadFromGoStream: read failed")
	ErrXd3WriteFail       = errors.New("Xd3WriteToGoStream: write failed")
	ErrXd3SeekInvalOrigin = errors.New("Xd3SeekGoStream: invalid origin")
	ErrXd3SeekFail        = errors.New("Xd3SeekGoStream: seek failed")
)

// Reads data via the provided io.Reader and fills xdelta's C buffer with up to
// `size` bytes, storing the number read in `nread`.
//
// Note: Xdelta requests data blocks of the size equal to XD3_DEFAULT_WINSIZE.
// Near EOF, fewer bytes may be available than the `size` param specifies.
// Thus, the function handle EOF and read all available bytes.
//
// Parameters:
//   - buf: Pointer to the C buffer to fill with data.
//   - size: Number of bytes requested by Xdelta.
//   - nread: Pointer to store the number of bytes actually read.
//   - goReaderHandle: cgo.Handle to the Go io.Reader.
//
// Returns 0 on success or:
//   - XD3_INVALID_CGO_HANDLE if the handle is invalid.
//   - XD3_INTERNAL on I/O errors (except EOF or ErrUnexpectedEOF, which are treated as success due to xdelta's window size behaviour).
//
//export xd3_ReadFromGoStream
func xd3_ReadFromGoStream(buf unsafe.Pointer, size C.size_t, nread *C.size_t, goReaderHandle C.uintptr_t) C.int {
	// Retrieve the reader from the cgo handle
	readerHandle := cgo.Handle(goReaderHandle)
	readerVal, ok := readerHandle.Value().(io.Reader)
	if !ok {
		return ErrXd3InvalCgoHandle
	}
	reader := readerVal

	// Convert the C buffer to a Go []byte slice
	bufSlice := unsafe.Slice((*byte)(buf), size)

	// Attempt to read up to `size` bytes into the buffer
	n, err := io.ReadFull(reader, bufSlice)
	*nread = C.size_t(n)

	// Note: xdelta may request more bytes than available due to window size
	// Thus, EOF and ErrUnexpectedEOF are treated as success.
	if err != nil && err != io.EOF && err != io.ErrUnexpectedEOF {
		return ErrXd3Internal
	}
	return 0
}

// Writes data from xdelta's C buffer via the provided io.Writer. Writes all
// `size` bytes or returns an error if unable to do so.
//
// Parameters:
//   - buf: Pointer to the C buffer containing data to write.
//   - size: Number of bytes to write.
//   - goWriterHandle: cgo.Handle to the Go io.Writer.
//
// Returns 0 on success or:
//   - XD3_INVALID_CGO_HANDLE if the handle is invalid.
//   - XD3_INTERNAL on I/O errors.
//
//export xd3_WriteToGoStream
func xd3_WriteToGoStream(buf unsafe.Pointer, size C.size_t, goWriterHandle C.uintptr_t) C.int {
	// Retrieve the writer from the cgo handle
	writerHandle := cgo.Handle(goWriterHandle)
	writerVal, ok := writerHandle.Value().(io.Writer)
	if !ok {
		return ErrXd3InvalCgoHandle
	}
	writer := writerVal

	// Convert the C buffer to a Go []byte slice
	bufSlice := unsafe.Slice((*byte)(buf), size)

	// Attempt to write all bytes from the buffer to the writer
	// In case of partial write, Write() returns an error
	_, err := writer.Write(bufSlice)

	// Check for errors
	if err != nil {
		return ErrXd3Internal
	}
	return 0
}

// Seeks to a position via the provided io.Seeker using the given offset and
// origin.
//
// Parameters:
//   - offset: The offset to seek to, interpreted based on origin.
//   - origin: The origin point for the offset (0: start, 1: current, 2: end).
//   - goReaderHandle: cgo.Handle to the Go io.Seeker.
//
// Returns 0 on success or:
//   - XD3_INVALID_CGO_HANDLE if the handle is invalid or not a seeker.
//   - XD3_INTERNAL if the origin value is invalid, or
//     if the seek operation fails.
//
//export xd3_SeekGoStream
func xd3_SeekGoStream(offset C.longlong, origin C.short, goReaderHandle C.uintptr_t) C.int {
	// Retrieve the seeker from the cgo ReaderHandle
	readerHandle := cgo.Handle(goReaderHandle)
	seekerVal, ok := readerHandle.Value().(io.Seeker)
	if !ok {
		return ErrXd3InvalCgoHandle
	}
	seeker := seekerVal

	// Map C origin value to Go io.Seeker constants
	goOrigin := io.SeekStart
	switch origin {
	case 0: // SEEK_SET
		goOrigin = io.SeekStart
	case 1: // SEEK_CUR
		goOrigin = io.SeekCurrent
	case 2: // SEEK_END
		goOrigin = io.SeekEnd
	default:
		fmt.Println("xd3_SeekGoStream: invalid offset origin value: ", origin)
		return ErrXd3Internal
	}

	// Perform seek operation with the given offset and origin
	_, err := seeker.Seek(int64(offset), goOrigin)

	if err != nil {
		return ErrXd3Internal
	}
	return 0
}

// <-------------- WRAPPERS FOR TESTING PURPOSE -------------->

func WrapXd3ReadFromGoStream(reader io.Reader, size int) error {
	handle := cgo.NewHandle(reader)
	defer handle.Delete()

	buf := C.malloc(C.size_t(size))
	if buf == nil {
		return ErrBufferAllocFail
	}
	defer C.free(buf)

	var nread C.size_t
	if ret := xd3_ReadFromGoStream(buf, C.size_t(size), &nread, C.uintptr_t(handle)); ret != 0 {
		return ErrXd3ReadFail
	}
	return nil
}

func WrapXd3WriteToGoStream(writer io.Writer, data []byte) error {
	handle := cgo.NewHandle(writer)
	defer handle.Delete()

	buf := C.CBytes(data)
	if buf == nil {
		return ErrBufferAllocFail
	}
	defer C.free(buf)

	if ret := xd3_WriteToGoStream(buf, C.size_t(len(data)), C.uintptr_t(handle)); ret != 0 {
		return ErrXd3WriteFail
	}
	return nil
}

func WrapXd3SeekGoStream(seeker io.Seeker, offset int64, origin int) error {
	handle := cgo.NewHandle(seeker)
	defer handle.Delete()

	var cOrigin C.short
	switch origin {
	case io.SeekStart:
		cOrigin = 0
	case io.SeekCurrent:
		cOrigin = 1
	case io.SeekEnd:
		cOrigin = 2
	default:
		return ErrXd3SeekInvalOrigin
	}

	if ret := xd3_SeekGoStream(C.longlong(offset), cOrigin, C.uintptr_t(handle)); ret != 0 {
		return ErrXd3SeekFail
	}
	return nil
}
