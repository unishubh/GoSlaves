package slaves

import (
	"testing"
	"time"
)

func BenchmarkLimitedSmallRun(b *testing.B) {
	b.ReportAllocs()

	sp := MakePool(10)

	var count = 0

	sp.Open(func(obj interface{}) interface{} {
		time.Sleep(time.Millisecond * 1000)
		return 1
	}, func(cw interface{}) {
		count += cw.(int)
	})
	defer func() {
		sp.Close()
		if count != 10 {
			panic("bad count")
		}
	}()

	for i := 0; i < 10; i++ {
		sp.SendWork(i)
	}
}
