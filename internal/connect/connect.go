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

func Connect() {
	address := "127.0.0.1:2106"

	fmt.Println("Connecting to server: " + address)
	conn, err := net.DialTimeout("tcp", address, time.Second)
	if err != nil {
		fmt.Printf("Error connecting to server: %v\n", err)
		return
	}
	defer conn.Close()

	fmt.Println("Connected to server. Receiving data...")

	// TODO(melg): Maybe wrong, and should do it in loop?
	buf := make([]byte, 1024)
	for {
		n, err := conn.Read(buf)
		if err == io.EOF {
			fmt.Println("Server doesn't send any more data.")
			break
		} else if err != nil {
			fmt.Printf("Error reading from connection: %v\n", err)
			return
		}
		ShowAsHexAndAscii(buf[:n])
	}
	fmt.Println("Connection closed.")
}
