// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

// See page 10.
//!+

// Dup2 prints the count and text of lines that appear more than once
// in the input.  It reads from stdin or from a list of named files.
package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	//map 是一个由 make 函数创建的数据结构的引用。map 作为参数传递给某函数时比如在下面的countLines方法，该函数接收这个引用的一份拷贝
	counts := make(map[string]int)
	files := os.Args[1:]
	if len(files) == 0 {
		countLines(os.Stdin, counts)
	} else {
		for _, arg := range files {
			f, err := os.Open(arg)
			if err != nil {
				//与表示任意类型默认格式值的动词 %v，向标准错误流打印一条信息
				fmt.Fprintf(os.Stderr, "dup2: %v\n", err)
				continue
			}
			//countLines 函数在其声明前被调用，函数和包级别的变量（package-level entities）可以任意顺序声明，并不影响其被调用
			countLines(f, counts)
			f.Close()
		}
	}
	for line, n := range counts {
		if n > 1 {
			fmt.Printf("%d\t%s\n", n, line)
		}
	}
}

func countLines(f *os.File, counts map[string]int) {
	//bufio.NewScanner除了能够读取标准输入，也能读取文件流。
	input := bufio.NewScanner(f)
	for input.Scan() {
		counts[input.Text()]++
	}
	// NOTE: ignoring potential errors from input.Err()
}

//!-
