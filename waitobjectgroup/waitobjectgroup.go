package waitobjectgroup

import (
	"github.com/rs/xid"
)

// idの型をwrap
type WaitObjectID xid.ID

// idをkeyとしてchを保持, userには渡さない
type WaitObjectGroup struct {
	chMap map[WaitObjectID](chan struct{})
}

// goroutineの実行と、chリスト(≒groutineリスト)への登録
// closeしたらidをdelete
// ユーザのch操作は許容しないためidのみ返す
func (wog *WaitObjectGroup) Go(f func()) WaitObjectID {

	if len(wog.chMap) == 0 {
		wog.chMap = make(map[WaitObjectID](chan struct{}))
	}

	id := WaitObjectID(xid.New())
	done := make(chan struct{})
	wog.chMap[id] = done

	go func() {
		f()
		close(done)
		delete(wog.chMap, id)
	}()

	return id
}

// 引数で受け取ったchに対応するgroutineが全て終わるまで待機
func (wog *WaitObjectGroup) Wait(idList ...WaitObjectID) {
	for _, id := range idList {
		//closeとdelete(id)の時差やwaitは呼ばれるまでの時差を考慮
		ch, found := wog.chMap[id]
		if found {
			<-ch
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
