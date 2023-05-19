package main

import (
	"encoding/json"
	"fmt"
	"haproxy-go/haproxy"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

type Config struct {
	InPort  int    `json:"inport"`
	OutPort int    `json:"outport"`
	Server  string `json:"server"`
}

func main() {
	var path string
	for idx, args := range os.Args {
		fmt.Println("parameters: "+strconv.Itoa(idx)+":", args)
		if strings.Contains("./configs.json", "./configs.json") {
			path = args
		} else {
			log.Println("not found configs.json")
			return
		}
	}

	//ReadFile函数会读取文件的全部内容，并将结果以[]byte类型返回
	data, err := os.ReadFile(path)
	if err != nil {
		return
	}
	//log.Println((string(data)))

	//读取的数据为json格式，需要进行解码
	var configs []Config
	err = json.Unmarshal(data, &configs)
	if err != nil {
		log.Println(err)
		return
	}
	for _, config := range configs {
		log.Println(config.Server)
		go haproxy.UdpForward(config.InPort, config.Server, config.OutPort)
		go haproxy.TcpForward(config.Server, ":"+strconv.Itoa(config.InPort), ":"+strconv.Itoa(config.OutPort))
	}

	for {
		time.Sleep(10 * time.Second)
	}
}
