# !bin/bash 
set -o errexit

## 测试模式
## 1. 正常模式      normal
## 2. 网路链接模式   connect
## 3. ping网关模式  ping 
export TEST_MODE=normal

# export AGENT_URL=192.168.50.171:5700
# export AGENT_URL=192.168.50.15:5700
export AGENT_URL=agent601-dev.game.melandworld.com

export CLIENT_NUM=1
export CLIENT_IDX_BEGIN=20000

go run src/main.go

