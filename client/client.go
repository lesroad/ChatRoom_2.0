package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"path"
	"strings"
)

var (
	ch   = make(chan bool, 0)
	conn net.Conn
)

//错误处理
func cHandleError(err error, msg string) {
	if err != nil {
		fmt.Println(msg, err)
		os.Exit(1)
	}
}

func main() {
	conn, _ = net.Dial("tcp", "127.0.0.1:9009")
	defer conn.Close()

	if len(os.Args) != 2 {
		fmt.Println("not 2 args!")
		os.Exit(1)
	}

	name := os.Args[1]
	//传送名字
	_, err := conn.Write([]byte(name))
	cHandleError(err, "send name error:")

	//独立协程处理发送消息
	go handleWrite() //这里一定要传入conn，否则

	//接受消息
	go handleRead()

	//优雅关闭
	<-ch
}

//独立协程处理发送消息
func handleWrite() {
	reader := bufio.NewReader(os.Stdin)
	for {
		//读取标准输入
		buf, _, _ := reader.ReadLine()

		if strings.Index(string(buf), "upload") == 0 { //上传文件
			handleUploadFile(buf)
		} else if strings.Index(string(buf), "download") == 0 { //下载文件
			handleDownFile(buf)
		} else if strings.Index(string(buf), "exit") == 0 {
			os.Exit(0)
		} else {
			conn.Write(buf)
		}
	}

}

//独立协程接受消息
func handleRead() {
	buf := make([]byte, 1024)

	for {
		n, _ := conn.Read(buf[0:])
		if n > 0 {
			fmt.Println(string(buf[:n]))
		}
	}
}

//处理文件发送
//upload#文件名#内容
func handleUploadFile(buff []byte) {

	strs := strings.Split(string(buff), "#")
	if len(strs) != 2 {
		fmt.Println("上传格式有误，重新上传！")
		return
	}

	buf := make([]byte, 0) //如果非0那么append在后边加了！
	buf = append(buf, []byte("upload#")...)
	filename := path.Base(strs[1]) //取出文件名
	filename = filename + "#"
	buf = append(buf, filename...)
	//文件内容
	filebytes, err := ioutil.ReadFile(strs[1])
	cHandleError(err, "readfile error")

	buf = append(buf, filebytes...)

	conn.Write(buf)
}

//处理下载文件
//download#文件名#地址
func handleDownFile(buf []byte) {
	conn.Write(buf)
}
