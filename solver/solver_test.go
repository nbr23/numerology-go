package solver

import (
	"reflect"
	"sort"
	"testing"
)

func TestDigitsToNumber(t *testing.T) {
	tests := []struct {
		digits []int
		want   int
	}{
		{[]int{1}, 1},
		{[]int{1, 2}, 12},
		{[]int{1, 2, 3}, 123},
		{[]int{0}, 0},
		{[]int{0, 1}, 1},
		{[]int{0, 0, 5}, 5},
		{[]int{9, 9, 9}, 999},
		{[]int{}, 0},
	}

	for _, tt := range tests {
		got := digitsToNumber(tt.digits)
		if got != tt.want {
			t.Errorf("digitsToNumber(%v) = %d, want %d", tt.digits, got, tt.want)
		}
	}
}

func TestReverse(t *testing.T) {
	tests := []struct {
		input []int
		want  []int
	}{
		{[]int{1, 2, 3}, []int{3, 2, 1}},
		{[]int{1, 2}, []int{2, 1}},
		{[]int{1}, []int{1}},
		{[]int{}, []int{}},
		{[]int{1, 2, 3, 4, 5}, []int{5, 4, 3, 2, 1}},
	}

	for _, tt := range tests {
		input := make([]int, len(tt.input))
		copy(input, tt.input)
		reverse(input)
		if !reflect.DeepEqual(input, tt.want) {
			t.Errorf("reverse(%v) = %v, want %v", tt.input, input, tt.want)
		}
	}
}

func TestNextPermutation(t *testing.T) {
	nums := []int{1, 2, 3}
	expected := [][]int{
		{1, 3, 2},
		{2, 1, 3},
		{2, 3, 1},
		{3, 1, 2},
		{3, 2, 1},
	}

	for i, want := range expected {
		if !nextPermutation(nums) {
			t.Errorf("nextPermutation returned false at iteration %d", i)
		}
		if !reflect.DeepEqual(nums, want) {
			t.Errorf("iteration %d: got %v, want %v", i, nums, want)
		}
	}

	if nextPermutation(nums) {
		t.Error("nextPermutation should return false after last permutation")
	}
}

func TestNextPermutationWithDuplicates(t *testing.T) {
	nums := []int{1, 1, 2}
	count := 1
	for nextPermutation(nums) {
		count++
	}
	if count != 3 {
		t.Errorf("expected 3 permutations of [1,1,2], got %d", count)
	}
}

func TestPermutationCount(t *testing.T) {
	digits := []int{1, 2, 3, 4}
	sort.Ints(digits)
	count := 1
	perm := make([]int, len(digits))
	copy(perm, digits)
	for nextPermutation(perm) {
		count++
	}
	if count != 24 {
		t.Errorf("expected 24 permutations of 4 distinct digits, got %d", count)
	}
}

func TestGenerateGroupings(t *testing.T) {
	tests := []struct {
		digits []int
		want   [][]int
	}{
		{
			[]int{1, 2, 3},
			[][]int{
				{1, 2, 3},
				{1, 23},
				{12, 3},
				{123},
			},
		},
		{
			[]int{1, 2},
			[][]int{
				{1, 2},
				{12},
			},
		},
		{
			[]int{1},
			[][]int{
				{1},
			},
		},
	}

	for _, tt := range tests {
		var got [][]int
		for g := range generateGroupings(tt.digits) {
			got = append(got, g)
		}
		if len(got) != len(tt.want) {
			t.Errorf("generateGroupings(%v): got %d groupings, want %d", tt.digits, len(got), len(tt.want))
			continue
		}
		for i, g := range got {
			if !reflect.DeepEqual(g, tt.want[i]) {
				t.Errorf("generateGroupings(%v)[%d] = %v, want %v", tt.digits, i, g, tt.want[i])
			}
		}
	}
}

func TestGenerateGroupingsCount(t *testing.T) {
	tests := []struct {
		n    int
		want int
	}{
		{1, 1},
		{2, 2},
		{3, 4},
		{4, 8},
		{5, 16},
		{6, 32},
	}

	for _, tt := range tests {
		digits := make([]int, tt.n)
		for i := range digits {
			digits[i] = i + 1
		}
		count := 0
		for range generateGroupings(digits) {
			count++
		}
		if count != tt.want {
			t.Errorf("generateGroupings with %d digits: got %d groupings, want %d", tt.n, count, tt.want)
		}
	}
}

func TestApplyOp(t *testing.T) {
	tests := []struct {
		a      int
		op     rune
		b      int
		want   int
		wantOk bool
	}{
		{2, '+', 3, 5, true},
		{5, '-', 3, 2, true},
		{4, '*', 3, 12, true},
		{12, '/', 3, 4, true},
		{10, '/', 0, 0, false},
		{7, '/', 3, 0, false},
		{0, '/', 5, 0, true},
		{-6, '/', 2, -3, true},
		{6, '/', -2, -3, true},
		{5, '-', 10, -5, true},
		{3, '*', -4, -12, true},
	}

	for _, tt := range tests {
		got, ok := applyOp(tt.a, tt.op, tt.b)
		if ok != tt.wantOk {
			t.Errorf("applyOp(%d, %c, %d) ok = %v, want %v", tt.a, tt.op, tt.b, ok, tt.wantOk)
		}
		if ok && got != tt.want {
			t.Errorf("applyOp(%d, %c, %d) = %d, want %d", tt.a, tt.op, tt.b, got, tt.want)
		}
	}
}

