# !bin/bash 
set -o errexit

# 测试模式
# 1. 正常模式      normal
# 2. 网路链接模式   connect
# 3. ping网关模式  ping 
# 4. 主数据服压测   main-data
export TEST_MODE=main-data

# 间隔多少MS 添加一个客户端
export ADD_CLIENT_CD_MS=200


# export AGENT_URL=agent601-dev.game.melandworld.com
export AGENT_URL=192.168.50.171:5700
# export AGENT_URL=192.168.50.15:5700


# 默认连接的scene service appid 
# 为空 则使用manager service 下发的 scene service
# e.g: game-service-world-735
export DEFAULT_SCENE_SER="game-service-world-735" 

# 压测的客户端数量上限
export CLIENT_NUM=1

# 压测的客户端模拟 Id 起始值 
export CLIENT_IDX_BEGIN=10000

go run src/main.go




