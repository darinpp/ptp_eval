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

import (
	"fmt"
	"os"
	"time"
	"syscall"
	"unsafe"
)

func main() {
	const count = 5_000_000
	start := time.Now()
	for i := 0; i < count; i++ {
		_ = time.Now()
	}
	fmt.Printf("Time per now() call %v\n", time.Now().Sub(start)/count)

	ptp_dev, err := os.Open("/dev/ptp0")
	if err == nil {
		ptp_fd := ptp_dev.Fd()
		fmt.Printf("Opened /dev/ptp0 with fd %d", ptp_fd)

		start = time.Now()
		for i := 0; i < count; i++ {
			_ = gettime(ptp_fd)
		}
		fmt.Printf("Time per gettime(/dev/ptp0) call %v\n", time.Now().Sub(start)/count)
	} else {
		fmt.Printf("Can't open /dev/ptp0: %+v", err)
	}


	start = time.Now()
	for i := 0; i < count; i++ {
		_ = gettime(0)
	}
	fmt.Printf("Time per gettime(CLOCK_REALTIME) call %v\n", time.Now().Sub(start)/count)

	start = time.Now()
	for i := 0; i < count; i++ {
		_ = gettime(1)
	}
	fmt.Printf("Time per gettime(CLOCK_MONOTONIC) call %v\n", time.Now().Sub(start)/count)

	
}

func gettime(clock_id uintptr) uint64 {
	var ts syscall.Timespec
	syscall.Syscall(syscall.SYS_CLOCK_GETTIME, 1, uintptr(unsafe.Pointer(&ts)), 0)
	return uint64(ts.Nano())
}
