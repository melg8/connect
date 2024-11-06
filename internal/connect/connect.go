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

func Connect(address string) {
	fmt.Println("Connecting to server: " + address)
	conn, err := net.DialTimeout("tcp", address, time.Second)
	if err != nil {
		fmt.Printf("Error connecting to server: %v\n", err)
		return
	}
	defer conn.Close()

	fmt.Println("Connected to server. Waiting for data")

	// TODO(melg): Maybe wrong, and should do it in loop?

	for {
		buf := make([]byte, 1024)
		n, err := conn.Read(buf)
		if err == io.EOF {
			fmt.Println("Server doesn't send any more data")
			break
		} else if err != nil {
			fmt.Printf("Error reading from connection: %v\n", err)
			return
		}

		fmt.Println("Received packet from login server:")
		ShowAsHexAndAscii(buf[:n])
	}
	fmt.Println("Connection closed")
}
