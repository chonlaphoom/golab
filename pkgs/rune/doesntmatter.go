package rune

import "fmt"

func Print() {
	s := "Hi!世"
	r := []rune(s) // Convert to slice of runes (int32)

	fmt.Println("Original String:", s)
	fmt.Printf("String Byte Length: %d\n", len(s))
	fmt.Printf("Rune Slice Length:  %d\n", len(r))

	fmt.Println("---")
	// 0. Looping over the String (With Jumps)
	fmt.Println("Looping over string (Index jumps for multi-byte chars):")
	for i := range s {
		// 'i' is the byte position, and it jumps for multi-byte characters.
		fmt.Printf("  Byte Index: %d, Character: %c\n", i, s[i])
	}

	fmt.Println("---")

	// 1. Looping over the Rune Slice (No Jumps)
	fmt.Println("Looping over []rune (Index does not jump):")
	for i := range r {
		// 'i' is the character position, and it increments by 1 every time.
		fmt.Printf("  Character Index: %d, Character: %c\n", i, r[i])
	}

	fmt.Println("---")

	// 2. Direct Indexing Access
	// We can reliably get the 4th character ('世') using index 3.
	fourthChar := r[3]
	fmt.Printf("The character at RUNE index 3 is: %c\n", fourthChar)
}

// Rune basically int32 representing a Unicode code point. Non-ASCII characters may take multiple bytes in UTF-8 encoding.
