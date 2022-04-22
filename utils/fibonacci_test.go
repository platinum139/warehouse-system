package utils

import "testing"

func TestGetFibonacciNumberIndex0(t *testing.T) {
	fibNumber := GetFibonacciNumber(0)
	if fibNumber != 1 {
		t.Fatalf(`GetFibonacciNumber(0) = %q, want 1, error`, fibNumber)
	}
}

func TestGetFibonacciNumberIndex1(t *testing.T) {
	fibNumber := GetFibonacciNumber(1)
	if fibNumber != 1 {
		t.Fatalf(`GetFibonacciNumber(1) = %q, want 1, error`, fibNumber)
	}
}

func TestGetFibonacciNumberIndex2(t *testing.T) {
	fibNumber := GetFibonacciNumber(2)
	if fibNumber != 2 {
		t.Fatalf(`GetFibonacciNumber(2) = %q, want 2, error`, fibNumber)
	}
}

func TestGetFibonacciNumberIndex3(t *testing.T) {
	fibNumber := GetFibonacciNumber(3)
	if fibNumber != 3 {
		t.Fatalf(`GetFibonacciNumber(3) = %q, want 3, error`, fibNumber)
	}
}

func TestGetFibonacciNumberIndex4(t *testing.T) {
	fibNumber := GetFibonacciNumber(4)
	if fibNumber != 5 {
		t.Fatalf(`GetFibonacciNumber(4) = %q, want 5, error`, fibNumber)
	}
}

func TestGetFibonacciNumberIndex5(t *testing.T) {
	fibNumber := GetFibonacciNumber(5)
	if fibNumber != 8 {
		t.Fatalf(`GetFibonacciNumber(5) = %q, want 8, error`, fibNumber)
	}
}
