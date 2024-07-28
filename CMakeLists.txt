#设置最低的编译版本
cmake_minimum_required(VERSION 3.10)
#项目名称
project(cppDemo)
set(CMAKE_C_STANDARD 11)#C11
set(CMAKE_CXX_STANDARD 11)#C++11
# add_compile_options(-Wno-unused-variable -Wno-unused-function -Wno-invalid-constexpr)

#头文件路径
INCLUDE_DIRECTORIES(
"./include/"
# "C:/Program Files/boost/boost_1_75_0"
)
#枚举头文件
file(GLOB_RECURSE INCLUDES "./include/*.h" "./include/*.hpp")
#枚举源文件
file(GLOB_RECURSE SOURCES "./src/*.c" "./src/*.cpp")
#加入一些-D
# # add_definitions(${EIGEN_DEFINITIONS})

#PROJECT_SOURCE_DIR 这是项目路径
#PROJECT_BINARY_DIR 这是build里面的文件夹
SET(EXECUTABLE_OUTPUT_PATH ${PROJECT_SOURCE_DIR}/bin) #设置二进制文件存放路径
SET(LIBRARY_OUTPUT_PATH ${PROJECT_SOURCE_DIR}/firmware) #设置各种库存放路径


#添加需要链接的库文件目录 新建好了 链接的时候用
LINK_DIRECTORIES(${LIBRARY_OUTPUT_PATH})
#生成可执行文件
add_executable(main ${SOURCES})
#添加一个run命令
#add_custom_target (run COMMAND ninja all  COMMAND ${EXECUTABLE_OUTPUT_PATH}/main)

