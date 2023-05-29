module github.com/Meland-Inc/meland-client

go 1.17

require (
	github.com/google/uuid v1.3.0 // indirect
	github.com/gorilla/websocket v1.5.0 // indirect
	google.golang.org/protobuf v1.28.1 // indirect
)

require (
	// 使用本地go代码仓库方式: https://zhuanlan.zhihu.com/p/109828249
	game-message-core v0.0.0
	github.com/golang/snappy v0.0.4
	github.com/spf13/cast v1.5.0
)

replace game-message-core => ./src/game-message-core/messageGo
