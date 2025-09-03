package utils

func CalculateSubsidy(count int) int {
	if count >= 0 && count <= 1000 {
		return 1000 - (100 * count)
	} else {
		return -(100 * (count - 10))
	}
}
