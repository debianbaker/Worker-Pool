package main

import (
	"fmt"
)

func main(){
	factory := make(chan bool)
	toys := make(chan int, 15)
	
	for i := 1; i <= 10; i++{
		if i == 7{
			close(factory)
		}
		select{
			case v := <-factory: // not ready as there is no sender, and the channel is open until i=6.
				fmt.Println("Taken out", v)
				return
			case toys <- i:     // not ready when the buffer becomes full.
			fmt.Printf("Toy %d produced\n", i)
			//select randomly selects if all the statements are ready.
			//A select blocks until one of its cases can run, then it executes that case.
		}
	}
}