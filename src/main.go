package main

import (
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/Meland-Inc/meland-client/src/client"
	"github.com/Meland-Inc/meland-client/src/common/serviceLog"
	"github.com/spf13/cast"
)

func main() {

	agentUrl := os.Getenv("AGENT_URL")

	clientNum, err := cast.ToIntE(os.Getenv("CLIENT_NUM"))
	if err != nil {
		panic(err)
	}

	clientIdxBegin, err := cast.ToIntE(os.Getenv("CLIENT_IDX_BEGIN"))
	if err != nil {
		panic(err)
	}

	testModel := os.Getenv("TEST_MODE")

	serviceLog.Init(int64(clientIdxBegin), true)

	wg := new(sync.WaitGroup)
	for i := 0; i < clientNum; i++ {
		go func(idx int) {
			wg.Add(1)
			cli := client.NewGameClient(testModel, agentUrl, fmt.Sprint(clientIdxBegin+idx), int64(clientIdxBegin+idx))
			cli.Run()
			wg.Done()
		}(i)
		time.Sleep(time.Millisecond * 200)
	}
	wg.Wait()
}
