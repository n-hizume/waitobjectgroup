package main

import (
	"context"
	"fmt"
	"time"

	"github.com/n-hizume/go_new_wait_group/waitobjectgroup"
)

func main() {
	// sample1()
	sample2()
}

func sample1() {
	var wog waitobjectgroup.WaitObjectGroup

	f := func(i int) func() {
		return func() {
			fmt.Printf("start: %v\n", i)
			time.Sleep(time.Duration(i*100) * time.Millisecond)
			fmt.Printf("end: %v\n", i)
		}
	}

	hoge1 := wog.Go(f(1))
	hoge2 := wog.Go(f(2))
	wog.Go(f(3))
	hoge4 := wog.Go(f(4))
	wog.Go(f(5))

	wog.Wait(hoge1)
	fmt.Println("Wait1 Finished")
	wog.Wait(hoge4, hoge2)
	fmt.Println("Wait2&4 Finished")
	wog.WaitAll()
	fmt.Println("WaitAll Finished")
}

func sample2() {
	wog, ctx := waitobjectgroup.CreateGroup(context.TODO())

	wog.Go(func() {
		for {
			select {
			case <-ctx.Done():
				fmt.Println("context Cancel")
				return
			case <-time.After(100 * time.Millisecond):
				fmt.Println("Wait...")
			}
		}
	})

	wog.Go(func() {
		time.Sleep(400 * time.Millisecond)
		panic("hogePanic")
	})

	wog.WaitAll()
}
