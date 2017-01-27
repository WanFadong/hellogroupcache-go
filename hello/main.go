package main

import (
	"os"
	"fmt"
	"github.com/golang/groupcache"
	"io/ioutil"
	"net/http"
	"log"
	"time"
)

var (
	peers_addrs = []string{"http://127.0.0.1:8001", "http://127.0.0.1:8002"}
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("\r\n Usage local_addr \t\n local_addr must in(8001, 8002)\r\n")
		os.Exit(1)
	}
	local_addr := os.Args[1]
	peers := groupcache.NewHTTPPool("http://127.0.0.1:" + local_addr)// 初始化节点
	peers.Set(peers_addrs...)

	// 初始化Group集群（逻辑，或者说是命名空间）
	var image_cache = groupcache.NewGroup("info", 8 << 30, groupcache.GetterFunc(
		func(ctx groupcache.Context, key string, dest groupcache.Sink) error {
			result, err := ioutil.ReadFile("data/" + key)// []byte
			if err != nil {
				fmt.Printf("read file error %s.\n", err.Error())
				return nil
			}
			fmt.Printf("asking for %s from local file system\n", key)
			dest.SetBytes(result)
			return nil
		}))

	go func() {
		ticker := time.NewTicker(time.Minute)
		for range ticker.C {
			fmt.Println(image_cache.Stats.Gets)
		}
	}()
	log.Fatal(http.ListenAndServe("127.0.0.1:" + local_addr, nil))
}
