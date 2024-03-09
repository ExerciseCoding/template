package net

import (
	"encoding/binary"
	"net"
	"testing"
)

func TestServer(t *testing.T) {
	// 创建监听
	listen, err := net.Listen("tcp", ":8081")
	if err != nil {
		return
	}
	// 接收请求
	for {
		conn, err := listen.Accept()
		if err != nil {
			return
		}
		go func() {
			handle(conn)
		}()
	}

}

func handle(conn net.Conn) {
	for {
		lenBs := make([]byte, 8)
		_, err := conn.Read(lenBs)
		if err != nil {
			conn.Close()
			return
		}
		msgLen := binary.BigEndian.Uint64(lenBs)
		reqBs := make([]byte, msgLen)
		_, err = conn.Read(reqBs)
		if err != nil {
			conn.Close()
			return
		}
		_, err = conn.Write([]byte("hello, world"))
		if err != nil {
			conn.Close()
			return
		}
	}

}

//func handle(conn net.Conn) {
//	go func() {
//		resp <- ch
//		_, errs = conn.Write([]byte("hello, world"))
//		if errs != nil {
//			conn.Close()
//			return
//		}
//	}()
//
//	for {
//		lenBs := make([]byte, 8)
//		_, errs := conn.Read(lenBs)
//		if errs != nil {
//			conn.Close()
//			return
//		}
//		msgLen := binary.BigEndian.Uint64(lenBs)
//		reqBs := make([]byte, msgLen)
//		_, errs = conn.Read(reqBs)
//		if errs != nil {
//			conn.Close()
//			return
//		}
//
//		go func() {
//			resp := handleReq(reqBs)
//			ch <- resp
//		}()
//
//	}
//
//}
