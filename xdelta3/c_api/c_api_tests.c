// Copyright (c) 2024-2025 Antmicro

#include <errno.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <sys/stat.h>
#include <unistd.h>

#include "unity.h"
#include "xdelta3-api.h"

#define TEST_DATA_DIR "../test_data"
#define FILEPATH_MAX 1024

static char source_file[FILEPATH_MAX];
static char input_file[FILEPATH_MAX];
static char output_file[FILEPATH_MAX];
static char decoded_file[FILEPATH_MAX];
static char temp_dir[FILEPATH_MAX];

static void CreateTempDir(void) {
  snprintf(temp_dir, FILEPATH_MAX, "/tmp/xd3_test_XXXXXX");
  if (mkdtemp(temp_dir) == NULL) {
    TEST_FAIL_MESSAGE("Failed to create temp directory");
  }
}

// Helper to check if a file exists and is non-empty
static void AssertFileValid(const char *filename, const char *context) {
  struct stat st;
  char msg[256];
  snprintf(msg, sizeof(msg), "%s: File not created", context);
  TEST_ASSERT_EQUAL_INT_MESSAGE(0, stat(filename, &st), msg);
  snprintf(msg, sizeof(msg), "%s: File is empty", context);
  TEST_ASSERT_GREATER_THAN_INT_MESSAGE(0, st.st_size, msg);
}

static int FilesAreEqual(const char *file1, const char *file2) {
  FILE *f1 = fopen(file1, "rb");
  FILE *f2 = fopen(file2, "rb");
  if (f1 == NULL || f2 == NULL) {
    if (f1)
      fclose(f1);
    if (f2)
      fclose(f2);
    return 0;
  }
  int equal = 1;
  while (1) {
    int c1 = fgetc(f1);
    int c2 = fgetc(f2);
    if (c1 != c2) {
      equal = 0;
      break;
    }
    if (c1 == EOF) {
      break;
    }
  }
  fclose(f1);
  fclose(f2);
  return equal;
}

void setUp(void) {
  // Ensure the test data directory exists
  struct stat st = {0};
  if (stat(TEST_DATA_DIR, &st) == -1) {
    TEST_FAIL_MESSAGE("../test_data directory doesn't exist");
  }
  // Create temporary directory
  CreateTempDir();
}

void tearDown(void) { rmdir(temp_dir); }

void test_xd3_encode_normal(void) {
  snprintf(source_file, sizeof(source_file), "%s/source_normal.bin", TEST_DATA_DIR);
  snprintf(input_file, sizeof(input_file), "%s/target_modified.bin",
           TEST_DATA_DIR);
  snprintf(output_file, sizeof(output_file), "%s/delta_output_normal.bin",
           temp_dir);

  int ret = xd3_encode(source_file, input_file, output_file);
  TEST_ASSERT_EQUAL_INT_MESSAGE(0, ret,
                                "xd3_encode failed for normal scenario");
  AssertFileValid(output_file, "Normal encode");
}

void test_xd3_encode_empty_source(void) {
  snprintf(source_file, sizeof(source_file), "%s/source_empty.bin",
           TEST_DATA_DIR);
  snprintf(input_file, sizeof(input_file), "%s/target_normal.bin",
           TEST_DATA_DIR);
  snprintf(output_file, sizeof(output_file), "%s/delta_output_empty_source.bin",
           temp_dir);

  int ret = xd3_encode(source_file, input_file, output_file);
  TEST_ASSERT_EQUAL_INT_MESSAGE(0, ret, "xd3_encode failed for empty source");
  AssertFileValid(output_file, "Empty source encode");
}

void test_xd3_encode_empty_target(void) {
  snprintf(source_file, sizeof(source_file), "%s/source_normal.bin", TEST_DATA_DIR);
  snprintf(input_file, sizeof(input_file), "%s/target_empty.bin",
           TEST_DATA_DIR);
  snprintf(output_file, sizeof(output_file), "%s/delta_output_empty_target.bin",
           temp_dir);

  int ret = xd3_encode(source_file, input_file, output_file);
  TEST_ASSERT_EQUAL_INT_MESSAGE(0, ret, "xd3_encode failed for empty target");
  AssertFileValid(output_file, "Empty target encode");
}

void test_xd3_decode_normal(void) {
  snprintf(source_file, sizeof(source_file), "%s/source_normal.bin", TEST_DATA_DIR);
  snprintf(output_file, sizeof(output_file), "%s/delta_normal.bin",
           TEST_DATA_DIR);
  snprintf(decoded_file, sizeof(decoded_file), "%s/decoded_normal.bin",
           temp_dir);

  int ret = xd3_decode(source_file, output_file, decoded_file);
  TEST_ASSERT_EQUAL_INT_MESSAGE(0, ret,
                                "xd3_decode failed for normal scenario");

  snprintf(input_file, sizeof(input_file), "%s/target_modified.bin",
           TEST_DATA_DIR);
  TEST_ASSERT_TRUE_MESSAGE(FilesAreEqual(input_file, decoded_file),
                           "Decoded file does not match target");
}

void test_xd3_decode_empty_delta(void) {
  snprintf(source_file, sizeof(source_file), "%s/source_normal.bin", TEST_DATA_DIR);
  snprintf(output_file, sizeof(output_file), "%s/delta_empty.bin",
           TEST_DATA_DIR);
  snprintf(decoded_file, sizeof(decoded_file), "%s/decoded_empty_delta.bin",
           temp_dir);

  int ret = xd3_decode(source_file, output_file, decoded_file);
  TEST_ASSERT_NOT_EQUAL_INT_MESSAGE(0, ret,
                                    "xd3_decode should fail for empty delta");
}

void test_xd3_decode_invalid_base(void) {
  snprintf(source_file, sizeof(source_file), "%s/source_empty.bin",
           TEST_DATA_DIR);
  snprintf(output_file, sizeof(output_file), "%s/delta_normal.bin",
           TEST_DATA_DIR);
  snprintf(decoded_file, sizeof(decoded_file), "%s/decoded_invalid_base.bin",
           temp_dir);

  int ret = xd3_decode(source_file, output_file, decoded_file);
  TEST_ASSERT_NOT_EQUAL_INT_MESSAGE(0, ret,
                                    "xd3_decode should fail for invalid base");
}

int main(void) {
  UNITY_BEGIN();
  RUN_TEST(test_xd3_encode_normal);
  RUN_TEST(test_xd3_encode_empty_source);
  RUN_TEST(test_xd3_encode_empty_target);
  RUN_TEST(test_xd3_decode_normal);
  RUN_TEST(test_xd3_decode_empty_delta);
  RUN_TEST(test_xd3_decode_invalid_base);
  return UNITY_END();
}
