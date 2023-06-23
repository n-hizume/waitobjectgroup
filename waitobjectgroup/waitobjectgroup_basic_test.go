package waitobjectgroup

import (
	"sync/atomic"
	"testing"
	"time"
)

func f(ms uint32, count *uint32, delta uint32) func() {
	return func() {
		time.Sleep(time.Duration(ms) * time.Millisecond)
		atomic.AddUint32(count, delta)
	}
}

func assertEqual(t *testing.T, real uint32, ideal uint32) {
	if ideal != real {
		t.Errorf("Error! {ideal: %v, real: %v}", ideal, real)
	}
}

/*
Wait関数のテスト
*/
func TestWait(t *testing.T) {
	var wog WaitObjectGroup

	var count uint32

	wo1 := wog.Go(f(10, &count, 1))
	wo2 := wog.Go(f(10, &count, 1))
	wo3 := wog.Go(f(20, &count, 3))

	wog.Wait(wo1, wo2)
	assertEqual(t, count, 2)

	wog.Wait(wo3)
	assertEqual(t, count, 5)

}

/*
Wait関数のテスト
とっくにcloseされたものをwaitしても問題ないかテスト
*/
func TestWait2(t *testing.T) {
	var wog WaitObjectGroup

	var count uint32

	wo1 := wog.Go(f(30, &count, 1))
	wo2 := wog.Go(f(10, &count, 1))
	wo3 := wog.Go(f(10, &count, 3))

	wog.Wait(wo1)
	assertEqual(t, count, 5)

	wog.Wait(wo2, wo3)
	assertEqual(t, count, 5)

}

/*
WaitAllのテスト
*/
func TestWaitAll(t *testing.T) {
	var wog WaitObjectGroup

	var count uint32

	wo1 := wog.Go(f(10, &count, 1))
	wog.Go(f(20, &count, 1))
	wog.Go(f(30, &count, 3))

	wog.Wait(wo1)
	assertEqual(t, count, 1)

	wog.WaitAll()
	assertEqual(t, count, 5)

}

/*
WaitAllのテスト
とっくにcloseされたものをwaitAllしても問題ないかテスト
*/
func TestWaitAll2(t *testing.T) {
	var wog WaitObjectGroup

	var count uint32

	wo1 := wog.Go(f(30, &count, 1))
	wog.Go(f(10, &count, 1))
	wog.Go(f(10, &count, 3))

	wog.Wait(wo1)
	assertEqual(t, count, 5)

	wog.WaitAll()
	assertEqual(t, count, 5)

}
