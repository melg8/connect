// SPDX-FileCopyrightText: 2024 Melg Eight <public.melg8@gmail.com>
//
// SPDX-License-Identifier: MIT

package connect

import (
	"fmt"
	"log"
	"net"

	"github.com/melg8/connect/internal/connect/helpers"
	fromauthserver "github.com/melg8/connect/internal/connect/packets/from_auth_server"
	"github.com/melg8/connect/internal/connect/packets/packet"
	toauthserver "github.com/melg8/connect/internal/connect/packets/to_auth_server"
)

func LogRecievedData(data []byte) {
	log.Println("Received: " + fmt.Sprint(len(data)) + " bytes")
	helpers.ShowAsHexAndASCII(data)
}

func LogSentData(data []byte) {
	log.Println("Sent: " + fmt.Sprint(len(data)) + " bytes")
	helpers.ShowAsHexAndASCII(data)
}

func LogInitPacket(initPacket *fromauthserver.InitPacket) {
	log.Println(initPacket.ToString())
}

// Extracts packet data from raw data. Decrypts if needed, returns packet id and packet data.
func ExtractPacketFromRawData(data []byte) (int32, []byte, error) {
	if len(data) < 3 {
		return 0, nil, fmt.Errorf("data is too small")
	}
	packetID := int32(data[2])
	if packetID == 0x00 {
		return packetID, data[3 : len(data)-4], nil
	}
	return 0, nil, fmt.Errorf("unexpected packet type")
}

// Reads full packet from connection.
func ReadPacket(conn net.Conn) ([]byte, error) {
	rawData := make([]byte, 1024)
	n, err := conn.Read(rawData)
	if err != nil {
		return nil, err
	}
	LogRecievedData(rawData[:n])

	return rawData[:n], nil
}

// Writes full packet to connection.
func WritePacket(conn net.Conn, data []byte) error {
	LogSentData(data)
	needToWrite := data
	for {
		n, err := conn.Write(needToWrite)
		if err != nil {
			return err
		}
		needToWrite = needToWrite[n:]
		if len(needToWrite) == 0 {
			break
		}
	}
	return nil
}

func RequestInit(rawData []byte) (*fromauthserver.InitPacket, error) {
	packetID, packetData, err := ExtractPacketFromRawData(rawData)
	if err != nil {
		return nil, err
	}
	if packetID != 0x00 {
		return nil, fmt.Errorf("unexpected packet type %v while waiting for init packet 0x00", packetID)
	}

	initPacket, err := fromauthserver.NewInitPacketFromBytes(packetData)
	if err != nil {
		return nil, err
	}
	LogInitPacket(initPacket)

	return initPacket, nil
}

func GGAuth(rawData []byte) (int, error) {
	packetID, packetData, err := ExtractPacketFromRawData(rawData)

	if err != nil {
		return 0, err
	}
	if packetID != 0x11 {
		return 0, fmt.Errorf("unexpected packet type %v while waiting for init packet 0x11", packetID)
	}

	log.Printf("got GGAuth packet with data:")
	helpers.ShowAsHexAndASCII(packetData)

	// ggAuthPacket, err := fromauthserver.NewGGAuthPacketFromBytes(packetData)
	// if err != nil {
	// 	return 0, err
	// }
	// LogGGAuthPacket(ggAuthPacket)

	// return ggAuthPacket, nil
	return 0, nil
}

func RequestGGAuth(conn net.Conn, initResponse *fromauthserver.InitPacket) (int, error) {
	requestGGAuth := toauthserver.NewDefaultRequestGGAuth(initResponse.SessionID)

	packetWriter := packet.NewWriter()
	err := requestGGAuth.ToBytes(packetWriter)
	packetData := packetWriter.Bytes()

	// assembler := toauthserver.NewAssembler()

	// requestGGAuth.ToWriter(assembler.Writer)

	// packetData, err = assembler.Assemble(0x11)

	// reserve 2 bytes for packet size
	// set packet id to proper value for each packet
	// insert packet data
	// reserve 0-7 bytes for padding depending on current size to be 8 byte aligned
	// claculate checksum from packet id to end of packet - 4 bytes
	// insert checksum to last 4 bytes
	// encrypt from 2 to end
	// calculate size of packet and insert to first 2 bytes

	if err != nil {
		return 0, err
	}
	err = WritePacket(conn, packetData)
	if err != nil {
		return 0, err
	}
	rawResponse, err := ReadPacket(conn)
	if err != nil {
		return 0, err
	}
	return GGAuth(rawResponse)
}

func AuthentificateConn(conn net.Conn) error {
	defer conn.Close()
	rawData, err := ReadPacket(conn)
	if err != nil {
		return err
	}
	initResponse, err := RequestInit(rawData)
	if err != nil {
		return err
	}

	_, err = RequestGGAuth(conn, initResponse)
	if err != nil {
		return err
	}

	// responseLogin, err := RequestLogin(conn, ggAuthResponse)
	// if err != nil {
	// 	return err
	// }

	// responseServerList, err := RequestServerList(conn, responseLogin)
	// if err != nil {
	// 	return err
	// }

	// responseServerLogin, err := RequestServerLogin(conn, responseServerList)
	// if err != nil {
	// 	return err
	// }
	return nil
}
