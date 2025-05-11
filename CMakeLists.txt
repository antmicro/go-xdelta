# Copyright (c) 2024-2025 Antmicro

cmake_minimum_required(VERSION 3.25.1)
project(xdelta3 VERSION 3.0.12)

include(CheckTypeSize)
check_type_size("size_t"             XD3_SIZEOF_SIZE_T)
check_type_size("unsigned int"       XD3_SIZEOF_UNSIGNED_INT)
check_type_size("unsigned long"      XD3_SIZEOF_UNSIGNED_LONG)
check_type_size("unsigned long long" XD3_SIZEOF_UNSIGNED_LONG_LONG)

option(APP
  "Set to ON to build standalone application instead of static library."
  OFF
)

option(ENCODER
  "Set to OFF to build only the decoder, without the encoder."
  ON
)

option(DEBUG_SYMBOLS
  "Set to ON to enable debug symbols."
  OFF
)

option(XD3_DEBUG
  "Set to 1, 2 or 3. XD3_DEBUG=1 enables assertions and various statistics.
  Levels >1 enable additional output useful during development and debugging."
  0
)

option(LZMA
  "Set to OFF to disable LZMA secondary compression."
  ON
)

option(VCDIFF_TOOLS
  "Set to ON to enable VCDIFF_TOOLS (merge/print delta headers)."
  OFF
)

option(DJW_FGK
  "Set to ON to include DJW and FGK secondary compression algorithms."
  OFF
)

option(BUILD_TESTS
  "Set to ON to build unit tests."
  OFF
)

# Force encoder on when tests are enabled
if(BUILD_TESTS)
  set(ENCODER ON CACHE BOOL "" FORCE)
endif()

# Main target
if(APP)
  add_executable(${PROJECT_NAME})
else()
  add_library(${PROJECT_NAME} STATIC)
  target_include_directories(${PROJECT_NAME} PUBLIC
    xdelta3
  )
endif()

set(xd3_sources
  xdelta3/xdelta3.c
)
if(NOT APP)
  list(APPEND xd3_sources
    xdelta3/c_api/xdelta3-api.c
  )
endif()

target_sources(${PROJECT_NAME} PRIVATE
  ${xd3_sources}
)

add_compile_options(
  $<$<BOOL:${DEBUG_SYMBOLS}>:-g>
)

target_compile_definitions(${PROJECT_NAME} PRIVATE
    SIZEOF_SIZE_T=${XD3_SIZEOF_SIZE_T}
    SIZEOF_UNSIGNED_INT=${XD3_SIZEOF_UNSIGNED_INT}
    SIZEOF_UNSIGNED_LONG=${XD3_SIZEOF_UNSIGNED_LONG}
    SIZEOF_UNSIGNED_LONG_LONG=${XD3_SIZEOF_UNSIGNED_LONG_LONG}
    XD3_USE_LARGEFILE64=1
    EXTERNAL_COMPRESSION=1

    XD3_MAIN=$<BOOL:${APP}>
    NOT_MAIN=$<NOT:$<BOOL:${APP}>>
    STATIC_LIB=$<NOT:$<BOOL:${APP}>>
    XD3_ENCODER=$<BOOL:${ENCODER}>
    XD3_DEBUG=$<BOOL:${XD3_DEBUG}>
    VCDIFF_TOOLS=$<BOOL:${VCDIFF_TOOLS}>
    SECONDARY_DJW=$<BOOL:${DJW_FGK}>
    SECONDARY_FGK=$<BOOL:${DJW_FGK}>
    SECONDARY_LZMA=$<BOOL:${LZMA}>
)

# Generate config.h with proper values
if(ENCODER)
  set(XD3_ENCODER_VALUE 1)
else()
  set(XD3_ENCODER_VALUE 0)
endif()

configure_file(
  "xdelta3/config.h.in"
  "${CMAKE_SOURCE_DIR}/xdelta3/config.h"
  @ONLY
)

# Link math library
target_link_libraries(${PROJECT_NAME} PRIVATE
  m
)

if(LZMA)
  find_package(LibLZMA 5.2.1 REQUIRED)
  target_link_libraries(${PROJECT_NAME} PRIVATE
    LibLZMA::LibLZMA
  )
endif()

# Install (library only)
include(GNUInstallDirs)
if(NOT APP)
  install(TARGETS ${PROJECT_NAME}
    ARCHIVE DESTINATION ${CMAKE_INSTALL_LIBDIR}
  )
  install(FILES
    "xdelta3/c_api/xdelta3-api.h"
    "xdelta3/config.h"
    DESTINATION ${CMAKE_INSTALL_INCLUDEDIR}/xdelta3
  )
endif()

# Unit tests
if(BUILD_TESTS)
  enable_testing()
  include(FetchContent)
  FetchContent_Declare(
    unity
    GIT_REPOSITORY https://github.com/ThrowTheSwitch/Unity.git
    GIT_TAG v2.5.2
    EXCLUDE_FROM_ALL
  )
  FetchContent_MakeAvailable(unity)

  add_executable(test_xd3_c_api
    xdelta3/c_api/c_api_tests.c
  )
  target_include_directories(test_xd3_c_api PRIVATE
    xdelta3
  )
  target_link_libraries(test_xd3_c_api PRIVATE
    ${PROJECT_NAME}
    unity
    LibLZMA::LibLZMA
    m
  )
  add_test(NAME test_xd3_c_api
    COMMAND test_xd3_c_api
    WORKING_DIRECTORY ${CMAKE_SOURCE_DIR}/xdelta3/c_api
  )
  set_tests_properties(test_xd3_c_api PROPERTIES
    TIMEOUT 300
  )
endif()
