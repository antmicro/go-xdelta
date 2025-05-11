# xdelta3 fork with updated, simplified API

This project is a fork of the original xdelta3 repository, available at [https://github.com/jmacd/xdelta](https://github.com/jmacd/xdelta).
It's a fork of the `release3_0_apl` branch.

The goal was to simplify the usage of xdelta3 by introducing a new API that achieves the same performance as the CLI version, while avoiding the excessive RAM usage of the `xd3_encode_memory` and `xd3_decode_memory` functions used in the original project.

## New C API

The fork introduces two new functions:
- **`xd3_encode`**: performs delta encoding, taking a source file, an input file, and an output file as arguments. It generates a delta file efficiently, matching the CLI’s performance.
- **`xd3_decode`**: performs delta decoding, reconstructing the target file from a source file and a delta file, also matching CLI performance.

These functions are defined in `c_api/xdelta3-api.h` and implemented in `c_api/xdelta3-api.c`.

## Requirements

- **liblzma-dev**: required for LZMA support when `LZMA=ON`. Install version 5.2.1 or higher.
- **cmake**: required for building. Install version 3.1.0 or higher.

## Build instructions

The project now uses CMake for building. To get started:

1. Ensure you have installed the requirements.
2. Run these commands from the project root:
   ```bash
   mkdir build
   cd build
   cmake ..
   make
   ```
3. Optionally, install the library:
   ```bash
   sudo make install
   ```

## Build options

You can customize the build by adjusting options in `CMakeLists.txt` or passing them with `-D` flags:
- `-DAPP=ON`: build xdelta3 as a standalone application (default: `OFF`, builds as a static library).
- `-DENCODER=OFF`: set to OFF to build only the decoder, without the encoder. (default: `ON`).
- `-DLZMA=OFF`: disable LZMA secondary compression (default: `ON`, requires LibLZMA 5.2.1+).
- `-DBUILD_TESTS=ON`: enable unit tests (default: `OFF`).

## Unit tests

Unit tests utilize the Unity framework. To build and run them:
1. Enable tests with `-DBUILD_TESTS=ON` in CMake.
2. Build the project:
   ```bash
   make
   ```
3. Run the tests:
   ```bash
   ctest
   ```
The tests validate `xd3_encode` and `xd3_decode` using random data and file comparisons.

## Notes

Check the original repository [https://github.com/jmacd/xdelta](https://github.com/jmacd/xdelta) for additional context or to compare this fork with the upstream version.
