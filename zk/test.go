package main

import (
    "fmt"

    "time"

    "github.com/samuel/go-zookeeper/zk"
)

func ZkStateString(s *zk.Stat) string {
    return fmt.Sprintf("Czxid:%d, Mzxid: %d, Ctime: %d, Mtime: %d, Version: %d, Cversion: %d, Aversion: %d, EphemeralOwner: %d, DataLength: %d, NumChildren: %d, Pzxid: %d",
        s.Czxid, s.Mzxid, s.Ctime, s.Mtime, s.Version, s.Cversion, s.Aversion, s.EphemeralOwner, s.DataLength, s.NumChildren, s.Pzxid)
}

func ZkStateStringFormat(s *zk.Stat) string {
    return fmt.Sprintf("Czxid:%d\nMzxid: %d\nCtime: %d\nMtime: %d\nVersion: %d\nCversion: %d\nAversion: %d\nEphemeralOwner: %d\nDataLength: %d\nNumChildren: %d\nPzxid: %d\n",
        s.Czxid, s.Mzxid, s.Ctime, s.Mtime, s.Version, s.Cversion, s.Aversion, s.EphemeralOwner, s.DataLength, s.NumChildren, s.Pzxid)
}

func ZKOperateTest() {
    fmt.Printf("ZKOperateTest\n")

    var hosts = []string{"localhost:2181"}
    conn, _, err := zk.Connect(hosts, time.Second*5)
    if err != nil {
        fmt.Println(err)
        return
    }
    defer conn.Close()

    var path = "/zk_test_go"
    var data = []byte("hello")
    var flags int32 = 0
    // permission
    var acls = zk.WorldACL(zk.PermAll)

    // create
    p, err_create := conn.Create(path, data, flags, acls)
    if err_create != nil {
        fmt.Println(err_create)
        return
    }
    fmt.Println("created:", p)

    // get
    v, s, err := conn.Get(path)
    if err != nil {
        fmt.Println(err)
        return
    }

    fmt.Printf("value of path[%s]=[%s].\n", path, v)
    fmt.Printf("state:\n")
    fmt.Printf("%s\n", ZkStateStringFormat(s))

    // exist
    exist, s, err := conn.Exists(path)
    if err != nil {
        fmt.Println(err)
        return
    }
    fmt.Printf("path[%s] exist[%t]\n", path, exist)
    fmt.Printf("state:\n")
    fmt.Printf("%s\n", ZkStateStringFormat(s))

    // update
    var new_data = []byte("zk_test_new_value")
    s, err = conn.Set(path, new_data, s.Version)
    if err != nil {
        fmt.Println(err)
        return
    }
    fmt.Printf("update state:\n")
    fmt.Printf("%s\n", ZkStateStringFormat(s))

    // get
    v, s, err = conn.Get(path)
    if err != nil {
        fmt.Println(err)
        return
    }

    fmt.Printf("new value of path[%s]=[%s].\n", path, v)
    fmt.Printf("state:\n")
    fmt.Printf("%s\n", ZkStateStringFormat(s))

    // delete
    err = conn.Delete(path, s.Version)
    if err != nil {
        fmt.Println(err)
        return
    }

    // check exist
    exist, s, err = conn.Exists(path)
    if err != nil {
        fmt.Println(err)
        return
    }
    fmt.Printf("after delete, path[%s] exist[%t]\n", path, exist)
    fmt.Printf("state:\n")
    fmt.Printf("%s\n", ZkStateStringFormat(s))
}

func callback(event zk.Event) {
    fmt.Println(">>>>>>>>>>>>>>>>>>>")
    fmt.Println("path:", event.Path)
    fmt.Println("type:", event.Type.String())
    fmt.Println("state:", event.State.String())
    fmt.Println("<<<<<<<<<<<<<<<<<<<")
}

func ZKOperateWatchTest() {
    fmt.Printf("ZKOperateWatchTest\n")

    option := zk.WithEventCallback(callback)
    var hosts = []string{"localhost:2181"}
    conn, _, err := zk.Connect(hosts, time.Second*5, option)
    if err != nil {
        fmt.Println(err)
        return
    }
    defer conn.Close()

    var path1 = "/zk_test_go1"
    var data1 = []byte("zk_test_go1_data1")
    exist, s, _, err := conn.ExistsW(path1)
    if err != nil {
        fmt.Println(err)
        return
    }

    fmt.Printf("path[%s] exist[%t]\n", path1, exist)
    fmt.Printf("state:\n")
    fmt.Printf("%s\n", ZkStateStringFormat(s))

    // try create
    var acls = zk.WorldACL(zk.PermAll)
    p, err_create := conn.Create(path1, data1, zk.FlagEphemeral, acls)
    if err_create != nil {
        fmt.Println(err_create)
        return
    }
    fmt.Printf("created path[%s]\n", p)

    time.Sleep(time.Second * 2)

    exist, s, _, err = conn.ExistsW(path1)
    if err != nil {
        fmt.Println(err)
        return
    }

    fmt.Printf("path[%s] exist[%t] after create\n", path1, exist)
    fmt.Printf("state:\n")
    fmt.Printf("%s\n", ZkStateStringFormat(s))

    // delete
    conn.Delete(path1, s.Version)

    exist, s, _, err = conn.ExistsW(path1)
    if err != nil {
        fmt.Println(err)
        return
    }
    fmt.Printf("path[%s] exist[%t] after delete\n", path1, exist)
    fmt.Printf("state:\n")
    fmt.Printf("%s\n", ZkStateStringFormat(s))
}

func ZKWatchPathTest() {
    fmt.Printf("ZKWatchPathTest")

    var hosts = []string{"localhost:2181"}
    conn, _, err := zk.Connect(hosts, time.Second*5)
    if err != nil {
        fmt.Println(err)
        return
    }
    defer conn.Close()

    // try create root path
    var root_path = "/test_root"

    // check root path exist
    exist, _, err := conn.Exists(root_path)
    if err != nil {
        fmt.Println(err)
        return
    }

    if !exist {
        fmt.Printf("try create root path: %s\n", root_path)
        var acls = zk.WorldACL(zk.PermAll)
        p, err := conn.Create(root_path, []byte("root_value"), 0, acls)
        if err != nil {
            fmt.Println(err)
            return
        }
        fmt.Printf("root_path: %s create\n", p)
    }

    // try create child node
    cur_time := time.Now().Unix()
    ch_path := fmt.Sprintf("%s/ch_%d", root_path, cur_time)
    var acls = zk.WorldACL(zk.PermAll)
    p, err := conn.Create(ch_path, []byte("ch_value"), zk.FlagEphemeral, acls)
    if err != nil {
        fmt.Println(err)
        return
    }
	fmt.Printf("ch_path: %s create\n", p)

    for  {
        _, _, ch, _ := conn.GetW(ch_path)
		select {
		case e := <-ch:
			if e.Err == nil {
				if e.Type == zk.EventNodeDataChanged {
                    fmt.Printf("has node[%s] data changed\n", e.Path)
                    fmt.Printf("%+v\n", e)
                    v, _, err := conn.Get(ch_path)
                    if err != nil {
                        fmt.Println(err)
                        return
                    }
        
                    fmt.Printf("value of path[%s]=[%s].\n", ch_path, v)
                }
			}
        }
	}
}

func main() {
	ZKOperateTest()
	ZKOperateWatchTest()
	ZKWatchPathTest()
}