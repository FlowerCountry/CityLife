# tests/unit/CMakeLists.txt
# CityLife 单元测试配置

# 创建测试可执行文件
add_executable(unit_tests
    test_cashier.cpp
    test_wallet_inspector.cpp
    ${CMAKE_SOURCE_DIR}/src/payment/cashier.cpp
    ${CMAKE_SOURCE_DIR}/src/object/wallet_inspector.cpp
    ${CMAKE_SOURCE_DIR}/src/locale/locale_manager.cpp
    ${CMAKE_SOURCE_DIR}/src/path/path_manager.cpp
)

# 设置头文件路径
target_include_directories(unit_tests PRIVATE
    ${CMAKE_SOURCE_DIR}/include
)

# 链接 GoogleTest
target_link_libraries(unit_tests
    GTest::gtest_main
)

# 注册测试到 CTest
include(GoogleTest)
gtest_discover_tests(unit_tests)

# 设置输出目录
set_target_properties(unit_tests PROPERTIES
    RUNTIME_OUTPUT_DIRECTORY ${CMAKE_SOURCE_DIR}/bin
)
