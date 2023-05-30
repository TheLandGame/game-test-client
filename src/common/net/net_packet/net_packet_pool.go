package net_packet

import (
	"sync"
)

var packetPool *NetPacketPool

func init() {
	packetPool = &NetPacketPool{}
	packetPool.Init()
}

func GetPool() *NetPacketPool {
	return packetPool
}

type NetPacketPool struct {
	pool sync.Pool
}

func (p *NetPacketPool) Init() {
	p.pool = sync.Pool{
		New: func() interface{} {
			return &NetPacket{}
		},
	}
}

func (p *NetPacketPool) Get() *NetPacket {
	data := p.pool.Get().(*NetPacket)
	return data
}

func (p *NetPacketPool) Put(netData *NetPacket) {
	// 将对象放回对象池中
	netData.Id = 0
	netData.Length = 0
	netData.Body = nil
	p.pool.Put(netData)
}
