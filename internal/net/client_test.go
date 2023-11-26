package net

import (
	"encoding/binary"
	"net"
	"testing"
)

func TestClient(t *testing.T) {
	conn, err := net.Dial("tcp", ":8081")
	if err != nil {
		t.Fatal(err)
	}

	//字节序，又称端序或尾序（英语中用单词：Endianness 表示），在计算机领域中，指电脑内存中或在数字通信链路中，占用多个字节的数据的字节排列顺序。
	//字节的排列方式有两个通用规则:
	//大端序（Big-Endian）将数据的低位字节存放在内存的高位地址，高位字节存放在低位地址。这种排列方式与数据用字节表示时的书写顺序一致，符合人类的阅读习惯。
	//小端序（Little-Endian），将一个多位数的低位放在较小的地址处，高位放在较大的地址处，则称小端序。小端序与人类的阅读习惯相反，但更符合计算机读取内存的方式，因为CPU读取内存中的数据时，是从低地址向高地址方向进行读取的。

	msg := "how are you"
	msgLen := len(msg)
	msgLenBs := make([]byte, 8)
	binary.BigEndian.PutUint64(msgLenBs, uint64(msgLen))
	data := append(msgLenBs, []byte(msg)...)
	_, err = conn.Write(data)
	if err != nil {
		conn.Close()
		return
	}
	respBs := make([]byte, 16)
	_, err = conn.Read(respBs)
	if err != nil {
		conn.Close()
	}

}
