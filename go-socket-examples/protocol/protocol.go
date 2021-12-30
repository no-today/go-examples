package protocol

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/cloudwego/netpoll"
	"io"
)

func (p *Message) Pack() ([]byte, error) {
	buffer := new(bytes.Buffer)

	// 内容长度
	// +1 是为了协议版本号
	if err := binary.Write(buffer, binary.LittleEndian, int32(len(p.content))+1); err != nil {
		return nil, err
	}

	// 协议版本号
	if err := binary.Write(buffer, binary.LittleEndian, p.version); err != nil {
		return nil, err
	}

	// 内容
	if err := binary.Write(buffer, binary.LittleEndian, p.content); err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}

func UnPack(reader io.Reader) (*Message, error) {
	buffer := make([]byte, 4)
	_, err := io.ReadFull(reader, buffer)
	if err != nil {
		return nil, err
	}

	// 解码内容长度
	contentLength := binary.LittleEndian.Uint32(buffer)

	buffer = make([]byte, contentLength)
	_, err = io.ReadFull(reader, buffer)
	if err != nil {
		return nil, err
	}

	return &Message{
		version: buffer[0],
		content: buffer[1:],
	}, nil
}

func UnPackZeroCopy(reader netpoll.Reader) (*Message, error) {
	// 读取内容长度
	buffer, err := reader.Next(4)
	if err != nil {
		return nil, err
	}

	// 解码内容长度
	contentLength := binary.LittleEndian.Uint32(buffer)

	// 读取版本号
	version, err := reader.Next(1)
	if err != nil {
		return nil, err
	}

	// 读取内容
	buffer, err = reader.Next(int(contentLength - 1))
	if err != nil {
		return nil, err
	}

	return &Message{
		version: version[0],
		content: buffer,
	}, nil
}

type Message struct {
	version uint8
	content []byte
}

func NewMessage(version uint8, content []byte) *Message {
	return &Message{version: version, content: content}
}

func (p *Message) Version() uint8 {
	return p.version
}

func (p *Message) SetVersion(version uint8) *Message {
	p.version = version
	return p
}

func (p *Message) Content() []byte {
	return p.content
}

func (p *Message) SetContent(content []byte) *Message {
	p.content = content
	return p
}

func (p *Message) String() string {
	return fmt.Sprintf("[%d/%d] %s", p.version, len(p.content), string(p.content))
}
