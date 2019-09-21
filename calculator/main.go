package main

import (
	"fmt"
	"math"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

// Function to find precedence of
// operators.
func precedence(op string) int {
	if op == "+" || op == "-" {
		return 1
	}
	if op == "*" || op == "/" {
		return 2
	}
	if op == "^" {
		return 3
	}
	return 0
}

//Function to evaluate a simple expression with two operand and one operator
func applyOp(a float64, b float64, op string, ch chan float64) {
	switch op {
	case "+":
		ch <- a + b
		break
	case "-":
		ch <- a - b
		break
	case "*":
		ch <- a * b
		break
	case "/":
		ch <- a / b
		break
	case "^":
		ch <- math.Pow(a, b)
		break
	}
	ch <- a
}

type numstack []float64

func (s numstack) Top() float64 {
	return s[len(s)-1]
}

func (s numstack) Push(v float64) numstack {
	return append(s, v)
}

func (s numstack) Pop() (numstack, float64) {
	l := len(s)
	return s[:l-1], s[l-1]
}

func (s numstack) Empty() bool {
	return len(s) == 0
}

type strstack []string

func (s strstack) Top() string {
	return s[len(s)-1]
}

func (s strstack) Push(v string) strstack {
	return append(s, v)
}

func (s strstack) Pop() (strstack, string) {
	l := len(s)
	return s[:l-1], s[l-1]
}

func (s strstack) Empty() bool {
	return len(s) == 0
}

// Function that returns value of
// expression after evaluation.
func evaluate(exp string) float64 {
	var i int

	// stack to store integer values.
	var values numstack
	var val, val1, val2 float64

	// stack to store operators.
	var ops strstack
	var op string

	evaluate_last := func() {
		values, val2 = values.Pop()
		values, val1 = values.Pop()
		ops, op = ops.Pop()
		ch := make(chan float64, 0)
		go applyOp(float64(val1), float64(val2), op, ch)
		result := <-ch
		values = values.Push(result)
	}

	for i = 0; i < len(exp); i += 1 {
		// if value at current index is whitespace then skip it
		if string(exp[i]) == "(" {
			// if value at current index is opening brace, push it to 'ops' stack
			ops = ops.Push(string(exp[i]))
		} else if _, err := strconv.Atoi(string(exp[i])); err == nil {
			// if value at current index is a number, push it to 'values' stack
			var val float64
			// There may be more than one digit in number.
			for i < len(exp) {
				x, numErr := strconv.Atoi(string(exp[i]))
				if numErr != nil {
					break
				}
				// forming the number and pushing it to values stack
				val = (val * 10) + float64(x)
				i += 1
			}
			i -= 1
			values = values.Push(val)
		} else if string(exp[i]) == ")" {
			// if value at current index is closing brace, solve the entire expression in the brace.
			for !ops.Empty() && ops.Top() != "(" {
				evaluate_last()
			}

			// remove opening brace from 'ops' stack
			if !ops.Empty() {
				ops, _ = ops.Pop()
			}
		} else {
			// Current token is an operator.
			// While top of 'ops' has same or greater
			// precedence to current token, which
			// is an operator. Apply operator on top
			// of 'ops' to top two elements in values stack.
			for !ops.Empty() && precedence(ops.Top()) >= precedence(string(exp[i])) {
				evaluate_last()
			}

			// Push current token to 'ops' stack
			ops = ops.Push(string(exp[i]))
		}
	}

	// Entire expression has been parsed at this
	// point, apply remaining ops to remaining
	// values.
	for !ops.Empty() {
		evaluate_last()
	}

	// Value on top of the 'value' stack is the result so returning it
	values, val = values.Pop()
	return val
}

//
func replaceVars(exp string, vars map[string]int) string {
	exp = strings.Replace(exp, " ", "", -1)
	no_insert_sym := "+-*/^("
	for i := 0; i < len(exp); i += 1 {
		val, ok := vars[string(exp[i])]
		if !ok {
			if i > 0 && string(exp[i]) == "(" && !strings.Contains(no_insert_sym, string(exp[i-1])) {
				exp = exp[:i] + "*" + exp[i:]
			}
		} else if i > 0 && !strings.Contains(no_insert_sym, string(exp[i-1])) {
			exp = exp[:i] + "*" + strconv.Itoa(val) + exp[i+1:]
		} else {
			exp = exp[:i] + strconv.Itoa(val) + exp[i+1:]
		}
	}
	return exp
}

func main() {

	var variables []string
	var input string

	fmt.Scanf("%s", &input)

	isAlphabet := regexp.MustCompile(`^[A-Za-z]+$`).MatchString
	for i := 0; i < len(input); i++ {
		if isAlphabet(string(input[i])) {
			variables = append(variables, string(input[i]))
		}
	}
	sort.Strings(variables)

	ma := make(map[string]int)

	for i := 0; i < len(variables); i++ {
		if _, v := ma[variables[i]]; !v {
			var temp int
			fmt.Print("value for ", variables[i], ":")
			fmt.Scanf("%d", &temp)
			ma[variables[i]] = temp
		}
	}
	expression := replaceVars(input, ma)
	answer := evaluate(expression)
	fmt.Println(answer)
}
