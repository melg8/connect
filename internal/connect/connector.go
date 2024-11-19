package connect

import (
	"fmt"
	"log"
	"net"
	"time"
)

type Connector interface {
	Connect() (net.Conn, error)
	Address() string
}

type TcpConnector struct {
	serverAddress string
	timeout       time.Duration
}

func NewTcpConnector(serverAddress string, timeout time.Duration) *TcpConnector {
	return &TcpConnector{
		serverAddress: serverAddress,
		timeout:       timeout,
	}
}

func (c *TcpConnector) Connect() (net.Conn, error) {
	return net.DialTimeout("tcp", c.serverAddress, c.timeout)
}

func (c *TcpConnector) Address() string {
	return c.serverAddress
}

type RateLimitedConnector struct {
	connector    Connector
	lastConnTime time.Time
	timeout      time.Duration
}

func NewRateLimitedConnector(connector Connector, timeout time.Duration) *RateLimitedConnector {
	return &RateLimitedConnector{
		connector:    connector,
		lastConnTime: time.Now().Add(-timeout),
		timeout:      timeout,
	}
}

func (c *RateLimitedConnector) Address() string {
	return c.connector.Address()
}

func (c *RateLimitedConnector) Connect() (net.Conn, error) {
	now := time.Now()

	if now.Sub(c.lastConnTime) < c.timeout {
		log.Printf("Sleeping for %v seconds before next attempt", c.timeout.Seconds())
		time.Sleep(c.timeout - time.Since(now))
	} else {
		log.Println("Connection successful, no delay needed")
	}

	conn, err := c.connector.Connect()
	if err != nil {
		log.Printf("Error connecting to server: %v", err)
		return nil, err
	}

	c.lastConnTime = now
	log.Printf("Connected to server: %s at %v", c.connector.Address(), now)
	return conn, nil
}

type RetryConnector struct {
	connector Connector
	retries   int
}

func NewRetryConnector(connector Connector, retries int) *RetryConnector {
	return &RetryConnector{connector: connector, retries: retries}
}

func (c *RetryConnector) Address() string {
	return c.connector.Address()
}

func (c *RetryConnector) Connect() (net.Conn, error) {
	for r := c.retries; r > 0; r-- {
		log.Printf("Attempting connection %d/%d", r, c.retries)
		conn, err := c.connector.Connect()
		if err != nil {
			log.Printf("Error connecting to server: %v (attempt %d)", err, r)
			if r > 1 {
				log.Println("Retrying connection...")
			}
			continue
		} else {
			return conn, nil
		}
	}
	return nil, fmt.Errorf("failed to connect to server after %d attempts", c.retries)
}

func ServerConnector(address string) (Connector, error) {
	tcpConnectorTimeout := time.Second * 10
	betweenAttemptsTimeout := time.Second + time.Millisecond*10
	tcpConnector := NewTcpConnector(address, tcpConnectorTimeout)
	connector := NewRateLimitedConnector(tcpConnector, betweenAttemptsTimeout)
	return NewRetryConnector(connector, 5), nil
}
