package utils

func GetFibonacciNumber(index int) int {
	prevNum := 0
	currentNum := 1
	for i := 0; i < index; i++ {
		prevNum, currentNum = currentNum, prevNum+currentNum
	}
	return currentNum
}
