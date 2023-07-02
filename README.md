# WaitObjectGroup
株式会社ナレッジワークさんの3daysインターン
「Enablement Internship for Gophers 」
で開発した、Goの並行処理を使ったOSSを開発するという課題の成果物です。

# 概要
`sync.waitGroup`や`errorGroup`のような、ゴルーチンの終了の待ち合わせをしたい時に使えます。

実行開始時に返り値としてオブジェクトを受け取り、そのオブジェクトをWait関数の引数に渡すことで、渡したオブジェクトに対応するゴルーチン処理全てが終了するのを待機することができます。

# 使用方法
```
go install github.com/n-hizume/waitobjectgroup
```
を実行することで、このライブラリをinstallできます。


## 使用例1

```
import "github.com/n-hizume/waitobjectgroup"

func main() {
    var wog waitobjectgroup.WaitObjectGroup

    hoge1 := wog.Go(func(){...})
    hoge2 := wog.Go(func(){...})
    hoge3 := wog.Go(func(){...})
    hoge4 := wog.Go(func(){...})
    hoge5 := wog.Go(func(){...})

    wog.Wait(hoge1)
    fmt.Println("Wait1 Finished")

    wog.Wait(hoge2, hoge4)
    fmt.Println("Wait2&4 Finished")

    wog.WaitAll()
    fmt.Println("WaitAll Finished")
}
```

このように `WaitObjectGroup`の`Go()`関数の引数に`func型`を渡すことでgoroutineとして実行され、返り値としてオブジェクト(`WaitObjectID型`)を受け取ります。

このオブジェクトを`Wait()`関数に一つ以上渡すことで、対応する全てのゴルーチンの処理の終了を待機することができ、`sync.WaitGroup`で行うような待ち合わせ処理が簡単に行えます。

また、`Go()`で実行した全ての処理を待ちたい場合は、`WaitAll()`関数を使うことができます。

## 使用例2
`errorGroup.WithCountext` のように、panicが起きた時にコンテキストをキャンセルすることもできます。

下の例のように、`waitobjectgroup.CreateGroup(ctx)`として初期化してください。

これにより、`Go()`で実行中のゴルーチンのどれかでpanicが起きると、`<- ctx.Done()`でキャンセルを知ることができます。

```
func main() {
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

    wog.Go(func(){ panic("hogePanic") })

    wog.WaitAll()
}
```

## 使用例3

下の例のように、`Go()`関数自体をゴルーチン内で呼び出すこともできます。
`WaitObjectGroup`の内部情報が適切に排他制御されているためです。

ただし、下の例では`WaitAll()`が処理を待たずに終了する可能性があることに注意してください。
ゴルーチン内の`Go()`関数が実行されるよりも先に、`WaitAll()`が実行される可能性があるためです。

```
func main() {
    wog, _ := CreateGroup(context.TODO())

    go wog.Go(func(){...})
    go wog.Go(func(){...})

    wog.WaitAll()
}
```
