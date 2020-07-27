package arrays

func Sum(nums []int) int {
	sum := 0
	for _, n := range nums {
		sum += n
	}
	return sum
}

func SumAllTails(numSlices ...[]int) []int {
	var sums []int
	for _, nums := range numSlices {
		if len(nums) == 0 {
			sums = append(sums, 0)
		} else {
			sums = append(sums, Sum(nums[1:]))
		}
	}
	return sums
}
