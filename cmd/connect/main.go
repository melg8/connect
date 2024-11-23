// SPDX-FileCopyrightText: 2024 Melg Eight <public.melg8@gmail.com>
//
// SPDX-License-Identifier: MIT

package main

import (
	"log"

	"github.com/melg8/connect/internal/connect"
)

func main() {
	log.Println("Starting connect bot...")

	connector, err := connect.ServerConnector("127.0.0.1:2106")
	if err != nil {
		log.Fatalf("Failed to create server connector: %v", err)
	}

	conn, err := connector.Connect()
	if err != nil {
		log.Fatalf("Failed to connect to server: %v", err)
	}

	if err := connect.AuthentificateConn(conn); err != nil {
		conn.Close()
		log.Fatal("Failed to authentificate connection: ", err)
	}
	log.Println("Connection authentificated")
	conn.Close()
}
