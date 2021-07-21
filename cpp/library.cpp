#include "library.h"
#include <chrono>

using namespace std::chrono;

const auto offset = system_clock::now().time_since_epoch() - high_resolution_clock::now().time_since_epoch();

int64_t now() {
    auto t = high_resolution_clock::now() + offset;
    return duration_cast<nanoseconds>(t.time_since_epoch()).count();
}