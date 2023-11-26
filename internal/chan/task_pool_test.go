package _chan

import (
	"fmt"
	"testing"
	"time"
)

func TestTaskPool_Do(t1 *testing.T) {
	tp := NewTaskPool(2)
	tp.Do(func() {
		time.Sleep(time.Second)
		fmt.Println("task1")
	})

	tp.Do(func() {
		time.Sleep(time.Second)
		fmt.Println("task2")
	})

	tp.Do(func() {
		fmt.Println("task3")
		})
}
