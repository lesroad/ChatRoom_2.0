package main

import "strconv"

var waitmsgfromown = make(chan int, 0)

type Group struct {
	//群昵称
	Name string
	//群主
	Owner *Client
	//群成员
	Members []*Client
}

/*
群昵称：xxx
群 主 ：xxx
群人数：xxx
*/

func setUpGroup(name string, client *Client) {
	group := new(Group)
	group.Name = name
	group.Owner = client
	group.Members = append(group.Members, client)
	allGroupsMap[name] = group
}

func waitMsgFromOwn(name string, client *Client) {
	result := <-waitmsgfromown
	if result == 1 {
		allGroupsMap[name].Members = append(allGroupsMap[name].Members, client)

		//反馈结果
		SendMsg2Client("系统消息：成功加入"+name+"群啦哈哈", client)

		//加入日志
		writeMsgToLog("加入群："+name, client)
	} else {
		//反馈结果
		SendMsg2Client("系统消息：唉 群主不让加群", client)
	}
}

//查询所有群信息
func queryAllInfo(client *Client) {
	info := ""
	for _, g := range allGroupsMap {
		info += "群昵称：" + g.Name + "\n"
		info += "群	主：" + g.Owner.name + "\n"
		info += "群人数：" + strconv.Itoa(len(g.Members)) + "人\n\n"
	}

	SendMsg2Client(info, client)
}

//查询单个群信息
func queryOneInfo(name string, client *Client) {
	g := allGroupsMap[name]
	info := "群昵称：" + g.Name + "\n"
	info += "群	主：" + g.Owner.name + "\n"
	info += "群人数：" + strconv.Itoa(len(g.Members)) + "人\n"

	SendMsg2Client(info, client)
}

//踢人
func deleteOne(groupname, clientname string) {
	g := allGroupsMap[groupname]
	var client *Client
	for index, mem := range g.Members {
		if mem.name == clientname {
			g.Members = append(g.Members[:index], g.Members[index+1:]...)
			client = mem
			break
		}
	}

	//告诉被T的人
	SendMsg2Client("系统消息：你被T出"+groupname+"群啦", client)

	//写入日志
	writeMsgToLog("被T出群："+groupname, client)
}

//解散群
func overGroup(name string, client *Client) {
	delete(allGroupsMap, name)

	SendMsg2All(name + "已被群主解散")

	//记录日志
	writeMsgToLog("解散群："+name, client)
}
