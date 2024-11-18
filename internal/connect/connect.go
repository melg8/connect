// SPDX-FileCopyrightText: © 2024 Melg Eight <public.melg8@gmail.com>
//
// SPDX-License-Identifier: MIT

package connect

import (
	"fmt"
	"io"
	"net"
	"time"

	"github.com/melg8/connect/internal/connect/helpers"
	"github.com/melg8/connect/internal/connect/packets/from_auth_server"
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

type ConnectionRequest struct {
	conn chan net.Conn
}

// Prevents too frequent connections to the login server from the same IP from
// triggering flood protection.
// https://gitlab.com/TheDnR/l2j-lisvus/-/blame/main/
// core/java/net/sf/l2j/util/IPv4Filter.java#L88

type RateLimitedConnector struct {
	serverAddress string
	requests      chan ConnectionRequest
	done          chan struct{}
	lastConnTime  time.Time
	timeout       time.Duration
}

func NewRateLimitedConnector(serverAddress string) *RateLimitedConnector {
	return &RateLimitedConnector{
		serverAddress: serverAddress,
		requests:      make(chan ConnectionRequest, 100),
		lastConnTime:  time.Now().Add(time.Second * -2),
		timeout:       time.Second + time.Millisecond*100,
	}
}

func (c *RateLimitedConnector) Start() {
	go func() {
		for request := range c.requests {
			select {
			case <-c.done:
				return
			default:
				c.openConnection(request.conn)
			}
		}
	}()
}

func (c *RateLimitedConnector) Stop() {
	go func() {
		c.done <- struct{}{}
	}()
}

func (c *RateLimitedConnector) RequestConnection() chan net.Conn {
	requested := make(chan net.Conn)
	c.requests <- ConnectionRequest{conn: requested}
	return requested
}

func (c *RateLimitedConnector) openConnection(requested_conn chan net.Conn) {
	for {
		select {
		case <-c.done:
			return
		default:
			now := time.Now()
			fmt.Println("Passed " + now.Sub(c.lastConnTime).String() +
				" from last attempt to connect to server")
			if now.Sub(c.lastConnTime) >= c.timeout {
				fmt.Println("Connecting to server: " + c.serverAddress + " time: " + now.String())
				conn, err := net.DialTimeout("tcp", c.serverAddress, time.Second)
				if err != nil {
					fmt.Printf("Error connecting to server: %v\n", err)
					time.Sleep(c.timeout)
					continue
				}
				c.lastConnTime = now
				fmt.Println("Connected to server: " + c.serverAddress)
				requested_conn <- conn
				fmt.Println("Connection to server send to requestor")
				return
			} else {
				fmt.Println("Sleeping for 1 sec")
				time.Sleep(c.timeout)
			}
		}
	}
}

type GameServerConnector struct {
	auth_provider *RateLimitedConnector
	conn_retries  int
}

func NewGameServerConnector(auth_provider *RateLimitedConnector) *GameServerConnector {
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

func (c *GameServerConnector) Stop() {
	c.auth_provider.done <- struct{}{}
}

type MultiAuth struct {
	provider   *RateLimitedConnector
	connectors []*GameServerConnector
}

func NewMultiAuth(address string, count int) *MultiAuth {
	provider := NewRateLimitedConnector(address)
	connectors := make([]*GameServerConnector, count)
	for i := 0; i < count; i++ {
		connectors[i] = NewGameServerConnector(provider)
	}
	return &MultiAuth{provider: provider, connectors: connectors}
}

func (multiAuth *MultiAuth) Start() {
	for _, connector := range multiAuth.connectors {
		connector.Start()
	}
	multiAuth.provider.Start()
}

func (m *MultiAuth) Stop() {
	m.provider.Stop()
	for _, connector := range m.connectors {
		connector.Stop()
	}
}

func StartBotsAt(address string, count int) *MultiAuth {
	multiAuth := NewMultiAuth(address, count)
	multiAuth.Start()
	return multiAuth
}
