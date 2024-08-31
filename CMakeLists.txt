# 设置最低的编译版本
cmake_minimum_required(VERSION 3.30)

# 设置目标平台与标准库文件
set(CMAKE_CXX_FLAGS "${CMAKE_CXX_FLAGS} --target=x86_64-w64-mingw32 -stdlib=libstdc++")
set(CMAKE_C_FLAGS "${CMAKE_C_FLAGS} --target=x86_64-w64-mingw32 -stdlib=libstdc++")

# 项目名称
project(CityFife)

# 设置C++版本
set(CMAKE_C_STANDARD 11)
set(CMAKE_CXX_STANDARD 11)

# 头文件路径
INCLUDE_DIRECTORIES(
    "./include/"
    "C:/Users/FlowerCity/CityLife/PDCurses-3.9" # 添加PDCurses头文件路径
)

# 枚举头文件
file(GLOB_RECURSE INCLUDES "./include/*.h" "./include/*.hpp")

# 枚举源文件
file(GLOB_RECURSE SOURCES "./src/*.c" "./src/*.cpp")

# PROJECT_SOURCE_DIR 这是项目路径
# PROJECT_BINARY_DIR 这是build里面的文件夹
SET(EXECUTABLE_OUTPUT_PATH ${PROJECT_SOURCE_DIR}/bin) # 设置二进制文件存放路径
SET(LIBRARY_OUTPUT_PATH ${PROJECT_SOURCE_DIR}/firmware) # 设置各种库存放路径

# 添加需要链接的库文件目录 新建好了 链接的时候用
LINK_DIRECTORIES(${LIBRARY_OUTPUT_PATH} "C:/Users/FlowerCity/CityLife/PDCurses-3.9/wincon")

# 生成可执行文件
add_executable(main ${SOURCES})

# 链接PDCurses库
target_link_libraries(main pdcurses.a)