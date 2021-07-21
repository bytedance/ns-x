#ifndef CPP_LIBRARY_H
#define CPP_LIBRARY_H

#ifdef __cplusplus

#include "cstdint"

extern "C" {
#else
#include "stdint.h"
#endif
int64_t now();
#ifdef __cplusplus
}
#endif

#endif //CPP_LIBRARY_H
