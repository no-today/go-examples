package main

// Sum 计算切片内元素的和
func Sum(numbers []int) (sum int) {
	for _, number := range numbers {
		sum += number
	}
	return sum
}

// SumAll 分别计算每个切片内元素的和
func SumAll(numbersToSum ...[]int) (sums []int) {
	for _, numbers := range numbersToSum {
		sums = append(sums, Sum(numbers))
	}
	return sums
}

// SumByIndexRange 计算索引范围内元素的和。不包含结束索引
func SumByIndexRange(start, end int, numbers []int) (sum int) {
	if start < 0 || end > len(numbers) {
		return 0
	}

	for i := start; i < end; i++ {
		sum += numbers[i]
	}
	return sum
}
