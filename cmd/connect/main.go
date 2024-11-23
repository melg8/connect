// SPDX-FileCopyrightText: 2024 Melg Eight <public.melg8@gmail.com>
//
// SPDX-License-Identifier: MIT

package main

import (
	"fmt"
	"log"

	"github.com/melg8/connect/internal/connect"
)

func connectAndAuthenticate() error {
	connector, err := connect.ServerConnector("127.0.0.1:2106")
	if err != nil {
		return fmt.Errorf("failed to create server connector: %v", err)
	}

	conn, err := connector.Connect()
	if err != nil {
		return fmt.Errorf("failed to connect to server: %v", err)
	}
	defer conn.Close()

	if err := connect.AuthentificateConn(conn); err != nil {
		return fmt.Errorf("failed to authentificate connection: %v", err)
	}

	log.Println("Connection authentificated")
	return nil
}

func main() {
	log.Println("Starting connect bot...")
	if err := connectAndAuthenticate(); err != nil {
		log.Fatal(err)
	}
}
