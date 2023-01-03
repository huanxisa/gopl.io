// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

// See page 244.

// Countdown implements the countdown for a rocket launch.
package main

import (
	"fmt"
	"os"
	"time"
)

//!+

//期望用户终止时终止
func main() {
	// ...create abort channel...

	//!-

	//!+abort
	abort := make(chan struct{})
	go func() {
		//等待用户输入
		os.Stdin.Read(make([]byte, 1)) // read a single byte
		abort <- struct{}{}
	}()
	//!-abort

	//如果像原来那样采用循环的形式，那就存在一个问题，如何在一个循环中检查两个channel?
	//	应该可以使用len来检查个数，然后做判断，但是这样有个问题就是如何关闭另外一个chan所在的协程？其实下面的协程也没有关。不用慌
	//!+
	fmt.Println("Commencing countdown.  Press return to abort.")
	//这时候我们需要多路复用（multiplex）这些操作了，为了能够多路复用，我们使用了select语句
	select {
	case <-time.After(10 * time.Second):
		// Do nothing.
	case <-abort:
		fmt.Println("Launch aborted!")
		return
	}
	launch()
}

//!-

func launch() {
	fmt.Println("Lift off!")
}
