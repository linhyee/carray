package carray

import (
	"errors"
	"fmt"
	"sync/atomic"
)

// ConcurrencyArray 代表并发安全的整数数组接口
type ConcurrencyArray interface {
	// Set 用于设置指定场比索引上的元素值
	Set(index uint32, elem int) (err error)
	// Get 用于获取指定索引上的元素值
	Get(index uint32) (elem int, err error)
	// Len 用于获取数组的长度
	Len() uint32
}

// intArray 代表ConcurrencyArray接口的实现类型
type intArray struct {
	length uint32
	val    atomic.Value
}

// NewConcurrencyArray 创建一个ConcurrencyArray类型值
func NewConcurrencyArray(length uint32) ConcurrencyArray {
	array := intArray{}
	array.length = length
	array.val.Store(make([]int, array.length))
	return &array
}

// Set 用于设置指定场比索引上的元素值
func (array *intArray) Set(index uint32, elem int) (err error) {
	if err = array.checkIndex(index); err != nil {
		return
	}
	if err = array.checkValue(); err != nil {
		return
	}

	// Don't do this. 会导致形成竞态条件,因类slice是引用类型,
	// 对其复制并不会复制底层数组
	//oldArray := array.val.Load().([]int)
	//oldArray[index] = elem
	//array.val.Store(oldArray)

	//写时复制(cow)
	newArray := make([]int, array.length)
	copy(newArray, array.val.Load().([]int))
	newArray[index] = elem
	array.val.Store(newArray)
	return
}

// Get 用于获取指定索引上的元素值
func (array *intArray) Get(index uint32) (elem int, err error) {
	if err = array.checkIndex(index); err != nil {
		return
	}
	if err = array.checkValue(); err != nil {
		return
	}
	elem = array.val.Load().([]int)[index]
	return
}

// Len 用于获取数组的长度
func (array *intArray) Len() uint32 {
	return array.length
}

// checkIndex 用于检查索引的有效性
func (array *intArray) checkIndex(index uint32) error {
	if index >= array.length {
		return fmt.Errorf("index out of range [0,%d)", array.length)
	}
	return nil
}

// checkValue 用于检查原子值是否已存在有值
func (array *intArray) checkValue() error {
	v := array.val.Load()
	if v == nil {
		return errors.New("invalid int array")
	}
	return nil
}
