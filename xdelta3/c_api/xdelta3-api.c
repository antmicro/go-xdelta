// Copyright (c) 2024-2025 Antmicro

#include "xdelta3-api.h"
#include "xdelta3-list.h"
#include "xdelta3-main.h"

#if XD3_ENCODER
int xd3_encode(const char *source_filename, const char *input_filename,
               const char *output_filename) {
  main_file src_file, in_file, out_file;
  int ret;
  // Reset config to default state
  reset_defaults();
  // Enable overwriting the output file by default
  option_force = 1;

  main_file_init(&src_file);
  main_file_init(&in_file);
  main_file_init(&out_file);

  src_file.filename = source_filename;
  in_file.filename = input_filename;
  out_file.filename = output_filename;

  // Note: SRC and OUT files are not opened here because
  // they are handled internally by main_input()
  ret = main_file_open(&in_file, in_file.filename, XO_READ);

  // Perform encoding
  if (ret == 0)
    ret = main_input(CMD_ENCODE, &in_file, &out_file, &src_file);

  // Cleanup and exit
  main_file_cleanup(&src_file);
  main_file_cleanup(&in_file);
  main_file_cleanup(&out_file);
  return ret;
}
#endif /* XD3_ENCODER */

int xd3_decode(const char *source_filename, const char *input_filename,
               const char *output_filename) {
  main_file src_file, in_file, out_file;
  int ret = 0;
  // Reset config to default state
  reset_defaults();
  // Enable overwriting the output file by default
  option_force = 1;

  main_file_init(&src_file);
  main_file_init(&in_file);
  main_file_init(&out_file);

  src_file.filename = source_filename;
  in_file.filename = input_filename;
  out_file.filename = output_filename;

  // Note: SRC and OUT files are not opened here because
  // they are handled internally by main_input().
  ret = main_file_open(&in_file, in_file.filename, XO_READ);

  // Perform encoding
  if (ret == 0)
    ret = main_input(CMD_DECODE, &in_file, &out_file, &src_file);

  // Cleanup and exit
  main_file_cleanup(&src_file);
  main_file_cleanup(&in_file);
  main_file_cleanup(&out_file);
  return ret;
}
