package main

import (
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
)

var (
	//客户端信息，昵称为键
	allClientsMap = make(map[string]*Client)

	//所有群,群名为键
	allGroupsMap = make(map[string]*Group)
	path         = "F:/go/gocode/src/聊天室2.0/"
)

//错误处理
func sHandleError(err error, msg string) {
	if err != nil {
		fmt.Println(msg, err)
		os.Exit(1)
	}
}

func main() {
	listerner, err := net.Listen("tcp", "127.0.0.1:9009")
	sHandleError(err, "listen error:")
	defer func() {
		SendMsg2All("系统消息：服务器即将关闭，都洗洗睡吧！")
		listerner.Close()
	}()

	for {
		//循环接入每一个用户
		conn, err := listerner.Accept()
		sHandleError(err, "accept error:")

		//获取名字
		buf := make([]byte, 1024)
		n, _ := conn.Read(buf)
		name := string(buf[:n])

		//提醒上线
		fmt.Printf(name + "上线了\n")
		for _, cli := range allClientsMap {
			cli.conn.Write([]byte(name + "上线了"))
		}

		//丢入map
		one := &Client{conn, name}
		allClientsMap[name] = one

		//单独开协程通信
		go withClient(one)
	}
}

//单独开协程通信
func withClient(client *Client) {
	buf := make([]byte, 1024)
	for {
		str := ""

		//传输文件需要多次读取！
		for {
			n, _ := client.conn.Read(buf)
			if n == 0 { //客户端关闭时n=0！
				SendMsg2All("系统消息：" + client.name + "下线了！")
				fmt.Println(client.name + "下线了")
				return
			}
			str += string(buf[:n])

			if n < 1024 { //防止阻塞
				break
			}
		}

		if len(str) > 0 {

			strs := strings.Split(str, "#")

			first := strs[0]
			switch first {
			case "upload":
				//upload#文件#内容
				name := strs[1]
				content := str[len(strs[0])+len(strs[1]):]

				//写入文件
				uploadFile(name, content, client)

			case "download":
				//download#文件名#地址
				name := strs[1]
				clipath := strs[2]

				downloadFile(name, clipath, client)

			case "all":
				//向所有人传达消息
				SendMsg2All(client.name + ":" + strs[1])

				//写入日志
				writeMsgToLog("发送全体消息："+strs[1], client)

			case "group_setup":
				groupname := strs[1]

				//搜索有无相同的群名
				for _, group := range allGroupsMap {
					if group.Name == groupname {
						SendMsg2Client("系统消息：骚年 有人用过该群名啦", client)
						break
					}
				}

				//处理建群
				setUpGroup(groupname, client)

				//反馈结果
				SendMsg2Client("系统消息：创建群成功！", client)

				//写入日志
				writeMsgToLog("创建了新群："+groupname, client)

			case "group_join":
				groupname := strs[1]

				//搜索有无此群
				if allGroupsMap[groupname] == nil {
					SendMsg2Client("系统消息：没有此群啦", client)
					break
				}

				//向群主发送加群通知
				SendMsg2Client("系统消息："+client.name+"申请加入群"+groupname, allGroupsMap[groupname].Owner)

				//另开协程群主等待同意
				go waitMsgFromOwn(groupname, client)

				//反馈结果
				SendMsg2Client("系统消息：已发送加群申请，请耐心等候！", client)

				//写入日志
				writeMsgToLog("申请加群："+groupname, client)

			case "group_joinreply":
				result, err := strconv.Atoi(strs[1])
				sHandleError(err, "atoi error")

				//向管道送群主意见
				waitmsgfromown <- result

			case "log":
				name := strs[1]

				if name == "all" {
					lookAllLog(client)
				} else {
					lookOneLog(name, client)
				}

			case "group_info":
				name := strs[1]

				if name == "all" {
					queryAllInfo(client)
				} else {
					//查询有无此群
					if allGroupsMap[name] == nil {
						SendMsg2Client("系统消息：没这个群啊", client)
						break
					}
					queryOneInfo(name, client)
				}
			case "group_delete":
				groupname := strs[1]
				clientname := strs[2]

				if allGroupsMap[groupname] == nil || allGroupsMap[groupname].Owner.name != client.name || allClientsMap[clientname] == nil {
					SendMsg2Client("系统消息：踢人失败", client)
					break
				}

				//执行踢人
				deleteOne(groupname, clientname)

				//写入日志
				writeMsgToLog("把"+clientname+"T出"+groupname, client)

			case "group_over":
				name := strs[1]

				//群主身份确认
				if allGroupsMap[name].Owner.name != client.name {
					SendMsg2Client("系统消息：你是群主吗啊你就解散", client)
					break
				}
				//执行灭群
				overGroup(name, client)
			default:
				//可能是输错的情况
				if allClientsMap[strs[0]] == nil {
					SendMsg2Client("系统消息：输入错误！", client)
					break
				}

				//向某个人发送信息
				msg := strs[1]
				SendMsg2One(strs[0], client.name+":"+msg)

				//写入日志
				writeMsgToLog("向"+strs[0]+"发送："+msg, client)
			}

		}

	}
}
