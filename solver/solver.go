package solver

import (
	"fmt"
	"sort"
)

var operators = []rune{'+', '-', '*', '/'}

func Solve(digits []int, target int) (string, bool) {
	sort.Ints(digits)

	return permuteAndSolve(digits, target)
}

func permuteAndSolve(digits []int, target int) (string, bool) {
	n := len(digits)
	if n == 0 {
		return "", false
	}

	perm := make([]int, n)
	copy(perm, digits)

	for {
		if result, found := solveForPermutation(perm, target); found {
			return result, true
		}
		if !nextPermutation(perm) {
			break
		}
	}
	return "", false
}

func nextPermutation(nums []int) bool {
	n := len(nums)
	i := n - 2
	for i >= 0 && nums[i] >= nums[i+1] {
		i--
	}
	if i < 0 {
		return false
	}
	j := n - 1
	for nums[j] <= nums[i] {
		j--
	}
	nums[i], nums[j] = nums[j], nums[i]
	reverse(nums[i+1:])
	return true
}

func reverse(nums []int) {
	for i, j := 0, len(nums)-1; i < j; i, j = i+1, j-1 {
		nums[i], nums[j] = nums[j], nums[i]
	}
}

func solveForPermutation(digits []int, target int) (string, bool) {
	for grouping := range generateGroupings(digits) {
		if result, found := solveForGrouping(grouping, target); found {
			return result, true
		}
	}
	return "", false
}

func generateGroupings(digits []int) <-chan []int {
	ch := make(chan []int)
	go func() {
		defer close(ch)
		generateGroupingsRecursive(digits, []int{}, ch)
	}()
	return ch
}

func generateGroupingsRecursive(remaining []int, current []int, ch chan<- []int) {
	if len(remaining) == 0 {
		result := make([]int, len(current))
		copy(result, current)
		ch <- result
		return
	}

	for prefixLen := 1; prefixLen <= len(remaining); prefixLen++ {
		num := digitsToNumber(remaining[:prefixLen])
		generateGroupingsRecursive(remaining[prefixLen:], append(current, num), ch)
	}
}

func digitsToNumber(digits []int) int {
	num := 0
	for _, d := range digits {
		num = num*10 + d
	}
	return num
}

func solveForGrouping(numbers []int, target int) (string, bool) {
	if len(numbers) == 1 {
		if numbers[0] == target {
			return fmt.Sprintf("%d", numbers[0]), true
		}
		return "", false
	}

	numOps := len(numbers) - 1
	totalCombos := 1
	for i := 0; i < numOps; i++ {
		totalCombos *= 4
	}

	for combo := 0; combo < totalCombos; combo++ {
		ops := make([]rune, numOps)
		c := combo
		for i := 0; i < numOps; i++ {
			ops[i] = operators[c%4]
			c /= 4
		}

		val, ok := evalWithPrecedence(numbers, ops)
		if ok && val == target {
			return buildExprString(numbers, ops), true
		}
	}
	return "", false
}

func evalWithPrecedence(numbers []int, ops []rune) (int, bool) {
	nums := make([]int, len(numbers))
	copy(nums, numbers)
	opsCopy := make([]rune, len(ops))
	copy(opsCopy, ops)

	for i := 0; i < len(opsCopy); {
		if opsCopy[i] == '*' || opsCopy[i] == '/' {
			val, ok := applyOp(nums[i], opsCopy[i], nums[i+1])
			if !ok {
				return 0, false
			}
			nums = append(nums[:i], append([]int{val}, nums[i+2:]...)...)
			opsCopy = append(opsCopy[:i], opsCopy[i+1:]...)
		} else {
			i++
		}
	}

	result := nums[0]
	for i, op := range opsCopy {
		val, ok := applyOp(result, op, nums[i+1])
		if !ok {
			return 0, false
		}
		result = val
	}

	return result, true
}

func buildExprString(numbers []int, ops []rune) string {
	result := fmt.Sprintf("%d", numbers[0])
	for i, op := range ops {
		result += fmt.Sprintf("%c%d", op, numbers[i+1])
	}
	return result
}

func applyOp(a int, op rune, b int) (int, bool) {
	switch op {
	case '+':
		return a + b, true
	case '-':
		return a - b, true
	case '*':
		return a * b, true
	case '/':
		if b == 0 || a%b != 0 {
			return 0, false
		}
		return a / b, true
	}
	return 0, false
}
