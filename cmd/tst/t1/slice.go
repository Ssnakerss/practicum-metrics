package main

import "fmt"

func main() {
	slice := make([]string, 1, 3)

	func(slice []string) {
		slice = slice[1:3]
		slice[0] = "b"
		slice[1] = "b"
		fmt.Print(len(slice))
		fmt.Print(slice)
	}(slice)
	fmt.Print(len(slice))
	fmt.Print(slice)
}
