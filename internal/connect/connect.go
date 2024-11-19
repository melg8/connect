// SPDX-FileCopyrightText: 2024 Melg Eight <public.melg8@gmail.com>
//
// SPDX-License-Identifier: MIT

package connect

import (
	"errors"
	"fmt"
	"io"
	"log"
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
		helpers.ShowAsHexAndASCII(data)
		initPacketData := data[3 : len(data)-4]
		initPacket, err := from_auth_server.NewInitPacketFromBytes(initPacketData)
		if err == nil {
			log.Println(initPacket.ToString())
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
			log.Println("Reading from connection")
			bytesRead, err := r.conn.Read(buf)
			log.Println("Read " + fmt.Sprint(bytesRead) + " bytes")
			if errors.Is(err, io.EOF) {
				log.Println("Server doesn't send any more data")
				break
			} else if err != nil {
				log.Printf("Error reading from connection: %v\n", err)
				return
			}
			log.Println("Received packet from login server:")
			r.dataChannel <- buf[:bytesRead]
		}
		log.Println("Connection closed")
	}
}

func (r *AuthDataReciever) Stop() {
	r.conn.Close()
}

func Authentificate(conn net.Conn) {
	defer conn.Close()
	log.Println("Connected to server. Waiting for data")
	dataChannel := make(chan []byte)
	protocol := NewAuthProtocol(dataChannel, conn)

	go protocol.Run()

	for {
		buf := make([]byte, 1024)
		log.Println("Reading from connection")
		bytesRead, err := conn.Read(buf)
		log.Println("Read " + fmt.Sprint(bytesRead) + " bytes")
		if errors.Is(err, io.EOF) {
			log.Println("Server doesn't send any more data")
			break
		} else if err != nil {
			log.Printf("Error reading from connection: %v\n", err)
			return
		}
		log.Println("Received packet from login server:")
		dataChannel <- buf[:bytesRead]
	}
	log.Println("Connection closed")
}
