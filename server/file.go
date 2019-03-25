package main

import (
	"fmt"
	"io/ioutil"
	"os"
)

func uploadFile(name, content string, client *Client) {
	fmt.Println(client.name + "上传文件：" + name)

	//检验文件夹是否存在
	_, err := os.Stat(path + "file")

	//不存在创建文件夹
	if os.IsNotExist(err) {
		os.Mkdir(path+"file", 0777)
	}

	//存储到指定地址
	fpath := path + "file/" + name
	err = ioutil.WriteFile(fpath, []byte(content), 0777)
	sHandleError(err, "write file error:")

	//反馈结果
	SendMsg2Client("系统消息：文件上传成功！", client)

	//写入日志
	writeMsgToLog("上传文件"+name, client)
}

func downloadFile(name, clipath string, client *Client) {
	fmt.Println(client.name + "下载文件：" + name)

	_, err := os.OpenFile(path+"file/"+name, os.O_RDONLY, 0777)
	if err != nil {
		SendMsg2Client("系统消息：没有此文件！", client)
		return
	}

	//查找到此文件
	files, _ := ioutil.ReadFile(path + "file/" + name)

	//写入文件
	err = ioutil.WriteFile(clipath+"/"+name, files, 0777)
	sHandleError(err, "writefile error:")

	//反馈结果
	SendMsg2Client("系统消息：文件下载成功！", client)

	//写入日志
	writeMsgToLog("下载文件"+name, client)
}
