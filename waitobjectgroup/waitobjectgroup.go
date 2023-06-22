package waitobjectgroup

type WaitObjectGroup struct {
	chanList []chan struct{}
}

func (wog *WaitObjectGroup) Go(f func()) chan struct{} {
	done := make(chan struct{})
	go func() {
		f()
		close(done)
	}()
	wog.chanList = append(wog.chanList, done)
	return done
}

func (wog *WaitObjectGroup) Wait(chList ...chan struct{}) {
	for _, ch := range chList {
		<-ch
	}
}

func (wog *WaitObjectGroup) WaitAll() {
	wog.Wait(wog.chanList...)
}
