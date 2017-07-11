package main

import (
	"fmt"
	"time"
	"github.com/coreos/etcd/clientv3"
	"log"
	"context"
)

func auto_connect()  {
	c, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"localhost:2379"},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		log.Fatal(err)
	}
	defer c.Close()

	fmt.Println("before restart")

	_, err = c.Put(context.TODO(), "k1", "v1")
	if err != nil {
		log.Fatal(err)
	}

	resp, err := c.Get(context.TODO(), "k1")
	if err != nil {
		log.Fatal(err)
	}
	for _, ev := range resp.Kvs {
		fmt.Printf("%s : %s\n", ev.Key, ev.Value)
	}


	time.Sleep(20 * time.Second)
	// restart etcd

	fmt.Println("after restart")

	// This method will block until it reconnect to the server
	resp, err = c.Get(context.TODO(), "k1")
	if err != nil {
		log.Fatal(err)
	}

	for _, ev := range resp.Kvs {
		fmt.Printf("%s : %s\n", ev.Key, ev.Value)
	}
}


func main() {
	auto_connect()
}

