package waitobjectgroup

type WaitObjectGroup struct {
	chanList []chan struct{}
}

// goroutineの実行と、chリスト(≒groutineリスト)への登録
func (wog *WaitObjectGroup) Go(f func()) chan struct{} {
	done := make(chan struct{})
	go func() {
		f()
		close(done)
	}()
	wog.chanList = append(wog.chanList, done)
	return done
}

// 引数で受け取ったchに対応するgroutineが全て終わるまで待機
func (wog *WaitObjectGroup) Wait(chList ...chan struct{}) {
	for _, ch := range chList {
		<-ch
	}
}

// List内の全てがcloseするまで待機
func (wog *WaitObjectGroup) WaitAll() {
	wog.Wait(wog.chanList...)
}
