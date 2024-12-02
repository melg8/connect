// SPDX-FileCopyrightText: 2024 Melg Eight <public.melg8@gmail.com>
//
// SPDX-License-Identifier: MIT

package packet

// ## RegisterNewPacket(packetID int32, packetSize int32)
// reserve 2 bytes for packet size
// set packet id to proper value for each packet

// ## Write...(...)
// insert packet data

// ## Assemble()
// reserve 0-7 bytes for padding depending on current size to be 8 byte aligned
// claculate checksum from packet id to end of packet - 4 bytes
// insert checksum to last 4 bytes
// encrypt from 2 to end
// calculate size of packet and insert to first 2 bytes

// ## Packet()

type Assembler struct {
	Writer
}

func NewAssembler() *Assembler {
	return &Assembler{}
}
