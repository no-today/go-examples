package main

import (
	"fmt"
	"reflect"
	"testing"
)

func TestSum(t *testing.T) {
	t.Run("计算切面内元素的和", func(t *testing.T) {
		numbers := []int{1, 2, 3, 4, 5}
		want := 15
		got := Sum(numbers)

		assertEquals(t, got, want)
	})

	t.Run("空的切面", func(t *testing.T) {
		var numbers []int
		want := 0
		got := Sum(numbers)

		assertEquals(t, got, want)
	})
}

func TestSumAll(t *testing.T) {
	t.Run("分别计算每个切面内元素的和", func(t *testing.T) {
		numbersToSum := [][]int{
			{1, 2, 3},
			{1, 2, 3, 4},
			{1, 2, 3, 4, 5},
		}

		want := []int{6, 10, 15}
		got := SumAll(numbersToSum...)

		assertDeepEquals(t, got, want)
	})

	t.Run("空的切面", func(t *testing.T) {
		numbersToSum := [][]int{{}, {}, {}}

		want := []int{0, 0, 0}
		got := SumAll(numbersToSum...)

		assertDeepEquals(t, got, want)
	})
}

func TestSumByIndexRange(t *testing.T) {
	t.Run("计算切片指定索引范围内元素的和", func(t *testing.T) {
		numbers := []int{1, 2, 3, 4, 5}
		start := 1
		end := 3

		want := 5
		got := SumByIndexRange(start, end, numbers)

		assertEquals(t, got, want)
	})

	t.Run("错误的索引范围:起始索引小于零", func(t *testing.T) {
		numbers := []int{1, 2, 3, 4, 5}
		start := -1
		end := 2

		want := 0
		got := SumByIndexRange(start, end, numbers)

		assertEquals(t, got, want)
	})

	t.Run("错误的索引范围:结束索引大于切面长度", func(t *testing.T) {
		numbers := []int{1, 2, 3, 4, 5}
		start := 1
		end := 100

		want := 0
		got := SumByIndexRange(start, end, numbers)

		assertEquals(t, got, want)
	})
}

func assertEquals(t *testing.T, got int, want int) {
	if got != want {
		t.Errorf("want: %v, got: %v", want, got)
	}
}

func assertDeepEquals(t *testing.T, got []int, want []int) {
	if !reflect.DeepEqual(got, want) {
		t.Errorf("want: %v, got: %v", want, got)
	}
}

func TestSlice(t *testing.T) {
	slice0 := []int{1, 2, 3}

	// 两个指针指向同一块内存
	slice1 := slice0[:]

	slice2 := slice0

	println(slice0)
	println(slice1)
	println(slice2)

	slice0[0] = 3

	fmt.Printf("%v\n", slice0)
	fmt.Printf("%v\n", slice1)
	fmt.Printf("%v\n", slice2)

	if !reflect.DeepEqual(slice0, slice1) {
		t.Errorf("预期相同,但是没有")
	}
}
