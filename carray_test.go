package carray

import (
	"math/rand"
	"sync"
	"testing"
)

func TestConcurrencyArray(t *testing.T) {
	t.Run("all", func(t *testing.T) {
		t.Run("New", testNew)
		array := NewConcurrencyArray(uint32(rand.Int31n(100)))
		maxI := uint32(1000)
		t.Run("Set", func(t *testing.T) {
			testSet(array, maxI, t)
		})
		t.Run("Get", func(t *testing.T) {
			testGet(array, maxI, t)
		})
	})
}

func testNew(t *testing.T) {
	expectedLen := uint32(rand.Int31n(1000))
	intArray := NewConcurrencyArray(expectedLen)
	if intArray == nil {
		t.Fatalf("unormal int array")
	}
	if intArray.Len() != expectedLen {
		t.Fatalf("incorrect int array length")
	}
}

func testSet(array ConcurrencyArray, maxI uint32, t *testing.T) {
	arrayLen := array.Len()
	var wg sync.WaitGroup
	wg.Add(int(maxI))
	for i := uint32(0); i < maxI; i++ {
		go func(i uint32) {
			defer wg.Done()
			for j := uint32(0); j < arrayLen; j++ {
				err := array.Set(j, int(j*i))
				if uint32(j) >= arrayLen && err == nil {
					t.Fatalf("unexpected nil error! (index: %d)", j)
				} else {
					if err != nil {
						t.Fatalf("unexpected error: %s (index: %d)", err, j)
					}
				}
			}
		}(i)
	}
	wg.Wait()
}

func testGet(array ConcurrencyArray, maxI uint32, t *testing.T) {
	arrayLen := array.Len()
	intMax := int((maxI - 1) * (arrayLen - 1))
	for i := uint32(0); i < arrayLen; i++ {
		elem, err := array.Get(i)
		if err != nil {
			t.Fatalf("unexpected error: %s (index:%d)", err, i)
		}
		if elem < 0 || elem > intMax {
			t.Fatalf("incorrect element: %d! (index: %d), expected max: %d", elem, i, intMax)
		}
	}
}
