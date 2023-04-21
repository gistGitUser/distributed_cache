package main

import (
	"fmt"
	"log"
	"net"
	"sync"
	"testing"
	"time"
)

func benchRedisClient() {
	startPopulating := time.Now()

	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		go func(i int) {

			conn, err := net.Dial("tcp", ":5000")
			if err != nil {
				log.Fatal(err)
			}

			for j := 0; j < 10; j++ {
				conn.Write([]byte("SET Foo Bar 2500"))
			}
			wg.Done()
		}(i)
		wg.Add(1)
	}
	wg.Wait()
	fmt.Printf("Populating db took %s\n", time.Since(startPopulating))

}

// https://github.com/dragonflydb/dragonfly/issues/81
func TestLoad(t *testing.T) {

	benchRedisClient()

}
