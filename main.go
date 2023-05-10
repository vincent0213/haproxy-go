package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
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
		if strings.Contains("./configs.json", "./configs.json") {
			path = args
		}
		fmt.Println("参数"+strconv.Itoa(idx)+":", args)
	}

	//ReadFile函数会读取文件的全部内容，并将结果以[]byte类型返回
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return
	}
	//log.Println((string(data)))

	//读取的数据为json格式，需要进行解码
	var configs []Config
	err = json.Unmarshal([]byte(data), &configs)
	if err != nil {
		log.Println(err)
		return
	}
	for _, config := range configs {
		log.Println(config.Server)
		go udpLocal(config.InPort, config.Server, config.OutPort)
		go tcpForward(config.Server, ":"+strconv.Itoa(config.InPort), ":"+strconv.Itoa(config.OutPort))
	}

	for {
		time.Sleep(10 * time.Second)
	}

	//flag.StringVar(&flags.Server1, "s", "18.179.166.28", "client connect address or url")
	//flag.Parse()
	//
	//go udpLocal(13531, flags.Server1, 13531)
	//go udpLocal(13532, flags.Server2, 13532)
	//go udpLocal(13533, flags.Server3, 13533)
	//go tcpForward(flags.Server, ":13531")
	//go tcpForward(flags.Server, ":13532")
	//tcpForward(flags.Server, ":13533")

}
