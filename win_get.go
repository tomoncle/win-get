// windows download tools
package main

import (
	"gopkg.in/alecthomas/kingpin.v2"
	"log"
)

func main() {
	var (
		url   = kingpin.Flag("url", " request url.").String()
		proxy = kingpin.Flag("proxy", "config socks5 proxy server addr.").String()
	)

	kingpin.Version("1.0.0")
	kingpin.HelpFlag.Short('h')
	kingpin.Parse() // 解析参数
	log.Println("request url: ", *url)
	log.Println("proxy: ", *proxy)

}
