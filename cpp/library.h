#ifndef CPP_LIBRARY_H
#define CPP_LIBRARY_H

#include <chrono>

extern "C" {
long long now() {
    return std::chrono::high_resolution_clock::now().time_since_epoch().count();
}
}

#endif //CPP_LIBRARY_H
