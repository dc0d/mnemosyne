package mnemosyne_test

import (
	"fmt"
	"time"

	"github.com/dc0d/mnemosyne"
)

func ExampleCache_GetOrInit() {
	cache := mnemosyne.New()

	key := "key-1"
	val := TV{Message: "val-1"}

	result := cache.GetOrInit(key, func() (mnemosyne.TV, time.Duration) { return val, 0 })

	fmt.Println("cached value:", result)

	// performing same operation will have no effect
	// because the value is already cached.
	result = cache.GetOrInit(key, func() (mnemosyne.TV, time.Duration) { panic("will not be called") })

	fmt.Println("cached value:", result)

	// Output:
	// cached value: {val-1}
	// cached value: {val-1}
}

func ExampleCache_put_get() {
	cache := mnemosyne.New()

	key := "key-a1"
	val := TV{Message: "val-a1"}

	// at this time nothing will be found
	result, found := cache.Get(key)
	fmt.Println(result, found)

	// this operation updates the entry and overwrites it
	// even if it's not expired yet
	cache.Put(key, func() (mnemosyne.TV, time.Duration) { return val, 0 })

	result, found = cache.Get(key)
	fmt.Println(result, found)

	// Output:
	// <nil> false
	// {val-a1} true
}

type TV struct{ Message string }
