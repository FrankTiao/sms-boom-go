package scans

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/gookit/color"
	"strconv"
	"strings"
)

func ScanPhone(scanner *bufio.Scanner) ([]string, error) {
	fmt.Print("请输入要轰炸的手机号，多个手机号时使用一个空格分隔: ")
	scanner.Scan()
	phone, err := CheckPhone(scanner.Text())
	if err != nil {
		return nil, err
	}

	return phone, nil
}

func CheckPhone(phone string) ([]string, error) {
	phoneList := strings.Split(phone, " ")

	var list []string
	for _, v := range phoneList {
		if len(v) > 0 {
			list = append(list, v)
		}
	}

	if len(list) > 0 {
		return list, nil
	}

	return []string{}, errors.New("手机号格式错误，多个手机号应使用一个空格分割")
}

func ScanFrequency(scanner *bufio.Scanner) (int, error) {
	fmt.Print("请输入循环轰炸几轮，默认1轮: ")
	scanner.Scan()
	input := scanner.Text()
	if len(input) <= 0 {
		return 1, nil
	}

	frequency, err := strconv.Atoi(input)
	if err != nil {
		return 0, err
	}

	if frequency <= 0 {
		return 0, errors.New("循环轰炸轮数不能小于1")
	}

	return frequency, nil
}

func ScanInterval(scanner *bufio.Scanner) (int, error) {
	fmt.Print("请输入每轮轰炸间隔几秒，默认10秒（每轮轰炸结束后的休息时间）: ")
	scanner.Scan()
	input := scanner.Text()
	if len(input) <= 0 {
		return 10, nil
	}

	interval, err := strconv.Atoi(input)
	if err != nil {
		return 0, err
	}

	if interval < 0 {
		return 0, errors.New("轰炸间隔不能小于0")
	}

	return interval, nil
}

func ScanCoroutine(scanner *bufio.Scanner) (int, error) {
	fmt.Print("请输入启动几个协程轰炸，默认1个（最终轰炸次数 = 协程数 * 轰炸轮数）: ")
	scanner.Scan()
	input := scanner.Text()
	if len(input) <= 0 {
		return 1, nil
	}

	count, err := strconv.Atoi(input)
	if err != nil {
		return 0, err
	}

	if count <= 0 {
		return 0, errors.New("协程数不能小于1")
	}

	if count > 128 {
		return 0, errors.New("协程数不能大于128，协程开启过多会对您的电脑造成较大压力")
	}

	return count, nil
}

func ScanHandelByStringSlice(scanner *bufio.Scanner, fun func(scanner *bufio.Scanner) ([]string, error)) []string {
	var err error
	var value []string
	for true {
		value, err = fun(scanner)
		if err == nil {
			break
		}
		color.Warn.Println(err.Error())
	}

	return value
}

func ScanHandelByInt(scanner *bufio.Scanner, fun func(scanner *bufio.Scanner) (int, error)) int {
	var err error
	var value int
	for true {
		value, err = fun(scanner)
		if err == nil {
			break
		}
		color.Warn.Println(err.Error())
	}

	return value
}

func Confirm(scanner *bufio.Scanner, phone []string, frequency, interval, coroutineCount int) bool {
	out := `
请确认是否有误：
===================================================
 手 机 号：` + strings.Join(phone, ", ") + `
===================================================
 循环轰炸：` + strconv.Itoa(frequency) + ` 轮
===================================================
 轰炸间隔：` + strconv.Itoa(interval) + ` 秒
===================================================
 协 程 数：` + strconv.Itoa(coroutineCount) + ` 个
===================================================
`
	fmt.Println()
	color.Info.Print(out)
	fmt.Print("是否开始轰炸（Y/N）: ")

	scanner.Scan()
	confirm := scanner.Text()

	fmt.Println()

	return len(confirm) == 0 || strings.ToLower(confirm) == "y" || strings.ToLower(confirm) == "yes"
}
