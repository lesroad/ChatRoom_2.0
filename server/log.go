package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"time"
)

//写入日志
func writeMsgToLog(str string, client *Client) {
	//检验文件夹是否存在
	_, err := os.Stat(path + "log")

	//不存在创建文件夹
	if os.IsNotExist(err) {
		os.Mkdir(path+"log", 0777)
	}

	//以名字创建文件
	f, err := os.OpenFile(path+"log/"+client.name+".txt", os.O_RDWR|os.O_CREATE, 0777)
	sHandleError(err, "open file error:")
	logMsg := fmt.Sprintln(time.Now().Format("2006-01-02 15:04:05"), client.name+str)
	f.Write([]byte(logMsg))
}

//查看所有日志
func lookAllLog(client *Client) {
	logMsg := ""

	//遍历文件夹所有文件
	dir, err := ioutil.ReadDir(path + "log")
	sHandleError(err, "readdir error:")
	for _, f := range dir {
		buf, _ := ioutil.ReadFile(path + "log/" + f.Name())

		logMsg += string(buf)
	}

	SendMsg2Client(logMsg, client)
}

//查看某人日志
func lookOneLog(name string, client *Client) {
	buf, err := ioutil.ReadFile(path + "log/" + name + ".txt")
	sHandleError(err, "打开日志失败")
	SendMsg2Client(string(buf), client)
}
