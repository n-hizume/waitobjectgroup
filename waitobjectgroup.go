package waitobjectgroup

import (
	"context"
	"fmt"
	"sync"

	"github.com/rs/xid"
)

// idの型をwrap
type WaitObjectID xid.ID

type WaitObject struct {
	ch chan struct{}
	// p  any
}

// idをkeyとしてchを保持, userには渡さない
type WaitObjectGroup struct {
	chMap  map[WaitObjectID](WaitObject)
	cancel func(error)
	mu     sync.Mutex
}

func CreateGroup(ctx context.Context) (*WaitObjectGroup, context.Context) {
	ctx, cancel := context.WithCancelCause(ctx)
	return &WaitObjectGroup{cancel: cancel}, ctx
}

// goroutineの実行と、chリスト(≒groutineリスト)への登録
// closeしたらidをdelete
// ユーザのch操作は許容しないためidのみ返す
func (wog *WaitObjectGroup) Go(f func()) WaitObjectID {

	id := WaitObjectID(xid.New())
	done := make(chan struct{})

	wog.mu.Lock()
	if wog.chMap == nil {
		wog.chMap = make(map[WaitObjectID](WaitObject))
	}
	wog.chMap[id] = WaitObject{done}
	wog.mu.Unlock()

	go func() {
		defer func() {
			close(done)
			wog.mu.Lock()
			delete(wog.chMap, id)
			wog.mu.Unlock()
			if wog.cancel != nil {
				if p := recover(); p != nil {
					fmt.Printf("panic was caused: %v\n", p)
					wog.cancel(fmt.Errorf("panic was caused: %v", p))
				}
			}
		}()
		f()
	}()

	return id
}

// 引数で受け取ったchに対応するgroutineが全て終わるまで待機
func (wog *WaitObjectGroup) Wait(idList ...WaitObjectID) {
	for _, id := range idList {
		wo, found := wog.chMap[id]
		if found {
			<-wo.ch
		}
	}
}

// List内の全てがcloseするまで待機
func (wog *WaitObjectGroup) WaitAll() {
	idList := []WaitObjectID{}
	for id := range wog.chMap {
		idList = append(idList, id)
	}
	wog.Wait(idList...)
}
