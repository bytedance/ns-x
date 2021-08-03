// +build cgo,time_compiled

package base

/*
#cgo CFLAGS: -I${SRCDIR}/cpp
#cgo LDFLAGS: -L${SRCDIR}/cpp -ltime -lstdc++
#include "cpp/library.h"
*/
import "C"
import "time"

func Now() time.Time {
	return time.Unix(0, int64(C.now()))
}
