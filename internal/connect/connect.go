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

func ShowAsHex(data []byte) {
	bytes_in_row := 16
	fmt.Printf("Got data [%d]:\n", len(data))

	for i := 0; i < len(data); i += bytes_in_row {
		for j := 0; j < bytes_in_row; j++ {
			if (i + j) < len(data) {
				fmt.Printf("%02x ", data[i+j])
			} else {
				fmt.Printf("   ")
			}
		}
		fmt.Println()
	}
	fmt.Println()
}

func ShowAsHexAndAscii(data []byte) {
	bytesInRow := 16
	fmt.Printf("Got data [%d]:\n", len(data))

	for i := 0; i < len(data); i += bytesInRow {
		// Print the offset
		fmt.Printf("%04x: ", i)

		// Print hex values
		for j := 0; j < bytesInRow; j++ {
			if (i + j) < len(data) {
				fmt.Printf("%02x ", data[i+j])
			} else {
				fmt.Printf("   ")
			}
		}

		// Print ASCII representation
		fmt.Print(" ")
		for j := 0; j < bytesInRow; j++ {
			if (i + j) < len(data) {
				// Printable ASCII range
				if data[i+j] >= 32 && data[i+j] <= 126 {
					fmt.Printf("%c", data[i+j])
				} else {
					fmt.Print(".")
				}
			}
		}

		fmt.Println()
	}
	fmt.Println()
}

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
