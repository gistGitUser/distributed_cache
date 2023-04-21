package tests

import (
	"cacher/client"
	"context"
	"fmt"
	"log"
	"sync"
	"testing"
	"time"
)

func benchRedisClient() {
	startPopulating := time.Now()

	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		go func(i int) {

			client, err := client.New(":5000", client.Options{})
			if err != nil {
				log.Fatal(err)
			}

			for j := 0; j < 1000; j++ {

				err = client.Set(context.Background(), []byte(fmt.Sprintf("%d", time.Now().Nanosecond())), []byte(fmt.Sprintf("%d", time.Now().Nanosecond())), 2_500_000_000)
				if err != nil {
					log.Println(err)
				}
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
