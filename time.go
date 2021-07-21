package networksimulator

/*
#cgo CFLAGS: -I${SRCDIR}/cpp
#cgo LDFLAGS: -L${SRCDIR}/cpp -ltime -lstdc++
#include "cpp/library.h"
*/
import "C"
import "time"

func Now() time.Time {
	t := int64(C.now())
	println(t)
	return time.Unix(0, t)
}
