// SPDX-FileCopyrightText: 2024 Melg Eight <public.melg8@gmail.com>
//
// SPDX-License-Identifier: MIT

package connect

import (
	"fmt"
	"io"
	"net"

	"github.com/melg8/connect/internal/connect/helpers"
	"github.com/melg8/connect/internal/connect/packets/from_auth_server"
)

type AuthProtocol struct {
	dataChannel chan []byte
	conn        net.Conn
}

func NewAuthProtocol(dataChannel chan []byte, conn net.Conn) *AuthProtocol {
	return &AuthProtocol{dataChannel: dataChannel, conn: conn}
}

func (p *AuthProtocol) Run() {
	for data := range p.dataChannel {
		helpers.ShowAsHexView(data)
		helpers.ShowAsHexAndAscii(data)
		initPacketData := data[3 : len(data)-4]
		initPacket, err := from_auth_server.NewInitPacketFromBytes(initPacketData)
		if err == nil {
			fmt.Println(initPacket.ToString())
		}
	}
}

type AuthDataReciever struct {
	conn        net.Conn
	dataChannel chan []byte
}

func NewAuthDataReciever(conn net.Conn) *AuthDataReciever {
	return &AuthDataReciever{conn: conn}
}

func (r *AuthDataReciever) Run() {
	for {
		for {
			buf := make([]byte, 1024)
			fmt.Println("Reading from connection")
			n, err := r.conn.Read(buf)
			fmt.Println("Read " + fmt.Sprint(n) + " bytes")
			if err == io.EOF {
				fmt.Println("Server doesn't send any more data")
				break
			} else if err != nil {
				fmt.Printf("Error reading from connection: %v\n", err)
				return
			}
			fmt.Println("Received packet from login server:")
			r.dataChannel <- buf[:n]
		}
		fmt.Println("Connection closed")
	}
}

func (r *AuthDataReciever) Stop() {
	r.conn.Close()
}

func Authentificate(conn net.Conn) {
	defer conn.Close()
	fmt.Println("Connected to server. Waiting for data")
	dataChannel := make(chan []byte)
	protocol := NewAuthProtocol(dataChannel, conn)
	go protocol.Run()

	for {
		buf := make([]byte, 1024)
		fmt.Println("Reading from connection")
		n, err := conn.Read(buf)
		fmt.Println("Read " + fmt.Sprint(n) + " bytes")
		if err == io.EOF {
			fmt.Println("Server doesn't send any more data")
			break
		} else if err != nil {
			fmt.Printf("Error reading from connection: %v\n", err)
			return
		}
		fmt.Println("Received packet from login server:")
		dataChannel <- buf[:n]
	}
	fmt.Println("Connection closed")
}
