// +build cgo,!time_compiled,darwin,amd64 cgo,!time_compiled,freebsd,amd64 cgo,!time_compiled,dragonfly,amd64

package base

/*
#cgo CFLAGS: -I${SRCDIR}/cpp
#cgo LDFLAGS: -L${SRCDIR}/binary/darwin -ltime -lstdc++
#include "cpp/library.h"
*/
import "C"
import "time"

func Now() time.Time {
	return time.Unix(0, int64(C.now()))
}
