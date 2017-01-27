package main

import (
	"os"
	"fmt"
	"github.com/golang/groupcache"
	"io/ioutil"
	"net/http"
	"log"
)


var (
	peers_addrs = []string{"http://127.0.0.1:8001", "http://127.0.0.1:8002", "http://127.0.0.1:8003"}
)

func test() {
	if len(os.Args) != 2 {
		fmt.Println("\r\n Usage local_addr \t\n local_addr must in(8001, 8002, 8003)\r\n")
		os.Exit(1)
	}
	local_addr := os.Args[1]
	peers := groupcache.NewHTTPPool("http://127.0.0.1:" + local_addr)// 初始化节点
	peers.Set(peers_addrs...)

	// 初始化Group集群（逻辑，或者说是命名空间）
	var image_cache = groupcache.NewGroup("imagegroup", 8 << 30, groupcache.GetterFunc(
		func(ctx groupcache.Context, key string, dest groupcache.Sink) error {
			result, err := ioutil.ReadFile(key)// []byte
			if err != nil {
				fmt.Printf("read file error %s.\n", err.Error())
				return nil
			}
			fmt.Printf("asking for %s from local file system\n", key)
			dest.SetBytes(result)
			return nil
		}))
	http.HandleFunc("/image", func(rw http.ResponseWriter, r *http.Request) {
		var data []byte
		k := r.URL.Query().Get("id")
		fmt.Printf("user get %s from groupcache\n", k)
		image_cache.Get(nil, k, groupcache.AllocatingByteSliceSink(&data))
		rw.Write([]byte(data))
	})
	log.Fatal(http.ListenAndServe("127.0.0.1:" + local_addr, nil))
}
