package main

import "fmt"

func main() {
	count := 0
	for {
		c := make(chan int)
		for i := 0; i < 10; i++ {
			go func() {
				c <- 1
			}()
		}
		for i := 0; i < 10; i++ {
			<-c
		}
		count++
		if count%10000 == 0 {
			fmt.Printf("%v\n ", count)
		}
	}
}
