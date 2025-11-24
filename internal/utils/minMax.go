package utils

func Min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func MinInSlice(nums []int) int {
	if len(nums) == 0 {
		panic("slice is empty")
	}

	min := nums[0]
	for _, num := range nums {
		if num < min {
			min = num
		}
	}
	return min
}

func Max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func MaxInSlice(nums []int) int {
	if len(nums) == 0 {
		panic("slice is empty")
	}

	max := nums[0]
	for _, num := range nums {
		if num > max {
			max = num
		}
	}
	return max
}
