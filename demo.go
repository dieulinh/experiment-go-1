package main

import (
	"fmt"
)

func main() {
	intArr := [...]int32{1, 2, 3}
	fmt.Println(intArr)
	fmt.Println("The length is %v, capacity is %v", len(intArr), cap(intArr))
	var intSlice []int32 = []int32{4, 5, 6}
	fmt.Println("The length is %v, capacity is %v", len(intSlice), cap(intSlice))
	intSlice = append(intSlice, 7)
	fmt.Println("The length is %v, capacity is %v", len(intSlice), cap(intSlice))
	fmt.Println(intSlice)
	var newIntArr []int32 = []int32{10, 11}
	intSlice = append(intSlice, newIntArr...)
	fmt.Println("The length is %v, capacity is %v", len(intSlice), cap(intSlice))
	intSlice = append(intSlice, 17)
	fmt.Println("The length is %v, capacity is %v", len(intSlice), cap(intSlice))
	var c = make(chan int)
	go process(c)
	for i := range c {
		fmt.Println(i)
	}

}
func process(c chan int) {
	for i := 0; i < 5; i++ {
		c <- i
	}
	close(c)
}
