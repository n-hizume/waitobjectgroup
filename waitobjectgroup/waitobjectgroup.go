package waitobjectgroup

import (
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
	chMap map[WaitObjectID](WaitObject)
}

// goroutineの実行と、chリスト(≒groutineリスト)への登録
// closeしたらidをdelete
// ユーザのch操作は許容しないためidのみ返す
func (wog *WaitObjectGroup) Go(f func()) WaitObjectID {

	if len(wog.chMap) == 0 {
		wog.chMap = make(map[WaitObjectID](WaitObject))
	}

	id := WaitObjectID(xid.New())
	done := make(chan struct{})
	// wog.chMap[id] = WaitObject{done, nil}
	wog.chMap[id] = WaitObject{done}

	go func() {
		f()
		close(done)
	}()

	return id
}

// 引数で受け取ったchに対応するgroutineが全て終わるまで待機
func (wog *WaitObjectGroup) Wait(idList ...WaitObjectID) {
	for _, id := range idList {
		wo, found := wog.chMap[id]
		if found {
			<-wo.ch
			delete(wog.chMap, id)
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
