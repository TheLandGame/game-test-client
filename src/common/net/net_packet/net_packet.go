package net_packet

import (
	"bytes"
	"encoding/binary"
	"fmt"

	"github.com/golang/snappy"
)

type NetPacket struct {
	Id     uint32
	Length int
	Body   []byte
}

const (
	PACKET_ID_SIZE     = 4
	PACKET_LENGTH_SIZE = 2
	PACKET_LIMIT       = 8192
)

// 压缩字节流数据
func compression(serialized []byte) []byte {
	compressed := snappy.Encode(nil, serialized)
	return compressed
}

// 解压缩字节流
func decompress(decompressed []byte) ([]byte, error) {
	body, err := snappy.Decode(nil, decompressed)
	return body, err
}

// 编码
func EncodeNetPacket(id uint32, data []byte) []byte {
	// 计算数据包长度
	length := uint16(len(data))

	// 封装数据包
	buffer := new(bytes.Buffer)
	binary.Write(buffer, binary.LittleEndian, id)
	binary.Write(buffer, binary.LittleEndian, length)
	binary.Write(buffer, binary.LittleEndian, data)
	return buffer.Bytes()
}

// 合并包(封装数据包)
func WritePacket(netPackets []*NetPacket) (bodyList [][]byte) {
	buffer := new(bytes.Buffer)
	for _, packet := range netPackets {
		body := EncodeNetPacket(packet.Id, packet.Body)

		buffLength := buffer.Len()
		if buffLength > 0 && buffLength+len(body) >= PACKET_LIMIT {
			compression := compression(buffer.Bytes())
			bodyList = append(bodyList, compression)
			buffer = new(bytes.Buffer)
		}

		binary.Write(buffer, binary.LittleEndian, body)
		GetPool().Put(packet)
	}
	if buffer.Len() > 0 {
		compression := compression(buffer.Bytes())
		bodyList = append(bodyList, compression)
	}
	return bodyList
}

// 解码
func DecodeNetPacket(body []byte) (*NetPacket, error) {
	if len(body) < (PACKET_ID_SIZE + PACKET_LENGTH_SIZE) {
		return nil, fmt.Errorf("net data length invalid")
	}

	typeBody := body[:PACKET_ID_SIZE]
	typeId := binary.LittleEndian.Uint32(typeBody)

	body = body[PACKET_ID_SIZE:]
	l := int(binary.LittleEndian.Uint16(body[:PACKET_LENGTH_SIZE]))

	body = body[PACKET_LENGTH_SIZE:]
	if len(body) < l {
		return nil, fmt.Errorf("invalid net data")
	}

	packet := GetPool().Get()
	packet.Id = typeId
	packet.Length = l
	packet.Body = body[:l]
	return packet, nil
}

// 读取数据包
func ReadPacket(netData []byte) (netPackets []*NetPacket, err error) {
	body, err := decompress(netData)
	if err != nil {
		return nil, err
	}

	packMinLength := PACKET_ID_SIZE + PACKET_LENGTH_SIZE
	if len(body) < (packMinLength) {
		return nil, fmt.Errorf("invalid net data length = %d", len(body))
	}

	for len(body) >= packMinLength {
		packet, err := DecodeNetPacket(body)
		if err != nil {
			return nil, err
		}
		netPackets = append(netPackets, packet)
		index := packMinLength + packet.Length
		body = body[index:]
	}
	return
}
