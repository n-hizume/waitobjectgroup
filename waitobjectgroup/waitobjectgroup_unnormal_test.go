package waitobjectgroup

import (
	"testing"
)

/*
wog.Go自体をgoroutineで実行
WaitAllの時にはまだGoが登録されていないのでスルーされるが、やむなし？
*/
func TestWaitWithGoroutine(t *testing.T) {
	var wog WaitObjectGroup

	var count uint32

	for i := 1; i <= 10; i++ {
		v := i
		go wog.Go(f(uint32(v*10), &count, uint32(v)))
	}

	wog.WaitAll()
	// assertEqual(t, count, 45)
}
