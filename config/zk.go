package config

import (
	"log"
	"time"
	"fmt"
	"errors"
	"os"

	"github.com/samuel/go-zookeeper/zk"
)

type zkClient struct {
	servers []string
	sessionTimeout time.Duration
}

func NewZKStore() KVStore {
	zkHost := os.Getenv("ZK_HOST")
	var hosts []string

	if zkHost != "" {
		hostPath := fmt.Sprintf("%s:%s", zkHost, "2181")
		hosts = []string{hostPath}
	} else {
		hosts = []string{"localhost:2181"}
	}
	
	var timeout = time.Second*5

	return &zkClient{hosts, timeout}
}

func (c *zkClient) GetConfig(dataId string, group string) (string, error) {
	log.Printf("GetConfig for dataId [%s] group [%s]", dataId, group)

	conn, _, err := zk.Connect(c.servers, c.sessionTimeout)
    if err != nil {
        log.Printf("Error : failed to connect to zk -> %s ", err)
        return "", err
    }
	defer conn.Close()

	var itemPath = fmt.Sprintf("/%s/%s", group, dataId)
	v, _, err := conn.Get(itemPath)
    if err != nil {
        log.Printf("Error : failed to get item from zk -> %s ", err)
        return "", err
    }

    log.Printf("value of path[%s]=[%s].\n", itemPath, v)
	
	return string(v), nil
}

func (c *zkClient) PublishConfig(dataId string, group string, content string) error {
	log.Printf("PublishConfig for dataId [%s] group [%s] content [%s]", dataId, group, content)

	conn, _, err := zk.Connect(c.servers, c.sessionTimeout)
    if err != nil {
        log.Printf("Error : failed to connect to zk -> %s", err)
        return err
    }
	defer conn.Close()

	rootPath := fmt.Sprintf("/%s", group)
	exist, _, err := conn.Exists(rootPath)
    if err != nil {
        log.Printf("Error : failed to check existence of item -> %s", err)
        return err
	}

	var flags int32 = 0
	var acls = zk.WorldACL(zk.PermAll)
	
	// Create root node if it does not exist
	if !exist {
		var rootData = []byte("")
		
		p, err_create := conn.Create(rootPath, rootData, flags, acls)
		if err_create != nil {
			log.Printf("Error : failed to create item -> %s", err_create)
			return err_create
		}
		log.Printf("root node created: %s", p)
	}

	var itemPath = fmt.Sprintf("/%s/%s", group, dataId)

	existItem, _, err := conn.Exists(itemPath)
    if err != nil {
        log.Printf("Error : failed to check existence of item -> %s", err)
        return err
	}

	if existItem {
		_, err_set := conn.Set(itemPath, []byte(content), -1)
		if err_set != nil {
			log.Printf("Error : failed to set item -> %s", err_set)
			return err_set
		}
		log.Printf("item node value set complete")

	} else {
		p, err_create := conn.Create(itemPath, []byte(content), flags, acls)
		if err_create != nil {
			log.Printf("Error : failed to create item -> %s", err_create)
			return err_create
		}
		log.Printf("item node created: %s", p)
	}

	return nil
}

func (c *zkClient) DeleteConfig(dataId string, group string) error {
	log.Printf("DeleteConfig for dataId [%s] group [%s]", dataId, group)

	conn, _, err := zk.Connect(c.servers, c.sessionTimeout)
    if err != nil {
        log.Printf("Error : failed to connect to zk -> %s", err)
        return err
    }
	defer conn.Close()

	var itemPath = fmt.Sprintf("/%s/%s", group, dataId)

	// check exist
    exist, s, err := conn.Exists(itemPath)
    if err != nil {
        log.Printf("Error : failed to check existence -> %s", err)
        return err
	}

	if !exist {
		log.Printf("Error : input config %s does not exist", itemPath)
        return errors.New("config item does not exist")
	}
	
	// delete
    err = conn.Delete(itemPath, s.Version)
    if err != nil {
        log.Printf("Error : failed to delete item -> %s", err)
        return err
    }
	return nil
}

func (c *zkClient) ListenConfig(dataId string, group string, timeout int) (string, error) {
	log.Printf("ListenConfig for dataId [%s] group [%s] timeout [%d]", dataId, group, timeout)

	conn, _, err := zk.Connect(c.servers, c.sessionTimeout)
    if err != nil {
        log.Printf("Error : failed to connect to zk -> %s", err)
        return "", err
    }
	defer conn.Close()

	waiter := time.Tick(time.Duration(timeout) * time.Millisecond)
	var itemPath = fmt.Sprintf("/%s/%s", group, dataId)

	_, _, ch, _ := conn.GetW(itemPath)
	select {
	case e := <-ch:
		if e.Err == nil {
			if e.Type == zk.EventNodeDataChanged {
				log.Printf("has node[%s] data changed\n", e.Path)
				log.Printf("%+v\n", e)
				v, _, err := conn.Get(itemPath)
				if err != nil {
					fmt.Println(err)
					return "", nil
				}
	
				log.Printf("value of path[%s]=[%s].\n", itemPath, v)
				return string(v), nil
			}
		}
	case ts := <-waiter:
		log.Printf("pulling timeout reached: %s", ts)
		return "", nil
	}

	// should not reach here
	return "", nil
}