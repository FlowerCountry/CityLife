# 设置最低的编译版本
cmake_minimum_required(VERSION 3.10)

# 项目名称
project(CityFife)

# 设置C++版本
set(CMAKE_C_STANDARD 11)
set(CMAKE_CXX_STANDARD 11)

# 公共头文件目录
include_directories("./include/")

# 枚举头文件
file(GLOB_RECURSE INCLUDES "./include/*.h" "./include/*.hpp")

# 枚举源文件
file(GLOB_RECURSE SOURCES "./src/*.c" "./src/*.cpp")

# PROJECT_SOURCE_DIR 这是项目路径
# PROJECT_BINARY_DIR 这是build里面的文件夹
SET(EXECUTABLE_OUTPUT_PATH ${PROJECT_SOURCE_DIR}/bin) # 设置二进制文件存放路径
SET(LIBRARY_OUTPUT_PATH ${PROJECT_SOURCE_DIR}/firmware) # 设置各种库存放路径

# 生成可执行文件
add_executable(main ${SOURCES})

# ==================== 单元测试支持 ====================
option(BUILD_TESTS "Build unit tests" OFF)

if(BUILD_TESTS)
    enable_testing()

    # 使用 FetchContent 自动下载 GoogleTest
    include(FetchContent)
    FetchContent_Declare(
        googletest
        GIT_REPOSITORY https://github.com/google/googletest.git
        GIT_TAG        v1.14.0
    )
    # 不安装 gtest 到系统
    set(gtest_force_shared_crt ON CACHE BOOL "" FORCE)
    set(INSTALL_GTEST OFF CACHE BOOL "" FORCE)
    FetchContent_MakeAvailable(googletest)

    # 添加单元测试子目录
    add_subdirectory(tests/unit)
endif()
