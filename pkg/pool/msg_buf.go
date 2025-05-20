/*
 * Copyright (C) 2020-2022, IrineSistiana
 *
 * This file is part of mosdns.
 *
 * mosdns is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * mosdns is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <https://www.gnu.org/licenses/>.
 */

package pool

import (
	"encoding/binary"
	"fmt"

	"github.com/miekg/dns"
)

const (
	// There is no such way to give dns.Msg.PackBuffer() a buffer
	// with a proper size.
	// Just give it a big buf and hope the buf will be reused in most scenes.
	packBufferSize = 8191

	uint16Size   = 2
	dummyPending = 10
)

// PackBuffer packs the dns msg m to wire format.
// Callers should release the buf by calling ReleaseBuf after they have done
// with the wire []byte.
func PackBuffer(m *dns.Msg) (*MsgBuffer, error) {
	buf, wire, err := packBuffer(m)
	if err != nil {
		return nil, err
	}
	defer ReleaseBuf(buf)

	msgBuf := GetBuf(len(wire))
	copy(*msgBuf, wire)
	return msgBuf, nil
}

// PackTCPBuffer packs the dns msg m to wire format, with to bytes length header.
// Callers should release the buf by calling ReleaseBuf.
func PackTCPBuffer(m *dns.Msg) (*MsgBuffer, error) {
	packBuf, wire, err := packBuffer(m)
	if err != nil {
		return nil, err
	}
	defer ReleaseBuf(packBuf)

	l := len(wire)
	if l > dns.MaxMsgSize {
		return nil, fmt.Errorf("dns payload size %d is too large", l)
	}

	msgBuf := GetBuf(uint16Size + len(wire))
	binary.BigEndian.PutUint16(*msgBuf, uint16(l))
	copy((*msgBuf)[uint16Size:], wire)
	return msgBuf, nil
}

func packBuffer(m *dns.Msg) (*MsgBuffer, []byte, error) {
	packBuf := GetBuf(m.Len() + dummyPending)
	wire, err := m.PackBuffer(*packBuf)
	if err != nil {
		ReleaseBuf(packBuf)
		return nil, nil, err
	}

	return packBuf, wire, nil
}
