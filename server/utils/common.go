package utils

import (
	"github.com/RoaringBitmap/roaring/v2"
	"github.com/RoaringBitmap/roaring/v2/roaring64"
	"sync"
)

func Merge32(a, b []uint32) []uint32 {
	set := roaring.New()
	set.AddMany(a)
	set.Or(roaring.BitmapOf(b...))
	return set.ToArray()
}

func Merge(a, b []uint64) []uint64 {
	set := roaring64.New()
	set.AddMany(a)
	set.Or(roaring64.BitmapOf(b...))
	return set.ToArray()
}

func Diff32(a, b []uint32) (onlyA, both, onlyB []uint32) {
	setA := roaring.BitmapOf(a...)
	setB := roaring.BitmapOf(b...)

	return roaring.AndNot(setA, setB).ToArray(),
		roaring.And(setA, setB).ToArray(),
		roaring.AndNot(setB, setA).ToArray()
}

func Diff(a, b []uint64) (onlyA, both, onlyB []uint64) {
	setA := roaring64.BitmapOf(a...)
	setB := roaring64.BitmapOf(b...)

	return roaring64.AndNot(setA, setB).ToArray(),
		roaring64.And(setA, setB).ToArray(),
		roaring64.AndNot(setB, setA).ToArray()
}

// WithLock 是一个辅助函数，用于在特定代码块内自动锁定和解锁
func WithLock(mu *sync.Mutex, fn func()) {
	mu.Lock()
	defer mu.Unlock()
	fn()
}

func WithTryLock(mu *sync.Mutex, fn func()) {
	if mu.TryLock() {
		defer mu.Unlock()
		fn()
	}
}

func Ptr[T any](v T) *T {
	return &v
}
