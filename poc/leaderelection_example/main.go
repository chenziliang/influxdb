package main

import (
	"splunk.com/etcd_study/leaderelection"
	"fmt"
	"time"
	"github.com/coreos/etcd/clientv3"
	"github.com/golang/glog"
)

type Controller struct {
	Identitiy string
}

func (c *Controller) onControllerFailover() {
	fmt.Printf("onControllerFailover=%s\n", c.Identitiy)
	time.Sleep(1000 * time.Second)
}

func (c *Controller) onControllerResignation() {
	fmt.Printf("onControllerResignation=%s\n", c.Identitiy)
}

func (c *Controller) onNewLeader(leader string) {
	fmt.Printf("onNewLeader(), Identitiy=%s, leader=%s\n", c.Identitiy, leader)
}

func (c *Controller) startup() {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"localhost:2379"},
	})
	if err != nil {
		glog.Fatal(err)
	}

	leaderelection.Startup(leaderelection.LeaderElectionConfig{
		Client: cli,
		Election: "controller",
		Identity: c.Identitiy,
		Callbacks: leaderelection.LeaderCallbacks{
			OnStartedLeading: c.onControllerFailover,
			OnStoppedLeading: c.onControllerResignation,
			OnNewLeader: c.onNewLeader,
		},
	})
}

func main() {
	c1 := Controller{
		Identitiy: "server1",
	}
	c2 := Controller{
		Identitiy: "server2",
	}
	c3 := Controller{
		Identitiy: "server3",
	}
	c4 := Controller{
		Identitiy: "server4",
	}
	go c1.startup()
	go c2.startup()
	go c3.startup()
	go c4.startup()
	time.Sleep(1000 * time.Second)
}