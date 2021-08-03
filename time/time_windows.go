// +build cgo,!time_compiled,windows,amd64

package time

/*
#cgo CFLAGS: -I${SRCDIR}/cpp
#cgo LDFLAGS: -L${SRCDIR}/binary/windows -ltime -lstdc++
#include "cpp/library.h"
*/
import "C"
import "time"

func Now() time.Time {
	return time.Unix(0, int64(C.now()))
}
