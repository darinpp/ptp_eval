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
	"time"
)

func main() {
	const count = 5_000_000
	start := time.Now()
	for i := 0; i < count; i++ {
		_ = time.Now()
	}
	fmt.Printf("Time per now() call %v\n", time.Now().Sub(start)/count)
}