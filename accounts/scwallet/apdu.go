// Copyright 2018 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

package scwallet

import (
	"bytes"
	"encoding/binary"
)

const (
	CLA_ISO7816 = 0

	INS_SELECT                = 0xA4
	INS_GET_RESPONSE          = 0xC0
	INS_PAIR                  = 0x12
	INS_UNPAIR                = 0x13
	INS_OPEN_SECURE_CHANNEL   = 0x10
	INS_MUTUALLY_AUTHENTICATE = 0x11

	SW1_GET_RESPONSE = 0x61
	SW1_OK           = 0x90
)

// CommandAPDU represents an application data unit sent to a smartcard
type CommandAPDU struct {
	Cla, Ins, P1, P2 uint8  // Class, Instruction, Parameter 1, Parameter 2
	Data             []byte // Command data
	Le               uint8  // Command data length
}

// serialize serializes a command APDU.
func (ca CommandAPDU) serialize() ([]byte, error) {
	buf := new(bytes.Buffer)

	if err := binary.Write(buf, binary.BigEndian, ca.Cla); err != nil {
		return nil, err
	}
	if err := binary.Write(buf, binary.BigEndian, ca.Ins); err != nil {
		return nil, err
	}
	if err := binary.Write(buf, binary.BigEndian, ca.P1); err != nil {
		return nil, err
	}
	if err := binary.Write(buf, binary.BigEndian, ca.P2); err != nil {
		return nil, err
	}
	if len(ca.Data) > 0 {
		if err := binary.Write(buf, binary.BigEndian, uint8(len(ca.Data))); err != nil {
			return nil, err
		}
		if err := binary.Write(buf, binary.BigEndian, ca.Data); err != nil {
			return nil, err
		}
	}
	if err := binary.Write(buf, binary.BigEndian, ca.Le); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// ResponseAPDU represents an application data unit received from a smart card
type ResponseAPDU struct {
	Data     []byte // response data
	Sw1, Sw2 uint8  // status words 1 and 2
}

// deserialize deserializes a response APDU
func (ra *ResponseAPDU) deserialize(data []byte) error {
	ra.Data = make([]byte, len(data)-2)

	buf := bytes.NewReader(data)
	if err := binary.Read(buf, binary.BigEndian, &ra.Data); err != nil {
		return err
	}
	if err := binary.Read(buf, binary.BigEndian, &ra.Sw1); err != nil {
		return err
	}
	if err := binary.Read(buf, binary.BigEndian, &ra.Sw2); err != nil {
		return err
	}
	return nil
}
