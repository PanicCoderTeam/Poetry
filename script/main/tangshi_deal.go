// client.go
package main

// import (
// 	"log"
// 	"time"

// 	"github.com/gorilla/websocket"
// 	"github.com/lonng/nano/serialize"
// 	"github.com/lonng/nano/serialize/json" // 使用JSON序列化
// )

// const (
// 	serverAddr = "ws://localhost:3250/nano"
// )

// // 协议结构体定义
// type (
// 	LoginRequest struct {
// 		Username string `json:"username"`
// 		Password string `json:"password"`
// 	}

// 	LoginResponse struct {
// 		Code  int    `json:"code"`
// 		Token string `json:"token"`
// 	}
// )

// func main() {
// 	// 创建WebSocket连接
// 	conn, _, err := websocket.DefaultDialer.Dial(serverAddr, nil)
// 	if err != nil {
// 		log.Fatal("连接失败:", err)
// 	}
// 	defer conn.Close()

// 	// 初始化序列化器
// 	serializer := json.NewSerializer()

// 	// 启动接收协程
// 	go receiveMessages(conn, serializer)

// 	// 发送登录请求
// 	sendRequest(conn, serializer, "game.login", &LoginRequest{
// 		Username: "player1",
// 		Password: "123456",
// 	})

// 	// 保持连接（示例心跳）
// 	ticker := time.NewTicker(10 * time.Second)
// 	defer ticker.Stop()

// 	for range ticker.C {
// 		if err := conn.WriteMessage(websocket.PingMessage, nil); err != nil {
// 			log.Println("心跳失败:", err)
// 			return
// 		}
// 	}
// }

// // 发送请求方法
// func sendRequest(conn *websocket.Conn, serializer serialize.Serializer, route string, data interface{}) {
// 	// 构造Nano协议包
// 	msg := map[string]interface{}{
// 		"type":  0, // 0=Request
// 		"id":    time.Now().UnixNano(),
// 		"route": route,
// 		"data":  data,
// 	}

// 	// 序列化消息
// 	payload, err := serializer.Marshal(msg)
// 	if err != nil {
// 		log.Println("序列化失败:", err)
// 		return
// 	}

// 	// 发送消息
// 	if err := conn.WriteMessage(websocket.TextMessage, payload); err != nil {
// 		log.Println("发送失败:", err)
// 	}
// }

// // 接收消息方法
// func receiveMessages(conn *websocket.Conn, serializer serialize.Serializer) {
// 	for {
// 		_, message, err := conn.ReadMessage()
// 		if err != nil {
// 			log.Println("接收错误:", err)
// 			return
// 		}

// 		// 反序列化消息
// 		var res map[string]interface{}
// 		if err := serializer.Unmarshal(message, &res); err != nil {
// 			log.Println("反序列化失败:", err)
// 			continue
// 		}

// 		// 处理不同消息类型
// 		switch res["type"].(float64) { // JSON数字默认解析为float64
// 		case 2: // Response
// 			handleResponse(res)
// 		case 3: // Push
// 			handlePush(res)
// 		}
// 	}
// }

// // 处理响应
// func handleResponse(res map[string]interface{}) {
// 	log.Printf("收到响应 ID=%v 数据=%+v", res["id"], res["data"])

// 	// 类型断言处理登录响应
// 	if res["route"] == "auth.Login" {
// 		var loginRes LoginResponse
// 		if data, ok := res["data"].(map[string]interface{}); ok {
// 			loginRes.Code = int(data["code"].(float64))
// 			loginRes.Token = data["token"].(string)
// 			log.Println("登录成功，Token:", loginRes.Token)
// 		}
// 	}
// }

// // 处理服务端推送
// func handlePush(res map[string]interface{}) {
// 	log.Printf("收到推送 Route=%s 数据=%+v", res["route"], res["data"])
// }
