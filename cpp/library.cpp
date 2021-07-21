#include "library.h"
#include <chrono>

using namespace std::chrono;

constexpr auto AlignPeriod = seconds(1);
auto lastSystemTime = system_clock::now();
auto alignedTime = high_resolution_clock::now();

int64_t now() {
    auto tick = high_resolution_clock::now();
    if (tick - alignedTime > AlignPeriod) {
        lastSystemTime = system_clock::now();
        alignedTime = tick;
    }
    return duration_cast<nanoseconds>(tick - alignedTime + lastSystemTime.time_since_epoch()).count();
}