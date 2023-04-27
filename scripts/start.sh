# !bin/bash 
set -o errexit

# 测试模式
# 1. 正常模式      normal
# 2. 网路链接模式   connect
# 3. ping网关模式  ping 
export TEST_MODE=normal

# 间隔多少MS 添加一个客户端
export ADD_CLIENT_CD_MS=200


export AGENT_URL=agent601-dev.game.melandworld.com
# export AGENT_URL=192.168.50.171:5700
# export AGENT_URL=192.168.50.15:5700

# 压测的客户端数量上限
export CLIENT_NUM=500

# 压测的客户端模拟 Id 起始值 
export CLIENT_IDX_BEGIN=10000

go run src/main.go




