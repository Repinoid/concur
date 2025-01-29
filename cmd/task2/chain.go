package main

import (
	"fmt"
	"time"
)

type chained struct {
	ch chan int
}

var t time.Time

func (s chained) filtering(threshold int) chained {
	chacha := chained{ch: make(chan int)}
	go func() {
		defer close(chacha.ch)
		for value := range s.ch {
			fmt.Printf("Filtering %d %v\n", value, time.Since(t).Microseconds())
			if value >= threshold {
				chacha.ch <- value
			}
		}
	}()
	return chacha
}
func (s chained) increasing(percent int) chained {
	chacha := chained{ch: make(chan int)}
	go func() {
		defer close(chacha.ch)
		for value := range s.ch {
			chacha.ch <- value * (100 + percent) / 100
			fmt.Printf("Increasing %d %v\n", value, time.Since(t).Microseconds())
		}
	}()
	return chacha
}
func (s chained) summator() chained {
	chacha := chained{ch: make(chan int)}
	go func() {
		defer close(chacha.ch)
		sum := 0
		for value := range s.ch {
			fmt.Printf("Summator %d %v\n", value, time.Since(t).Microseconds())
			sum += value
		}
		chacha.ch <- sum
		fmt.Printf("Summ %d\n", sum)
	}()
	return chacha
}
func (s chained) multip(mul float64) chained {
	chacha := chained{ch: make(chan int)}
	go func() {
		defer close(chacha.ch)
		value := int(mul * float64(<-s.ch))
		chacha.ch <- value
		fmt.Printf("multip %d %v\n", value, time.Since(t).Microseconds())
		// for value := range s.ch {
		// 	chacha.ch <- int(mul * float64(value))
		// }
	}()
	return chacha
}

func main() {
	sales := []int{10, 20, 30, 40, 50}
	threshold := 30
	percent := 10
	mul := 0.3

	t = time.Now()

	s := chained{ch: make(chan int)}

	go func() {
		defer close(s.ch)
		for _, num := range sales {
			s.ch <- num
		}
	}()

	resultCh := s.filtering(threshold).
		increasing(percent).
		summator().
		multip(mul)

	for res := range resultCh.ch {
		fmt.Print(res, " ")
	}
}
