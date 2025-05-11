// Copyright (c) 2024-2025 Antmicro

#ifndef XDELTA3_API_H
#define XDELTA3_API_H

#include "config.h"

#if XD3_ENCODER
/**
 * @brief Delta encoding function
 *
 * @details This function replaces the CLI interface and can be directly called
 * as `xd3_encode()`, achieving the same performance as its CLI counterpart.
 * It offers a better alternative to `xd3_encode_memory()`, which loads entire
 * files into RAM and may lead to high memory usage.
 *
 * @return 0 if success, otherwise error code
 */
int xd3_encode(const char *source_filename, const char *input_filename,
               const char *output_filename);
#endif /* XD3_ENCODER */

/**
 * @brief Delta decoding function
 *
 * @details This function replaces the CLI interface and can be directly called
 * as `xd3_decode()`, achieving the same performance as its CLI counterpart.
 * It offers a better alternative to `xd3_decode_memory()`, which loads entire
 * files into RAM and may lead to high memory usage.
 *
 * @return 0 if success, otherwise error code
 */
int xd3_decode(const char *source_filename, const char *input_filename,
               const char *output_filename);
#endif /* XDELTA3_API */
