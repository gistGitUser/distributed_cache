package tests

import (
	"cacher/client"
	"context"
	"fmt"
	"log"
	"testing"
)

func TestClientSet(t *testing.T) {

	client, err := client.New(":5000", client.Options{})
	if err != nil {
		log.Fatal(err)
	}

	err = client.Set(context.Background(), []byte("key"), []byte("value"), 2_500_000_000)
	if err != nil {
		log.Fatal(err)
	}

}

func TestClientGET(t *testing.T) {
	client, err := client.New(":5000", client.Options{})
	if err != nil {
		log.Fatal(err)
	}

	key, err := client.Get(context.Background(), []byte("key"))
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("KEY %q", string(key))
}
