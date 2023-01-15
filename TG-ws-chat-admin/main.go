package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"sort"
)

var (
	// 升级成 WebSocket 协议
	upgrader = websocket.Upgrader{
		// 允许CORS跨域请求
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	conn *websocket.Conn
	err  error

	// 应用一运行，就初始化 CenterHandler 处理中心对象
	handler = CenterHandler{
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
	}
)

// CenterHandler 处理中心，关联着每个 Client 的注册、注销、广播通道，相当于每个用户的中心通讯的中介。
type CenterHandler struct {
	// 广播通道，有数据则循环每个用户广播出去
	broadcast chan []byte
	// 注册通道，有用户进来 则推到用户集合map中
	register chan *Client
	// 注销通道，有用户关闭连接 则将该用户剔出集合map中
	unregister chan *Client
	// 用户集合，每个用户本身也在跑两个协程，监听用户的读、写的状态
	clients map[*Client]bool
}

// 处理中心的一个接口，监控状态
func (ch *CenterHandler) monitoring() {
	for {
		select {
		// 注册，新用户连接过来会推进注册通道，这里接收推进来的用户指针
		case client := <-ch.register:
			ch.clients[client] = true
			// 注销，关闭连接或连接异常会将用户推出群聊
		case client := <-ch.unregister:
			delete(ch.clients, client)
			// 消息，监听到有新消息到来
		case message := <-ch.broadcast:
			// 推送给每个用户的通道，每个用户都有跑协程起了writePump的监听

			for client := range ch.clients {
				println("收到消息，message：" + string(message))

				client.send <- message
			}
		}
	}
}

// Client 抽象出来的 Client，里面有这个 websocket 连接的 读 和 写 操作
type Client struct {
	handler *CenterHandler
	conn    *websocket.Conn
	// 每个用户自己的循环跑起来的状态监控
	send chan []byte
	UID  string
}

// Message 消息体
type Message struct {
	MessageType string `json:"type"`
	SenderUID   string `json:"senderUID"`
	SendTime    string `json:"sendTime"`
	Content     string `json:"content"`
}

// HostCount 用户存活体
type HostCount struct {
	MessageType string   `json:"type"`
	UIDList     []string `json:"UIDList"`
	Count       int      `json:"count"`
}

// 读，监听客户端是否有推送内容过来服务端
func (c *Client) readPump() {
	defer func() {
		c.handler.unregister <- c
		c.conn.Close()
	}()
	for {
		// 循环监听是否该用户是否要发言
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			// 异常关闭的处理
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}

		//消息处理
		//绑定到结构体
		var messageObj = &Message{}
		err = json.Unmarshal(message, &messageObj)
		if err != nil {
			fmt.Println("消息处理出错：", err)
		}

		//检查type类型
		//如果为hostCount则检查所有用户的UID和数量并发送回客户端
		if messageObj.MessageType == "hostCount" {
			UIDList, count := getHostCount()

			//转成JSON
			hostCount := &HostCount{
				MessageType: "hostCount",
				UIDList:     UIDList,
				Count:       count,
			}
			hostCountJSON, err := json.Marshal(hostCount)
			if err != nil {
				fmt.Println("hostCountJSON转换出错：", err)
			}

			//发送回客户端
			go func() {
				if err := c.conn.WriteMessage(websocket.TextMessage, hostCountJSON); err != nil {
					return
				}
			}()
		}

		//推给广播中心，广播中心再推给每个用户
		c.handler.broadcast <- message
	}
}

// 写，主动推送消息给客户端
func (c *Client) writePump() {
	defer func() {
		c.handler.unregister <- c
		c.conn.Close()
	}()
	for {
		// 广播推过来的新消息，马上通过websocket推给自己
		message, _ := <-c.send
		if err := c.conn.WriteMessage(websocket.TextMessage, message); err != nil {
			return
		}
	}
}

// 获取所有用户的UID以及数量
func getHostCount() (UIDSlice []string, hostCount int) {
	//获取数量
	hostCount = len(handler.clients)

	//定义切片
	for i, _ := range handler.clients {
		UIDSlice = append(UIDSlice, i.UID)
	}
	//排序
	sort.Strings(UIDSlice)
	fmt.Println(UIDSlice, hostCount)
	return
}

func main() {
	// 起个协程跑起来，监听注册、注销、消息 3 个 channel
	go handler.monitoring()

	// websocket 请求，建立双工通讯连接
	http.HandleFunc("/ws", func(writer http.ResponseWriter, request *http.Request) {
		// 由 http 升级成为 websocket 服务
		if conn, err = upgrader.Upgrade(writer, request, nil); err != nil {
			log.Println(err)
			return
		}

		//获取用户的UID
		_, message, err := conn.ReadMessage()
		if err != nil {
			// 异常关闭的处理
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
		}

		//获取发来的的UID
		UID := string(message)
		fmt.Printf("UID %s 加入了服务器\n", UID)

		// 为每个连接创建一个 Client 实例，（实际上这里应该还有绑定用户真实信息的操作）
		client := &Client{&handler, conn, make(chan []byte, 256), UID}
		// 推给监控中心注册到用户集合中
		handler.register <- client
		// 每个 client 都挂起 2 个新的协程，监控读、写状态
		go client.writePump()
		go client.readPump()
	})

	if err := http.ListenAndServe(":8888", nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
