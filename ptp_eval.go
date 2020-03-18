// Copyright {YEAR} The Cockroach Authors.
//
// Use of this software is governed by the Business Source License
// included in the file licenses/BSL.txt.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the Apache License, Version 2.0, included in the file
// licenses/APL.txt.

package main

/*
#include <time.h>
#include <unistd.h>
extern int clock_gettime(clockid_t clock_id, struct timespec *tp);
*/
import "C"

import (
	"fmt"
	"os"
	"time"
	"syscall"
	"unsafe"
)
const count = 5_000_000

//go:linkname walltime runtime.walltime
func walltime() (sec int64, nsec int32)

//go:linkname nanotime runtime.nanotime
func nanotime() uint64

func main() {
	start := time.Now()
	for i := 0; i < count; i++ {
		_ = time.Now()
	}
	end := time.Now()
	fmt.Printf("Time per now() call %v\n", end.Sub(start)/count)

	start = time.Now()
	sec, nsec := walltime()
	startNSec := uint64(sec)*1e9 + uint64(nsec)
	for i := 0; i < count; i++ {
		sec, nsec = walltime()
	}
	end = time.Now()
	endNSec := uint64(sec)*1e9 + uint64(nsec)
	fmt.Printf("Time per walltime() call %v, nsec diff: %v\n", end.Sub(start)/count, (endNSec-startNSec)/count)

	start = time.Now()
	startNSec = nanotime()
	for i := 0; i < count; i++ {
		endNSec = nanotime()
	}
	end = time.Now()
	fmt.Printf("Time per nanotime() call %v, nsec diff: %v\n", end.Sub(start)/count, (endNSec-startNSec)/count)

	ptp_dev, err := os.Open("/dev/ptp0")
	if err == nil {
		ptp_fd := ptp_dev.Fd()
		ptp_mod_fd := (^ptp_fd << 3) | 3
		fmt.Printf("Opened /dev/ptp0 with fd %x, mod_fd %x \n", ptp_fd, ptp_mod_fd)
		TryGetTimeCGO(ptp_mod_fd, "C.clock_gettime(/dev/ptp0)")
		TryGetTimeSyscall(ptp_mod_fd, "gettime(/dev/ptp0)")
	} else {
		fmt.Printf("Can't open /dev/ptp0: %+v\n", err)
	}

	TryGetTimeCGO(C.CLOCK_REALTIME, "C.clock_gettime(CLOCK_REALTIME)")
	TryGetTimeSyscall(C.CLOCK_REALTIME, "gettime(CLOCK_REALTIME)")
	TryGetTimeCGO(C.CLOCK_MONOTONIC, "C.clock_gettime(CLOCK_MONOTONIC)")
	TryGetTimeSyscall(C.CLOCK_MONOTONIC, "gettime(CLOCK_MONOTONIC)")
	fmt.Printf("now is %s\n", time.Now())
}

func TryGetTimeCGO(clockId uintptr, text string) {
	start := time.Now()
	var ts C.struct_timespec
	_ = C.clock_gettime(C.clockid_t(clockId), &ts)
	startNSec := uint64(ts.tv_sec)*1e9 + uint64(ts.tv_nsec)
	for i := 0; i < count; i++ {
		_ = C.clock_gettime(C.clockid_t(clockId), &ts)
	}
	end := time.Now()
	endNSec := uint64(ts.tv_sec)*1e9 + uint64(ts.tv_nsec)
	fmt.Printf("CGO %s call %v, end now: %s, end get time: %s, nsec diff: %v\n",
		text,
		end.Sub(start)/count,
		end,
		time.Unix(int64(ts.tv_sec), int64(ts.tv_nsec)),
		(endNSec-startNSec)/count,
	)

}
func TryGetTimeSyscall(clockId uintptr, text string) {
	start := time.Now()
	startGetTime := gettime(clockId)
	for i := 0; i < count; i++ {
		_ = gettime(clockId)
	}
	end := time.Now()
	endGetTime := gettime(clockId)
	fmt.Printf("Syscall %s call %v, end now: %s, end get time: %s, nsec diff: %v\n",
		text,
		end.Sub(start)/count,
		end,
		time.Unix(endGetTime.Sec, endGetTime.Nsec),
		(endGetTime.Nano()-startGetTime.Nano())/count,
	)
}

func gettime(clock_id uintptr) syscall.Timespec {
	var ts syscall.Timespec
	syscall.Syscall(228, 1, uintptr(unsafe.Pointer(&ts)), 0)
	return ts
}
