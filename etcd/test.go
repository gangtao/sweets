package main

import (
	"time"
	"context"
	"fmt"
	"github.com/coreos/etcd/clientv3"
)

func main() {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"localhost:3000"},
		// Endpoints: []string{"localhost:2379", "localhost:22379", "localhost:32379"}
		DialTimeout: 5 * time.Second,
	})

	if err != nil {
		panic(err)
	}

	fmt.Println("connected to etcd - ")

	defer cli.Close()

	kv := clientv3.NewKV(cli)

	// put a key value
	putResp, err := kv.Put(context.TODO(),"/test/key1", "Hello etcd!")
	if err != nil {
		panic(err)
	}
	fmt.Printf("%+v\n",putResp)

	// get a key value
	getResp, err := kv.Get(context.TODO(), "/test/key1")
	if err != nil {
		panic(err)
	}
	fmt.Printf("%+v\n",getResp)
	for _, ev := range getResp.Kvs {
		fmt.Printf("%s : %s\n", ev.Key, ev.Value)
	}

	// put a key value to the same key
	putResp, err = kv.Put(context.TODO(),"/test/key1", "Hello etcd again!")
	if err != nil {
		panic(err)
	}
	fmt.Printf("%+v\n",putResp)

	// get key value again
	getResp, err = kv.Get(context.TODO(), "/test/key1")
	if err != nil {
		panic(err)
	}
	fmt.Printf("%+v\n",getResp)
	for _, ev := range getResp.Kvs {
		fmt.Printf("%s : %s\n", ev.Key, ev.Value)
	}

	// delete the key here
	deleteResp, err := kv.Delete(context.TODO(), "/test/key1")
	if err != nil {
		panic(err)
	}
	fmt.Printf("%+v\n",deleteResp)
	fmt.Printf("key deleted %d\n",deleteResp.Deleted)

	// put a key value to the same key for watch
	putResp, err = kv.Put(context.TODO(),"/test/key1", "Hello etcd again!")
	if err != nil {
		panic(err)
	}
	fmt.Printf("%+v\n",putResp)

	watchCh := cli.Watch(context.TODO(), "/test/key1")

	go func() {
        for res := range watchCh {
            value := res.Events[0].Kv.Value
            if err != nil {
				fmt.Println("Failed to watch")
                continue
            }
            fmt.Println("now", time.Now(), "value", string(value))
        }
	}()
	
	time.Sleep(time.Duration(60)*time.Second)
	
}