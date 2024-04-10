package service

import (
	"fmt"
)

func GetUserInput() string {
	var input string

InputLoop:
	for {
		fmt.Println("請輸入欲查看電錶参数：(電壓/電流/用電/全部)")
		fmt.Scanln(&input)

		switch input {
		case "電壓", "電流", "用電", "全部":
			break InputLoop
		default:
			fmt.Println("輸入參數:", input, "為無效參數, 請重新輸入")
			fmt.Println("=================================")
		}
	}

	return input
}