func TestEvalWithPrecedence(t *testing.T) {
	tests := []struct {
		numbers []int
		ops     []rune
		want    int
		wantOk  bool
	}{
		{[]int{1, 2, 3}, []rune{'+', '+'}, 6, true},
		{[]int{1, 2, 3}, []rune{'+', '*'}, 7, true},
		{[]int{1, 2, 3}, []rune{'*', '+'}, 5, true},
		{[]int{1, 2, 3}, []rune{'*', '*'}, 6, true},
		{[]int{2, 3, 4}, []rune{'+', '*'}, 14, true},
		{[]int{2, 3, 4}, []rune{'*', '+'}, 10, true},
		{[]int{10, 2, 3}, []rune{'-', '*'}, 4, true},
		{[]int{10, 2, 3}, []rune{'/', '-'}, 2, true},
		{[]int{1, 2, 3, 4}, []rune{'+', '*', '-'}, 3, true},
		{[]int{1, 2, 3, 4, 5, 6}, []rune{'*', '*', '*', '+', '-'}, 23, true},
		{[]int{10, 0}, []rune{'/'}, 0, false},
		{[]int{10, 3}, []rune{'/'}, 0, false},
	}

	for _, tt := range tests {
		got, ok := evalWithPrecedence(tt.numbers, tt.ops)
		if ok != tt.wantOk {
			t.Errorf("evalWithPrecedence(%v, %v) ok = %v, want %v", tt.numbers, string(tt.ops), ok, tt.wantOk)
		}
		if ok && got != tt.want {
			t.Errorf("evalWithPrecedence(%v, %v) = %d, want %d", tt.numbers, string(tt.ops), got, tt.want)
		}
	}
}

func TestBuildExprString(t *testing.T) {
	tests := []struct {
		numbers []int
		ops     []rune
		want    string
	}{
		{[]int{1, 2, 3}, []rune{'+', '-'}, "1+2-3"},
		{[]int{10, 20}, []rune{'*'}, "10*20"},
		{[]int{5}, []rune{}, "5"},
		{[]int{1, 2, 3, 4}, []rune{'+', '*', '/'}, "1+2*3/4"},
	}

	for _, tt := range tests {
		got := buildExprString(tt.numbers, tt.ops)
		if got != tt.want {
			t.Errorf("buildExprString(%v, %v) = %s, want %s", tt.numbers, string(tt.ops), got, tt.want)
		}
	}
}

func TestSolveForGrouping(t *testing.T) {
	tests := []struct {
		numbers []int
		target  int
		wantOk  bool
	}{
		{[]int{23}, 23, true},
		{[]int{20, 3}, 23, true},
		{[]int{25, 2}, 23, true},
		{[]int{99}, 23, false},
		{[]int{1, 1}, 23, false},
	}

	for _, tt := range tests {
		_, ok := solveForGrouping(tt.numbers, tt.target)
		if ok != tt.wantOk {
			t.Errorf("solveForGrouping(%v, %d) ok = %v, want %v", tt.numbers, tt.target, ok, tt.wantOk)
		}
	}
}

func TestSolve(t *testing.T) {
	tests := []struct {
		name   string
		digits []int
		target int
		wantOk bool
	}{
		{"simple sum", []int{2, 0, 3}, 23, true},
		{"needs permutation", []int{3, 2, 0}, 23, true},
		{"direct number", []int{2, 3}, 23, true},
		{"complex", []int{1, 2, 3, 4, 5, 6}, 23, true},
		{"impossible", []int{1, 9}, 23, false},
		{"zeros", []int{0, 0, 2, 3}, 23, true},
		{"all same", []int{1, 1, 1, 1}, 23, false},
		{"date format", []int{1, 1, 0, 6, 1, 5}, 23, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, ok := Solve(tt.digits, tt.target)
			if ok != tt.wantOk {
				t.Errorf("Solve(%v, %d) ok = %v, want %v", tt.digits, tt.target, ok, tt.wantOk)
			}
			if ok {
				verified := verifyExpression(result, tt.target)
				if !verified {
					t.Errorf("Solve(%v, %d) = %s does not evaluate to %d", tt.digits, tt.target, result, tt.target)
				}
			}
		})
	}
}

func TestSolveEmptyInput(t *testing.T) {
	_, ok := Solve([]int{}, 23)
	if ok {
		t.Error("Solve with empty input should return false")
	}
}

func TestSolveSingleDigit(t *testing.T) {
	_, ok := Solve([]int{5}, 23)
	if ok {
		t.Error("Solve([5], 23) should return false")
	}

	result, ok := Solve([]int{2, 3}, 23)
	if !ok {
		t.Error("Solve([2,3], 23) should find 23")
	}
	if result != "23" {
		t.Errorf("expected '23', got '%s'", result)
	}
}

func TestSolveVerifyStandardPrecedence(t *testing.T) {
	result, ok := Solve([]int{1, 2, 3, 4, 5, 6}, 23)
	if !ok {
		t.Fatal("expected solution for 123456")
	}

	if !verifyExpression(result, 23) {
		t.Errorf("expression %s does not evaluate to 23 with standard precedence", result)
	}
}

func verifyExpression(expr string, target int) bool {
	nums, ops := parseExpression(expr)
	if len(nums) == 0 {
		return false
	}
	if len(ops) == 0 {
		return nums[0] == target
	}
	result, ok := evalWithPrecedence(nums, ops)
	return ok && result == target
}

func parseExpression(expr string) ([]int, []rune) {
	var nums []int
	var ops []rune
	currentNum := 0
	hasNum := false

	for _, c := range expr {
		if c >= '0' && c <= '9' {
			currentNum = currentNum*10 + int(c-'0')
			hasNum = true
		} else if c == '+' || c == '-' || c == '*' || c == '/' {
			if hasNum {
				nums = append(nums, currentNum)
				currentNum = 0
				hasNum = false
			}
			ops = append(ops, c)
		}
	}
	if hasNum {
		nums = append(nums, currentNum)
	}
	return nums, ops
}
