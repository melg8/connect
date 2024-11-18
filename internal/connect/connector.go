// SPDX-FileCopyrightText: © 2024 Melg Eight <public.melg8@gmail.com>
//
// SPDX-License-Identifier: MIT

package connect

import (
	"errors"
	"fmt"
	"net"
	"time"
)

type Connector interface {
	Connect() (net.Conn, error)
	Address() string
}

type TcpConnector struct {
	serverAddress string
}

func NewTcpConnector(serverAddress string) *TcpConnector {
	return &TcpConnector{
		serverAddress: serverAddress,
	}
}

func (c *TcpConnector) Connect() (net.Conn, error) {
	return net.DialTimeout("tcp", c.serverAddress, time.Second)
}

func (c *TcpConnector) Address() string {
	return c.serverAddress
}

// Prevents too frequent connections to the login server from the same IP from
// triggering flood protection.
// https://gitlab.com/TheDnR/l2j-lisvus/-/blame/main/
// core/java/net/sf/l2j/util/IPv4Filter.java#L88

type RateLimitedConnector struct {
	connector    Connector
	lastConnTime time.Time
	timeout      time.Duration
}

func NewRateLimitedConnector(connector Connector) *RateLimitedConnector {
	return &RateLimitedConnector{
		connector:    connector,
		lastConnTime: time.Now().Add(time.Second * -2),
		timeout:      time.Second + time.Millisecond*100,
	}
}

func (c *RateLimitedConnector) Address() string {
	return c.connector.Address()
}

func (c *RateLimitedConnector) Connect() (net.Conn, error) {
	now := time.Now()
	fmt.Println("Passed " + now.Sub(c.lastConnTime).String() +
		" from last attempt to connect to server")
	if now.Sub(c.lastConnTime) < c.timeout {
		fmt.Println("Sleeping for 1 sec")
		time.Sleep(c.timeout)
	}
	fmt.Println("Connecting to server: " + c.connector.Address() + " time: " + now.String())
	conn, err := c.connector.Connect()
	if err != nil {
		fmt.Printf("Error connecting to server: %v\n", err)
		return nil, err
	}
	c.lastConnTime = now
	fmt.Println("Connected to server: " + c.connector.Address() + " time: " + now.String())
	return conn, nil
}

type RetryConnector struct {
	connector Connector
	retries   int
}

func NewRetriesConnector(connector Connector, retries int) *RetryConnector {
	return &RetryConnector{connector: connector, retries: retries}
}

func (c *RetryConnector) Address() string {
	return c.connector.Address()
}

func (c *RetryConnector) Connect() (net.Conn, error) {
	for r := c.retries; r > 0; r-- {
		fmt.Println("Requesting connection, retries left: " + fmt.Sprint(r))
		conn, err := c.connector.Connect()
		if err != nil {
			fmt.Printf("Error connecting to server: %v\n", err)
			if r > 1 {
				fmt.Println("Retrying connection...")
			}
			continue
		} else {
			return conn, nil
		}
	}
	return nil, errors.New("Failed to connect to server: " + c.Address() +
		"after " + fmt.Sprint(c.retries) + " retries")
}

func ServerConnector(address string) *RetryConnector {
	if address == "" {
		return nil
	}
	return NewRetriesConnector(NewRateLimitedConnector(NewTcpConnector(address)), 5)
}
