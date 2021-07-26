cmake_minimum_required(VERSION 3.19)
project(high-resolution-time)

set(CMAKE_CXX_STANDARD 11)

add_library(time library.h library.cpp)
