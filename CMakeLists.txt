cmake_minimum_required(VERSION 3.21)
project(git_clones)

include_directories(/usr/local/opt/curl/include/)

set(CMAKE_CXX_STANDARD 17)

set(CURL_LIBRARY "-lcurl")
find_package(CURL REQUIRED)

add_executable(2700 git_clones.cpp)

include_directories(${CURL_INCLUDE_DIR})
target_link_libraries(2700 ${CURL_LIBRARIES})