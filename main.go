package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"strings"
)

func processJSON(data interface{}, prefix string, index int, fieldOrder []string) {
	switch v := data.(type) {
	case map[string]interface{}:
		for _, key := range fieldOrder {
			if value, ok := v[key]; ok {
				fullKey := fmt.Sprintf("%s%s,", prefix, key)
				processJSON(value, fullKey, index, fieldOrder)
			}
		}
	case []interface{}:
		for i, item := range v {
			fullKey := fmt.Sprintf("%s[第%d行],", prefix, i)
			processJSON(item, fullKey, i, fieldOrder)
		}
	default:
		printValue(prefix, data, index)
	}
}

func printValue(prefix string, value interface{}, index int) {
	switch v := value.(type) {
	case float64:
		str := fmt.Sprintf("%v", v)
		if strings.Contains(str, "e+") {
			if v > float64(int(v)) {
				fmt.Printf("%s%v\n", prefix, strconv.FormatFloat(v, 'f', 2, 64))
			} else {
				fmt.Printf("%s%d\n", prefix, int(v))
			}
		} else {
			fmt.Printf("%s%v\n", prefix, value)
		}
	default:
		fmt.Printf("%s%v\n", prefix, value)
	}
}

func clearTerminal() error {
	cmd := exec.Command("cmd", "/c", "cls") // 或者 "cls"（Windows 系统）
	cmd.Stdout = os.Stdout
	return cmd.Run()
}

func getFieldOrder(data interface{}) []string {
	fieldSet := make(map[string]bool)
	var fieldOrder []string

	var traverse func(interface{})
	traverse = func(value interface{}) {
		switch v := value.(type) {
		case map[string]interface{}:
			for key := range v {
				fieldSet[key] = true
				traverse(v[key])
			}
		case []interface{}:
			for _, item := range v {
				traverse(item)
			}
		}
	}

	traverse(data)

	for key := range fieldSet {
		fieldOrder = append(fieldOrder, key)
	}

	sort.Strings(fieldOrder)
	return fieldOrder
}

func main() {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("\n请输入 JSON 字符串> ")
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
			fmt.Fprintln(os.Stderr, "解析 JSON 数据失败:", err)
			continue
		}

		fieldOrder := getFieldOrder(data)

		// 将 HEAD 字段移到第一个位置
		headIndex := -1
		for i, field := range fieldOrder {
			if field == "HEAD" {
				headIndex = i
				break
			}
		}
		if headIndex > 0 {
			fieldOrder[0], fieldOrder[headIndex] = fieldOrder[headIndex], fieldOrder[0]
		}

		processJSON(data, "", -1, fieldOrder)

		fmt.Println("===================================")
	}
}
