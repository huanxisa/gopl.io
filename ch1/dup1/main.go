// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

// See page 8.
//!+

// Dup1 prints the text of each line that appears more than
// once in the standard input, preceded by its count.
package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	counts := make(map[string]int)
	//如何从屏幕上获取输入,window终止输入为ctrl+D
	input := bufio.NewScanner(os.Stdin)
	//每次调用 input.Scan()，即读入下一行，并移除行末的换行符；读取的内容可以调用 input.Text() 得到
	for input.Scan() {
		//map 中不含某个键时不用担心，首次读到新行时，等号右边的表达式 counts[line] 的值将被计算为其类型的零值，对于 int 即 0。
		counts[input.Text()]++
	}
	// NOTE: ignoring potential errors from input.Err()、
	//map 的迭代顺序并不确定，从实践来看，该顺序随机，每次运行都会变化
	for line, n := range counts {
		if n > 1 {
			//fmt.Printf 函数对一些表达式产生格式化输出,详细见笔记表格
			fmt.Printf("%d\t%s\n", n, line)
		}
	}
}

//!-
