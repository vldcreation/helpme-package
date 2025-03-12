package testdata

import "fmt"

func HelloWorld() {
	var n int
	fmt.Scan(&n)
	for i := 0; i < n; i++ {
		fmt.Println("Hello World")
	}
}
