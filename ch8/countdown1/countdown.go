// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

// See page 244.

// Countdown implements the countdown for a rocket launch.
package main

import (
	"fmt"
	"time"
)

//!+
func main() {
	fmt.Println("Commencing countdown.")
	//这个返回一个channel,会周期性地像一个节拍器一样向这个channel发送事件,节拍控制交给了time这个进程
	tick := time.Tick(1 * time.Second)
	//等待10个节拍
	for countdown := 10; countdown > 0; countdown-- {
		fmt.Println(countdown)
		<-tick

	}
	launch()
}

//!-

func launch() {
	fmt.Println("Lift off!")
}
