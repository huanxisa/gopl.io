// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

// See page 4.
//!+

// Echo1 prints its command-line arguments. 注释和其他语言一致，无须在意
package main

import (
	"fmt"
	"os"
)

func main() {
	//注意这里的变量声明，这个地方和其他语言不太一样，多了一个var，并且会在变量后面写上类型
	//变量会在声明时直接初始化。如果变量没有显式初始化，则被隐式地赋予其类型的 零值（zero value）
	var s, sep string
	//符号 := 是 短变量声明
	//所以 j=i++ 非法，而且 ++ 和 -- 都只能放在变量名后面，因此 --i 也非法
	//for这三个部分都可以省略，全部省略即为无限循环，同时for还有遍历效果，可见echo2
	for i := 0; i < len(os.Args); i++ {
		s += sep + os.Args[i]
		sep = " "
	}
	fmt.Println(s)
}

//!-
