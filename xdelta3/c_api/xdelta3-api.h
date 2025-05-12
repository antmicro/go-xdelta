// Copyright (c) 2024-2025 Antmicro

#ifndef XDELTA3_API_H
#define XDELTA3_API_H

#include <stddef.h>
#include <stdint.h>

#include "config.h"

#if XD3_ENCODER
#if CGO_INTEGRATION
/**
 * @brief Delta encoding function using Go interfaces for streaming.
 *
 * @details Replaces the CLI interface with equivalent performance, streaming
 * data via Go interfaces (e.g., io.Reader, io.Writer) to avoid loading whole
 * files into RAM, unlike xd3_encode_memory(). Uses xd3_ReadFromGoStream() and
 * xd3_WriteToGoStream() for CGO integration.
 *
 * @param srcReader CGO handle to a Go reading interface retrieving source data
 * (e.g., io.Reader).
 * @param inReader CGO handle to a Go reading interface retrieving input data
 * (e.g., io.Reader).
 * @param outWriter CGO handle to a Go writing interface for output delta (e.g.,
 * io.Writer).
 *
 * @return 0 on success, otherwise an error code.
 */
int xd3_go_encode(uintptr_t srcReader, uintptr_t inReader, uintptr_t outWriter);
#else
/**
 * @brief Delta encoding function for file-based I/O.
 *
 * @details Replaces the CLI interface with equivalent performance.
 * Offers a better alternative to xd3_encode_memory(), which loads entire
 * files into RAM and may lead to high memory usage.
 *
 * @return 0 if success, otherwise error code
 */
int xd3_encode(const char *source_filename, const char *input_filename,
               const char *output_filename);
#endif /* CGO_INTEGRATION */
#endif /* XD3_ENCODER */

#if CGO_INTEGRATION
/**
 * @brief Reads data blocks from a Go-based data stream via interface like
 * io.Reader, for xdelta3. The implementation is up to user in Go.
 *
 * @details Serves as a C-GO interface to read from a data source via a Go
 * interface (e.g., io.Reader, io.ByteReader) using CGO. Replaces Xdelta's
 * default POSIX/STDIO I/O function (e.g., xd3_posix_io()).
 * The implementation can be found in `go_api/helpers.go`.
 *
 * The function is responsible for:
 *   - Fetching data blocks from the underlying data stream,
 *   - Filling the provided C buffer with the requested number of bytes,
 * Note: Xdelta requests data blocks of the size equal to XD3_DEFAULT_WINSIZE.
 * Near EOF, fewer bytes may be available than the `size` param specifies.
 * Thus, the function should handle EOF and read all available bytes.
 *
 * @param buf C buffer where the read bytes will be stored.
 * @param size Number of bytes requested by the Xdelta engine.
 * @param nread Pointer to a C variable where the function will store the actual
 * number of bytes read.
 * @param goReaderHandle CGO handle to a Go reading interface (e.g., io.Reader).
 *
 * @return 0 on success (including EOF), XD3_INVALID_CGO_HANDLE for invalid
 * handle, or XD3_INTERNAL for I/O errors.
 */
int xd3_ReadFromGoStream(void *buf, size_t size, size_t *nread,
                         uintptr_t goReaderHandle);

/**
 * @brief Writes data blocks to a Go-based data stream via interface like
 * io.Writer, for xdelta3. The implementation is up to user in Go.
 *
 * @details Serves as a C-Go interface to write to a data destination via a Go
 * interface (e.g., io.Writer, io.ByteWriter) using CGO. Replaces Xdelta's
 * default POSIX/STDIO I/O function (e.g., xd3_posix_io()). The implementation
 * can be found in `go_api/helpers.go`.
 *
 * The function is responsible for:
 *   - Writing all `size` bytes from the C buffer to the underlying data
 * destination.
 *
 * @param buf C buffer containing data to write.
 * @param size Bytes to write.
 * @param goWriterHandle CGO handle to a Go writing interface (e.g., io.Writer).
 *
 * @return 0 on success, XD3_INVALID_CGO_HANDLE for invalid handle,
 * or XD3_INTERNAL for I/O errors (e.g., partial writes).
 */
int xd3_WriteToGoStream(void *buf, size_t size, uintptr_t goWriterHandle);

/**
 * @brief Seeks to a position in a Go-based data stream via interface like
 * io.Seeker, for xdelta3. The implementation is up to user in Go.
 *
 * @details Serves as a C-Go interface to perform seek operations on a data
 * source via a Go interface (e.g., io.Seeker) using CGO. Replaces Xdelta's
 * default POSIX/STDIO seek function (e.g., lseek()). The implementation can
 * be found in `go_api/helpers.go`.
 *
 * The function is responsible for:
 *   - Mapping C origin values (0: SEEK_SET, 1: SEEK_CUR, 2: SEEK_END) to Go's
 *     io.SeekStart, io.SeekCurrent, or io.SeekEnd.
 *   - Performing the seek operation on the underlying data source.
 *
 * @param offset Offset to seek to, based on origin.
 * @param origin Reference point: 0 (start), 1 (current), 2 (end).
 * @param goReaderHandle CGO handle to a Go interface of seekable data
 * (e.g., io.Seeker).
 *
 * @return 0 on success, XD3_INVALID_CGO_HANDLE for invalid handle,
 * or XD3_INTERNAL for invalid origin or seek errors.
 */
int xd3_SeekGoStream(long long offset, short origin, uintptr_t goReaderHandle);

/**
 * @brief Delta decoding function using Go interfaces for streaming.
 *
 * @details Replaces the CLI interface with equivalent performance, streaming
 * data via Go interfaces (e.g., io.Reader, io.Writer) to avoid loading whole
 * files into RAM, unlike xd3_decode_memory(). Uses xd3_ReadFromGoStream(),
 * xd3_WriteToGoStream(), and xd3_SeekGoStream for CGO integration.
 *
 * @param srcReader CGO handle to a Go reading interface retrieving source data
 * (e.g., io.Reader).
 * @param inReader CGO handle to a Go reading interface retrieving input (delta)
 * data (e.g., io.Reader).
 * @param outWriter CGO handle to a Go writing interface for output
 * (reconstructed target) data (e.g., io.Writer).
 *
 * @return 0 if success, otherwise error code
 */
int xd3_go_decode(uintptr_t srcReader, uintptr_t inReader, uintptr_t outWriter);
#else
/**
 * @brief Delta decoding function for file-based I/O.
 *
 * @details Replaces the CLI interface with equivalent performance.
 * Offers a better alternative to xd3_decode_memory(), which loads entire
 * files into RAM and may lead to high memory usage.
 *
 * @return 0 if success, otherwise error code
 */
int xd3_decode(const char *source_filename, const char *input_filename,
               const char *output_filename);
#endif /* CGO_INTEGRATION */
#endif /* XDELTA3_API */
