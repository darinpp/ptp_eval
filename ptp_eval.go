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
		fmt.Printf("Opened /dev/ptp0 with fd %d\n", ptp_fd)

		start = time.Now()
		startNSec := gettime(ptp_fd)
		for i := 0; i < count; i++ {
			_ = gettime(ptp_fd)
		}
		end := time.Now()
		endNSec := gettime(ptp_fd)
		fmt.Printf("Time per gettime(/dev/ptp0) call %v, nsec diff: %v\n", end.Sub(start)/count, (endNSec-startNSec)/count)
	} else {
		fmt.Printf("Can't open /dev/ptp0: %+v\n", err)
	}


	start = time.Now()
	startNSec := gettime(0)
	for i := 0; i < count; i++ {
		_ = gettime(0)
	}
	end := time.Now()
	endNSec := gettime(0)
	fmt.Printf("Time per gettime(CLOCK_REALTIME) call %v, nsec diff: %v\n", end.Sub(start)/count, (endNSec-startNSec)/count)

	start = time.Now()
	startNSec = gettime(1)
	for i := 0; i < count; i++ {
		_ = gettime(1)
	}
	end = time.Now()
	endNSec = gettime(1)
	fmt.Printf("Time per gettime(CLOCK_MONOTONIC) call %v, nsec diff: %v\n", end.Sub(start)/count, (endNSec-startNSec)/count)

	
}

func gettime(clock_id uintptr) uint64 {
	var ts syscall.Timespec
	syscall.Syscall(syscall.SYS_CLOCK_GETTIME, 1, uintptr(unsafe.Pointer(&ts)), 0)
	return uint64(ts.Nano())
}
