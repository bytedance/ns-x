#include "library.h"
#include <chrono>
#include <mutex>

using namespace std::chrono;

constexpr auto AlignPeriod = milliseconds(50);
auto lastSystemTime = system_clock::now();
auto alignedTime = high_resolution_clock::now();
std::mutex lock{};

void align(decltype(alignedTime) &tick) {
    std::lock_guard<std::mutex> lockGuard{lock};
    if (tick - alignedTime > AlignPeriod) {
        lastSystemTime = system_clock::now();
        alignedTime = tick;
    }
}

int64_t now() {
    auto tick = high_resolution_clock::now();
    if (tick - alignedTime > AlignPeriod) {
        align(tick);
    }
    return duration_cast<nanoseconds>(tick - alignedTime + lastSystemTime.time_since_epoch()).count();
}