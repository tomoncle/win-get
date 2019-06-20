// windows download tools
package main

import (
	"fmt"
	"gopkg.in/alecthomas/kingpin.v2"
	"gopkg.in/cheggaaa/pb.v1"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

// getFileName : validate url and get file name
func getFileName(requestUrl string) (string, error) {
	fileUrl, err := url.Parse(requestUrl)
	if err != nil {
		return "", err
	}
	filePath := fileUrl.Path
	segments := strings.Split(filePath, "/")
	fileName := segments[len(segments)-1]
	if strings.Contains(fileName, "?") {
		return strings.Split(filePath, "?")[0], nil
	}
	return fileName, nil
}

// downloader : download file
func downloader(fileName string, requestUrl string, proxy string) (int64, error) {
	file, err := os.Create(fileName)
	if err != nil {
		return 0, err
	}
	defer file.Close()

	request := http.Client{
		CheckRedirect: func(r *http.Request, via []*http.Request) error {
			r.URL.Opaque = r.URL.Path
			return nil
		},
	}
	response, err := request.Get(requestUrl)
	if err != nil {
		return 0, err
	}
	defer response.Body.Close()
	// console request status
	fmt.Println(response.Status)
	// get file size
	fileSize := response.ContentLength
	go func() {
		n, err := io.Copy(file, response.Body)
		if n != fileSize {
			fmt.Println("Truncated")
		}
		if err != nil {
			fmt.Printf("Error: %v", err)
		}
	}()

	countSize := int(fileSize)
	bar := pb.StartNew(countSize)
	var fi os.FileInfo
	for fi == nil || fi.Size() < fileSize {
		fi, _ = file.Stat()
		bar.Set(int(fi.Size()))
		time.Sleep(time.Millisecond)
	}
	bar.FinishPrint("The End!")
	return fileSize, nil
}

func main() {
	var (
		requestUrl   = kingpin.Flag("url", " request url.").String()
		requestProxy = kingpin.Flag("proxy", "config socks5 proxy server addr.").String()
	)
	kingpin.Version("1.0.0")
	kingpin.HelpFlag.Short('h')
	kingpin.Parse() // 解析参数
	if *requestUrl == "" {
		fmt.Println("Usage: win_get -h")
		os.Exit(0)
	}
	fmt.Println("request url: ", *requestUrl)
	if *requestProxy != "" {
		fmt.Println("proxy: ", *requestProxy)
	}
	log.Println("Downloading file.")

	fileName, err := getFileName(*requestUrl)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	fileSize, err := downloader(fileName, *requestUrl, *requestProxy)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	fmt.Printf("%s with %v bytes downloaded", fileName, fileSize)

}
