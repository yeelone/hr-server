package formula

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/xcltapestry/xclpkg/algorithm"
)

var valueMap map[string]float64

func Resolve(formula string, values map[string]float64) float64 {

	valueMap = values
	exp, err := ExpConvert(covertStr(formula))
	if err != nil {
		fmt.Println("中序表达式转后序表达式失败! ", err)
	} else {
		return Exp(exp)
	}
	return 0.00
}

func covertStr(str string) (sarr []string) {
	//先去除所有的空格
	str = strings.Replace(str, " ", "", -1)
	//同时去除换行符
	str = strings.Replace(str, "\n", "", -1)
	word := ""

	//在当前系统中，[]中括号包括的内容属于公式依赖的对应值，如果在[]存在着 () 小括号，不能把中括号中的小括号当成表达式运算符

	key := ""
	isKey := false
	for i, s := range str {
		ch := string(s)
		if ch == "[" {
			key += ch
			isKey = true
			continue
		}
		if ch == "]" {
			key += ch
			sarr = append(sarr, key)
			key = ""
			isKey = false
			continue
		}

		if isKey {
			key += ch
			continue
		}

		if IsOperator(ch) {
			if word != "" {
				sarr = append(sarr, word)
			}
			sarr = append(sarr, ch)
			word = ""
			continue
		}

		word += ch

		if i == len(str)-1 {
			sarr = append(sarr, word)
		}

	}

	return sarr
}

func ExpConvert(str []string) ([]string, error) {

	var result []string
	stack := algorithm.NewStack()
	for _, s := range str {
		ch := s
		if IsOperator(ch) { //是运算符
			if stack.Empty() || ch == "(" {
				stack.Push(ch)
			} else {
				if ch == ")" { //处理括号
					for {
						if stack.Empty() {
							return nil, errors.New("表达式有问题! 没有找到对应的\"(\"号")
						}
						if stack.Top().(string) == "(" {
							break
						}
						result = append(result, stack.Top().(string))
						stack.Pop()
					}

					//弹出"("
					stack.Top()
					stack.Pop()
				} else { //非括号
					//比较优先级
					for {
						if stack.Empty() {
							break
						}
						m := stack.Top().(string)
						if Priority(ch) > Priority(m) {
							break
						}
						// result += m
						result = append(result, m)
						stack.Pop()
					}
					stack.Push(ch)
				}
			}

		} else { //非运算符
			// result += ch

			result = append(result, ch)
		} //end IsOperator()

	} //end for range str

	for {
		if stack.Empty() {
			break
		}
		// result += stack.Top().(string)
		result = append(result, stack.Top().(string))
		stack.Pop()
	}

	return result, nil
}

func Exp(str []string) float64 {
	stack := algorithm.NewStack()
	for _, s := range str {
		ch := s
		if IsOperator(ch) { //是运算符
			if stack.Empty() {
				break
			}
			//stack.Print()
			b := stack.Top().(string)
			stack.Pop()

			a := stack.Top().(string)
			stack.Pop()

			// ia,ib := convToInt32(a,b)

			sv := fmt.Sprintf("%f", Calc(ch, a, b))

			// fmt.Println("Exp() => ", a, "", ch, "", b, "=", sv)

			stack.Push(sv)
			// stack.Print()
		} else {
			stack.Push(ch)
			// stack.Print()

		} //end IsOperator
	}
	// stack.Print()
	if !stack.Empty() {
		result := stack.Top().(string)
		value, err := strconv.ParseFloat(result, 32)
		if err != nil {
			// do something sensible
		}
		stack.Pop()
		return value
	}
	return 0.00
}

func IsOperator(op string) bool {
	switch op {
	case "(", ")", "+", "-", "*", "/":
		return true
	default:
		return false
	}
}

func Priority(op string) int {
	switch op {
	case "*", "/":
		return 3
	case "+", "-":
		return 2
	case "(":
		return 1
	default:
		return 0
	}
}

func covert2Int(a string) int {
	i, _ := strconv.Atoi(a)
	return i
}

func covert2Float(a string) float64 {
	i, _ := strconv.ParseFloat(a, 64)
	return i
}

func Calc(op string, a, b string) float64 {
	var ia = -1.00
	var ib = -1.00

	if val, ok := valueMap[a]; ok {
		ia = val
	}

	if val, ok := valueMap[b]; ok {
		ib = val
	}
	if ia == -1 {
		ia = covert2Float(a)
	}

	if ib == -1 {
		ib = covert2Float(b)
	}
	switch op {
	case "*":
		return (ia * ib)
	case "/":
		return (ia / ib)
	case "+":
		return (ia + ib)
	case "-":
		return (ia - ib)
	default:
		return 0
	}
}
