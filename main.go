package main

import (
	"fmt"
	"os"
	"time"
)

type Cache interface {
	Get(k string)
	Set(k, v string)
	DeleteOldKeys()
}

type CacheItem struct {
	value      string
	expiration time.Time
}

type cacheImpl struct {
	data map[string]CacheItem
}

func newCacheImpl() *cacheImpl {
	return &cacheImpl{
		data: make(map[string]CacheItem),
	}
}

func (c *cacheImpl) Get(k string) {
	item, ok := c.data[k]
	if !ok {
		fmt.Println("Key not found in cache")
		return
	}

	if time.Now().After(item.expiration) {
		delete(c.data, k)
		fmt.Println("Key deleted as it has expired")
	} else {
		fmt.Println("Value:", item.value)
		delete(c.data, k)
		fmt.Println("Key deleted after retrieval")
	}
}

func (c *cacheImpl) Set(k, v string) {
	expiration := time.Now().Add(1 * time.Minute)
	c.data[k] = CacheItem{value: v, expiration: expiration}
	fmt.Println("Value set in cache")
}

func (c *cacheImpl) DeleteOldKeys() {
	for k, item := range c.data {
		if time.Now().After(item.expiration) {
			delete(c.data, k)
			fmt.Println("Deleted expired key:", k)
		}
	}
}

type dbImpl struct {
	cache Cache
}

func newDbImpl(cache Cache) *dbImpl {
	return &dbImpl{
		cache: cache,
	}
}

func (d *dbImpl) Get(k string) {
	d.cache.Get(k)
	fmt.Println("Value not found in database")
}

func main() {
	cache := newCacheImpl()
	db := newDbImpl(cache)

	go func() {
		for {
			cache.DeleteOldKeys()
			time.Sleep(30 * time.Second)
		}
	}()

	for {
		fmt.Println("Menu:")
		fmt.Println("1. Enter key and value")
		fmt.Println("2. Enter key to get value")
		fmt.Println("3. Exit")

		var choice int
		_, err := fmt.Scanln(&choice)
		if err != nil {
			fmt.Println("Invalid choice. Please try again.")
			continue
		}

		switch choice {
		case 1:
			var key, value string
			fmt.Println("Enter key:")
			fmt.Scanln(&key)
			fmt.Println("Enter value:")
			fmt.Scanln(&value)
			cache.Set(key, value)
			fmt.Println("Key and value added to database")

		case 2:
			var key string
			fmt.Println("Enter key to get value:")
			fmt.Scanln(&key)
			db.Get(key)

		case 3:
			fmt.Println("Exiting the program...")
			os.Exit(0)

		default:
			fmt.Println("Invalid choice. Please try again.")
		}
	}
}
