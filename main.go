package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func processJSON(data interface{}, prefix string) {
	switch v := data.(type) {
	case map[string]interface{}:
		for key, value := range v {
			fullKey := fmt.Sprintf("%s%v,", prefix, key)
			processJSON(value, fullKey)
		}
	case []interface{}:
		for i, item := range v {
			fullKey := fmt.Sprintf("%s[第%d行],", prefix, i)
			processJSON(item, fullKey)
		}
	default:
		switch data.(type) {
		case float64:
			var str = fmt.Sprint(data)
			if strings.Contains(str, "e+") {
				if data.(float64) > float64(int(data.(float64))) {
					fmt.Printf("%s%.2f\n", prefix, data)
				} else {
					fmt.Printf("%s%d\n", prefix, int(data.(float64)))
				}
			} else {
				fmt.Printf("%s%v\n", prefix, data)
			}
		default:
			fmt.Printf("%s%v\n", prefix, data)
		}
	}
}

func clearTerminal() error {
	cmd := exec.Command("cmd", "/c", "cls") // 或者 "cls"（Windows 系统）
	cmd.Stdout = os.Stdout
	return cmd.Run()
}

func main() {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("\n请输入json字符串> ")
		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			continue
		}

		input = input[:len(input)-2] // 去除换行符

		if input == "clear" {
			if err := clearTerminal(); err != nil {
				fmt.Fprintln(os.Stderr, err)
			}
			continue
		}

		var data interface{}
		err = json.Unmarshal([]byte(input), &data)
		if err != nil {
			panic(err)
		}
		fmt.Println("\n生成结果如下:\n")
		processJSON(data, "")
	}
}
