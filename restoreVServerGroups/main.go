package main

import (
	"aliyun-clb/apis"
	"aliyun-clb/client"
	"flag"
	"fmt"
	slb20140515 "github.com/alibabacloud-go/slb-20140515/v4/client"
	"os"
	"path"
	"strings"
)

func request(file *string, client *slb20140515.Client) {
	rs, err := apis.CreateVServerGroup(*file, client)
	if err != nil {
		fmt.Printf("处理文件错误：%s，原因：%s\n", *file, err.Error())
	} else {
		fmt.Printf("处理成功：%s, 结果：%s\n", *file, rs.String())
	}
}

func main() {
	file := flag.String("file", "", "备份文件")
	dir := flag.String("dir", "", "备份文件目录")
	flag.Parse()
	c, err := client.CreateClient()
	if err != nil {
		panic(err)
	}
	if strings.TrimSpace(*file) != "" {
		request(file, c)
	}
	if strings.TrimSpace(*dir) != "" {
		entry, err := os.ReadDir(*dir)
		if err != nil {
			panic(err)
		}
		for _, f := range entry {
			if !f.IsDir() {
				dirFile := path.Join(*dir, f.Name())
				request(&dirFile, c)
			}
		}
	}
}
