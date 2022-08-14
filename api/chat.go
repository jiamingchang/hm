package api

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"time"
)

type User struct {
	new chan mess
	conn *websocket.Conn
}
type mess struct {
	PostId  int    `json:"postId"`
	Message string `json:"message"`
}

type Hub struct {
	//用户列表，保存所有用户
	userList map[*User]bool
	//注册chan，用户注册时添加到chan中
	register chan *User
	//注销chan，用户退出时添加到chan中，再从map中删除
	unregister chan *User
	//广播消息，将消息广播给所有连接
	broadcast chan mess
	//谁发的消息，消息不广播回去
	whosend chan *User
}
var hub = &Hub{
	userList:   make(map[*User]bool),
	register:   make(chan *User),
	unregister: make(chan *User),
	broadcast:  make(chan mess),
	whosend:    make(chan *User),
}

var pingTicker  = time.NewTicker(time.Second * 5)

var upgrader = websocket.Upgrader{
	ReadBufferSize:   1024,
	WriteBufferSize:  1024,
	HandshakeTimeout: 5 * time.Second,
	// 取消ws跨域校验
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}
func readdata(user *User){
	for{
		var me mess
		conn:= user.conn
		//_, data, err := conn.ReadMessage()
		err := conn.ReadJSON(&me)
		if err !=nil {
			log.Println(err)
			return
		}
		//log.Println(me)
		hub.broadcast <- me
		hub.whosend <- user
	}
}
func writedata(user *User){
	for {
		select {
		case n, _ := <-user.new:
			conn := user.conn
			err := conn.WriteJSON(n)
			if err != nil {
				log.Println(err)
				return
			}
		}
	}
}
func live(){
	for{
		select {
		case <-pingTicker.C:
			// 服务端心跳:每5秒ping一次客户端，查看其是否在线
			for user:= range hub.userList {
				conn := user.conn
				conn.SetWriteDeadline(time.Now().Add(time.Second * 6))
				err := conn.WriteMessage(websocket.PingMessage, []byte{})
				if err != nil {
					log.Println("send ping err:", err)
					hub.unregister <- user
					conn.Close()
				}
			}
		}
	}
}

func wshandle(w gin.ResponseWriter,r *http.Request){
	conn, _ := upgrader.Upgrade(w, r, nil)
	user := &User{
		new: make(chan mess),
		conn: conn,
	}
	hub.register <- user
	go readdata(user)
	go writedata(user)
}

func (h *Hub) run() {
	for {
		select {
		//从注册chan中取数据
		case user := <-h.register:
			//取到数据后将数据添加到用户列表中
			h.userList[user] = true
			log.Println(user.conn.RemoteAddr().String())
		case user := <-h.unregister:
			//从注销列表中取数据，判断用户列表中是否存在这个用户，存在就删掉
			if _, ok := h.userList[user]; ok {
				delete(h.userList, user)
			}
		case data := <-h.broadcast:
			//从广播chan中取消息，然后遍历给每个用户，发送到用户的msg中
			who := <- h.whosend
			for u := range h.userList {
				if u.conn.RemoteAddr().String() != who.conn.RemoteAddr().String() {
					u.new <- data
				}
			}
		}
	}
}
func Chat(context *gin.Context){
	// 将普通的http GET请求升级为websocket请求
	wshandle(context.Writer, context.Request)
}
func init(){
	go hub.run()
	go live()
}

