package config

import (
	"log"
	"time"
	"fmt"
	"os"
	"context"

	"github.com/coreos/etcd/clientv3"
)

type etcdClient struct {
	servers []string
	sessionTimeout time.Duration
}

func NewEtcdStore() KVStore {
	etcdHost := os.Getenv("ETCD_HOST")
	etcdPort := os.Getenv("ETCD_PORT")
	var hosts []string

	if etcdHost != "" {
		hostPath := fmt.Sprintf("%s:%s", etcdHost, etcdPort)
		hosts = []string{hostPath}
	} else {
		hosts = []string{"localhost:3000"}
	}
	
	var timeout = time.Second*5

	return &etcdClient{hosts, timeout}
}

func (c *etcdClient) GetConfig(dataId string, group string) (string, error) {
	log.Printf("GetConfig for dataId [%s] group [%s]", dataId, group)

	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   c.servers,
		DialTimeout: c.sessionTimeout,
	})

	if err != nil {
		panic(err)
	}
	log.Println("connected to etcd ")
	defer cli.Close()

	// get a key value
	kv := clientv3.NewKV(cli)
	itemPath := fmt.Sprintf("/%s/%s", group, dataId)
	getResp, err := kv.Get(context.TODO(), itemPath)
	if err != nil {
		panic(err)
	}
	log.Printf("%+v\n",getResp)
	for _, ev := range getResp.Kvs {
		log.Printf("%s : %s\n", ev.Key, ev.Value)
		// TODO: check list of result
		return string(ev.Value), nil
	}

	return "", nil
}

func (c *etcdClient) PublishConfig(dataId string, group string, content string) error {
	log.Printf("PublishConfig for dataId [%s] group [%s] content [%s]", dataId, group, content)
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   c.servers,
		DialTimeout: c.sessionTimeout,
	})

	if err != nil {
		panic(err)
	}
	log.Println("connected to etcd ")

	defer cli.Close()

	kv := clientv3.NewKV(cli)
	itemPath := fmt.Sprintf("/%s/%s", group, dataId)

	putResp, err := kv.Put(context.TODO(),itemPath, content)
	if err != nil {
		log.Printf("Failed to put key")
		return err
	}
	log.Printf("%+v\n",putResp)
	log.Println("Publish config success!")
	return nil
}

func (c *etcdClient) DeleteConfig(dataId string, group string) error {
	log.Printf("DeleteConfig for dataId [%s] group [%s]", dataId, group)
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   c.servers,
		DialTimeout: c.sessionTimeout,
	})

	if err != nil {
		panic(err)
	}
	log.Println("connected to etcd ")

	defer cli.Close()

	// delete the key here
	kv := clientv3.NewKV(cli)
	itemPath := fmt.Sprintf("/%s/%s", group, dataId)
	deleteResp, err := kv.Delete(context.TODO(), itemPath)
	if err != nil {
		panic(err)
	}
	log.Printf("%+v\n",deleteResp)
	log.Printf("key deleted %d\n",deleteResp.Deleted)

	return nil
}