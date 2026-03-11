package main

import (
	"fmt"
	"unsafe"
)

func main() {
	numbers := make([]int, 3)
	fmt.Printf("numbers: %v, len: %d, cap: %d\n", numbers, len(numbers), cap(numbers))
	for i := range numbers {
		un := unsafe.Pointer(&numbers[i])
		fmt.Println("address of numbers[", i, "]: ", un)
	}

	fmt.Println("append numbers exceed the capacity")
	// append numbers exceeding the capacity
	// expect the new slice to be allocated and the old slice to be copied to the new slice
	copied_numbers := append(numbers, 4, 5)
	fmt.Printf("copied_numbers %v, len: %d, cap: %d\n", copied_numbers, len(copied_numbers), cap(copied_numbers))
	for i := range copied_numbers {
		un := unsafe.Pointer(&copied_numbers[i])
		fmt.Println("address of numbers[", i, "]: ", un)
	}

	numbers[0] = 99
	// expect copied_numbers[0] to be the same not changed because the old slice is copied to the new slice
	fmt.Printf("numbers at index 0: %d\n", numbers[0])
	fmt.Printf("copied_numbers at index 0: %d\n", copied_numbers[0])

	same_numbers := numbers[:]
	fmt.Printf("same_numbers: %v, len: %d, cap: %d\n", same_numbers, len(same_numbers), cap(same_numbers))

	*(*int)(unsafe.Pointer(&same_numbers[0])) = 88

	fmt.Printf("numbers at index 0: %d\n", numbers[0])
	fmt.Printf("same_numbers at index 0: %d\n", same_numbers[0])
}
