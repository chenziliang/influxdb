package main

import (
	"fmt"
	"time"
	"github.com/coreos/etcd/clientv3/concurrency"
	"github.com/coreos/etcd/clientv3"
	"log"
	"context"
)

func session_timeout() {
	c, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"localhost:2379"},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		log.Fatal(err)
	}

	session, err := concurrency.NewSession(c, concurrency.WithTTL(15))
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

	go print_time()

	<-session.Done()
	fmt.Println("session timeout happened")

	// We'd better reconstruct the client to avoid hang forever
	c, err = clientv3.New(clientv3.Config{
		Endpoints:   []string{"localhost:2379"},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		log.Fatal(err)
	}

	c.Close()
}

func print_time() {
	for {
		time.Sleep(1 * time.Second)
		t := time.Now().Unix()
		fmt.Println(t)
	}
}

func main() {
	session_timeout()
}
