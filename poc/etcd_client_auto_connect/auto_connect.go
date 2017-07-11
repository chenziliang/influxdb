package main

import (
	"fmt"
	"time"
	"github.com/coreos/etcd/clientv3"
	"log"
	"context"
)

func new_client1() {
	// Construct a client object without DialTimeout.
	// Always return success since it won't try to connect to the etcd server
	_, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"localhost:2379"},
	})
	if err != nil {
		// This step will not be executed
		log.Fatal(err)
	} else {
		// always here
		fmt.Println("success")
	}
}

func new_client2() {
	// Construct a client object with DialTimeout.
	// It will try to connect to the etcd server in 10 seconds.
	c, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"localhost:2379"},
		DialTimeout: 10 * time.Second,
	})
	if err != nil {
		// If the client cannot connect to the etcd server in 10 seconds.
		log.Fatal(err)
	} else {
		// successfully connect to the etcd server in 10 seconds.
		fmt.Println("success")
	}
	c.Close()
}

func auto_connect1() {
	c, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"localhost:2379"},
	})
	if err != nil {
		// This step will not be executed
		log.Fatal(err)
	}

	fmt.Println("here1")
	// these methods will hang up until the client connect to the etcd successfully
	// no matter whether we specify the DialTimeout when we construct the client object
	_, err = c.Put(context.TODO(), "k1", "v1")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("here2")
	resp, err := c.Get(context.TODO(), "k1")
	if err != nil {
		log.Fatal(err)
	}
	for _, ev := range resp.Kvs {
		fmt.Printf("%s : %s\n", ev.Key, ev.Value)
	}
}

func auto_connect2() {
	c, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"localhost:2379"},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		// This step will not be executed
		log.Fatal(err)
	}

	fmt.Println("here1")
	time.Sleep(10 * time.Second)
	// stop etcd

	fmt.Println("here2")
	// these methods will hang up until the client connect to the etcd successfully
	// no matter whether we specify the DialTimeout when we construct the client object
	_, err = c.Put(context.TODO(), "k1", "v1")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("here3")
	resp, err := c.Get(context.TODO(), "k1")
	if err != nil {
		log.Fatal(err)
	}
	for _, ev := range resp.Kvs {
		fmt.Printf("%s : %s\n", ev.Key, ev.Value)
	}
}

func main() {
	//new_client1()
	//new_client2()
	//auto_connect1()
	auto_connect2()
}

