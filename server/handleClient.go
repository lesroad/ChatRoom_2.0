package main

import "net"

type Client struct {
	//客户端连接
	conn net.Conn
	//昵称
	name string
}

func SendMsg2Client(str string, client *Client) {
	client.conn.Write([]byte(str))
}

func SendMsg2All(str string) {
	for _, cli := range allClientsMap {
		cli.conn.Write([]byte(str))
	}
}

func SendMsg2One(name, msg string) {
	for _, c := range allClientsMap {
		if c.name == name {
			SendMsg2Client(msg, c)
		}
	}
}
