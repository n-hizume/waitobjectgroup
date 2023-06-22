package main

import (
	"fmt"
	"time"

	"github.com/n-hizume/go_new_wait_group/waitobjectgroup"
)

func main() {
	// contextをとってcancelできるようにする
	var wog waitobjectgroup.WaitObjectGroup

	/*
		他では作れないchを作る ex. 非公開objのch
		-> waitで, Goで作ったものだけを受けつけるため
		キャンセルが効くようにcontextを受け取る
		hoge1は ch or fun or id
	*/
	hoge1 := wog.Go(func() {
		fmt.Println("start hoge1")
		time.Sleep(1 * time.Second)
		fmt.Println("end hoge1")
	})
	hoge2 := wog.Go(func() {
		fmt.Println("start hoge2")
		time.Sleep(2 * time.Second)
		fmt.Println("end hoge2")
	})
	wog.Go(func() {
		fmt.Println("start hoge3")
		time.Sleep(3 * time.Second)
		fmt.Println("end hoge3")
	})
	hoge4 := wog.Go(func() {
		fmt.Println("start hoge4")
		time.Sleep(4 * time.Second)
		fmt.Println("end hoge4")
	})
	wog.Go(func() {
		fmt.Println("start hoge5")
		time.Sleep(5 * time.Second)
		fmt.Println("end hoge5")
	})

	wog.Wait(hoge1)
	fmt.Println("Wait1 Finished")
	wog.Wait(hoge4, hoge2)
	fmt.Println("Wait2 Finished")
	wog.WaitAll()
	fmt.Println("WaitAll Finished")
}
