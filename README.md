# TRPC-Go与Nano结合的房间游戏示例

这是一个使用TRPC-Go和Nano游戏服务器框架结合实现的简单房间游戏示例。

## 功能特点

- 基于TRPC-Go的服务框架
- 集成Nano游戏服务器框架
- 支持房间创建和加入
- 支持房间内消息广播
- 提供WebSocket客户端示例

## 项目结构

```
.
├── examples/
│   └── client.html          # WebSocket客户端示例
├── src/
│   ├── cmd/
│   │   └── main.go          # 主程序入口
│   └── internal/
│       └── game/            # 游戏相关代码
│           ├── room.go              # 房间实体
│           ├── room_manager.go      # 房间管理器
│           └── game_component.go    # 游戏组件
└── README.md                # 项目说明
```

## 快速开始

### 1. 运行服务器

```bash
go run src/cmd/main.go
```

服务器将在以下端口启动：
- TRPC服务: 根据trpc_go.yaml配置（默认8000端口）
- Nano游戏服务器: 3250端口

### 2. 使用客户端

打开`examples/client.html`文件在浏览器中访问，或者使用以下命令启动一个简单的HTTP服务器：

```bash
cd examples
python -m http.server 8080
```

然后在浏览器中访问 http://localhost:8080/client.html

## 客户端使用说明

1. 创建房间：输入房间名称，点击"创建房间"按钮
2. 加入房间：在房间列表中点击要加入的房间
3. 发送消息：在聊天输入框中输入消息，点击"发送"按钮或按回车键发送
4. 刷新房间列表：点击"刷新房间列表"按钮

## API说明

### 游戏组件API

- `game.createRoom`: 创建新房间
  - 参数: `{ name: string }`
  - 返回: `{ id: number, name: string, players: number }`

- `game.joinRoom`: 加入房间
  - 参数: `{ roomId: number }`
  - 返回: `{ id: number, name: string, players: number }`

- `game.sendMessage`: 发送消息
  - 参数: `{ content: string }`
  - 无返回值，消息通过广播发送

- `game.listRooms`: 获取房间列表
  - 无参数
  - 返回: `[{ id: number, name: string, players: number }, ...]`

### 事件通知

- `onPlayerJoin`: 当新玩家加入房间时触发
  - 数据: `{ playerId: number }`

- `onPlayerLeave`: 当玩家离开房间时触发
  - 数据: `{ playerId: number }`

- `onMessage`: 当收到新消息时触发
  - 数据: `{ playerId: number, message: { content: string } }`

## 扩展建议

这个示例实现了基本的房间和消息功能，可以进一步扩展：

1. 添加用户认证和登录功能
2. 实现更复杂的游戏逻辑
3. 添加房间设置和游戏配置
4. 实现游戏状态同步
5. 添加数据持久化