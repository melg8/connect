// SPDX-FileCopyrightText: © 2024 Melg Eight <public.melg8@gmail.com>
//
// SPDX-License-Identifier: MIT

package connect

import (
	"fmt"
	"io"
	"net"
	"time"
)

// Define connection type data structure:
type Connection struct {
	net.Conn
}

type AuthProtocol struct {
	dataChannel chan []byte
	conn        net.Conn
}

func NewAuthProtocol(dataChannel chan []byte, conn net.Conn) *AuthProtocol {
	return &AuthProtocol{dataChannel: dataChannel, conn: conn}
}

func (p *AuthProtocol) Run() {
	for data := range p.dataChannel {
		ShowAsHexAndAscii(data)
	}
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

type ConnectionRequest struct {
	conn chan net.Conn
}

// Prevents too frequent connections to the login server from the same IP from
// triggering flood protection.
// https://gitlab.com/TheDnR/l2j-lisvus/-/blame/main/
// core/java/net/sf/l2j/util/IPv4Filter.java#L88

type AuthServerConnectionProvider struct {
	serverAddress string
	requests      chan ConnectionRequest
	done          <-chan struct{}
	lastConnTime  time.Time
	timeout       time.Duration
}

func NewAuthServerConnectionProvider(serverAddress string) *AuthServerConnectionProvider {
	return &AuthServerConnectionProvider{
		serverAddress: serverAddress,
		requests:      make(chan ConnectionRequest, 100),
		lastConnTime:  time.Now().Add(time.Second * -2),
		timeout:       time.Second + time.Millisecond*100,
	}
}

func (p *AuthServerConnectionProvider) Start() {
	go func() {
		for request := range p.requests {
			select {
			case <-p.done:
				return
			default:
				p.openConnection(request.conn)
			}
		}
	}()
}

func (p *AuthServerConnectionProvider) RequestConnection() chan net.Conn {
	requested := make(chan net.Conn)
	p.requests <- ConnectionRequest{conn: requested}
	return requested
}

func (p *AuthServerConnectionProvider) openConnection(requested_conn chan net.Conn) {
	for {
		select {
		case <-p.done:
			return
		default:
			now := time.Now()
			fmt.Println("Passed " + now.Sub(p.lastConnTime).String() +
				" from last attempt to connect to server")
			if now.Sub(p.lastConnTime) >= p.timeout {
				fmt.Println("Connecting to server: " + p.serverAddress + " time: " + now.String())
				conn, err := net.DialTimeout("tcp", p.serverAddress, time.Second)
				if err != nil {
					fmt.Printf("Error connecting to server: %v\n", err)
					time.Sleep(p.timeout)
					continue
				}
				p.lastConnTime = now
				fmt.Println("Connected to server: " + p.serverAddress)
				requested_conn <- conn
				fmt.Println("Connection to server send to requestor")
				return
			} else {
				fmt.Println("Sleeping for 1 sec")
				time.Sleep(p.timeout)
			}
		}
	}
}

type GameServerConnector struct {
	auth_provider *AuthServerConnectionProvider
	conn_retries  int
}

func NewGameServerConnector(auth_provider *AuthServerConnectionProvider) *GameServerConnector {
	return &GameServerConnector{auth_provider: auth_provider, conn_retries: 2}
}

func (c *GameServerConnector) Start() {
	go func() {
		for r := c.conn_retries; r > 0; r-- {
			fmt.Println("Requesting connection, retries left: " + fmt.Sprint(r))
			conn := c.auth_provider.RequestConnection()
			Authentificate(<-conn)
		}
	}()
}

func StartBotsAt(address string, count int) {
	provider := NewAuthServerConnectionProvider(address)
	connectors := make([]*GameServerConnector, count)
	for i := 0; i < count; i++ {
		connectors[i] = NewGameServerConnector(provider)
		connectors[i].Start()
	}
	provider.Start()
}
