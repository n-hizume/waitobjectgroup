package waitobjectgroup

import (
	"context"
	"fmt"
	"testing"
	"time"
)

/*
wog.Go自体をgoroutineで実行
WaitAllの時にはまだGoが登録されていないのでスルーされるが、やむなし
*/
func TestWaitWithGoroutine(t *testing.T) {
	wog, _ := CreateGroup(context.TODO())

	var count uint32

	for i := 1; i <= 10; i++ {
		v := i
		go wog.Go(f(uint32(v*10), &count, uint32(v)))
	}

	wog.WaitAll()
	t.Logf("ok!!!")
	// assertEqual(t, count, 45)
}

/*
contextによるcancelの伝搬処理のテスト
panicが起きればrecovcerされctxをcancelするため、処理は終了しない
*/
func TestWaitWithContext(t *testing.T) {
	wog, ctx := CreateGroup(context.TODO())

	wog.Go(func() {
		for {
			select {
			case <-ctx.Done():
				fmt.Println("context Cancel")
				return
			case <-time.After(1 * time.Millisecond):
				fmt.Println("Wait...")
			}
		}
	})

	wog.Go(func() {
		time.Sleep(4 * time.Millisecond)
		panic("hogePanic")
	})

	wog.WaitAll()
	t.Logf("ok: panic not caused.")
}
