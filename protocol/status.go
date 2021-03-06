package protocol

/*
 * go-raknet
 *
 * Copyright (c) 2018 beito
 *
 * This software is released under the MIT License.
 * http://opensource.org/licenses/mit-license.php
 */

import (
	"github.com/beito123/go-raknet"
	"github.com/beito123/go-raknet/identifier"
)

type ConnectedPing struct {
	BasePacket
	Timestamp int64
}

func (pk ConnectedPing) ID() byte {
	return IDConnectedPing
}

func (pk *ConnectedPing) Encode() error {
	err := pk.BasePacket.Encode(pk)
	if err != nil {
		return err
	}

	err = pk.PutLong(pk.Timestamp)
	if err != nil {
		return err
	}

	return nil
}

func (pk *ConnectedPing) Decode() error {
	err := pk.BasePacket.Decode(pk)
	if err != nil {
		return err
	}

	pk.Timestamp, err = pk.Long()
	if err != nil {
		return err
	}

	return nil
}

func (pk *ConnectedPing) New() raknet.Packet {
	return new(ConnectedPing)
}

type ConnectedPong struct {
	BasePacket
	Timestamp     int64
	TimestampPong int64
}

func (pk ConnectedPong) ID() byte {
	return IDConnectedPong
}

func (pk *ConnectedPong) Encode() error {
	err := pk.BasePacket.Encode(pk)
	if err != nil {
		return err
	}

	err = pk.PutLong(pk.Timestamp)
	if err != nil {
		return err
	}

	err = pk.PutLong(pk.TimestampPong)
	if err != nil {
		return err
	}

	return nil
}

func (pk *ConnectedPong) Decode() error {
	err := pk.BasePacket.Decode(pk)
	if err != nil {
		return err
	}

	pk.Timestamp, err = pk.Long()
	if err != nil {
		return err
	}

	pk.TimestampPong, err = pk.Long()
	if err != nil {
		return err
	}

	return nil
}

func (pk *ConnectedPong) New() raknet.Packet {
	return new(ConnectedPong)
}

type UnconnectedPing struct {
	BasePacket
	Timestamp  int64
	Magic      bool
	PingID     int64
	Connection *raknet.ConnectionType
}

func (pk UnconnectedPing) ID() byte {
	return IDUnconnectedPing
}

func (pk *UnconnectedPing) Encode() error {
	err := pk.BasePacket.Encode(pk)
	if err != nil {
		return err
	}

	err = pk.PutLong(pk.Timestamp)
	if err != nil {
		return err
	}

	err = pk.PutMagic()
	if err != nil {
		return err
	}

	err = pk.PutLong(pk.PingID)
	if err != nil {
		return err
	}

	err = pk.PutConnectionType(pk.Connection)
	if err != nil {
		return err
	}

	return nil
}

func (pk *UnconnectedPing) Decode() error {
	err := pk.BasePacket.Decode(pk)
	if err != nil {
		return err
	}

	pk.Timestamp, err = pk.Long()
	if err != nil {
		return err
	}

	pk.Magic = pk.CheckMagic()

	pk.PingID, err = pk.Long()
	if err != nil {
		return err
	}

	pk.Connection, err = pk.ConnectionType()
	if err != nil {
		return err
	}

	return nil
}

func (pk *UnconnectedPing) New() raknet.Packet {
	return new(UnconnectedPing)
}

type UnconnectedPingOpenConnections struct {
	UnconnectedPing
}

func (pk UnconnectedPingOpenConnections) ID() byte {
	return IDUnconnectedPingOpenConnections
}

func (pk *UnconnectedPingOpenConnections) New() raknet.Packet {
	return new(UnconnectedPingOpenConnections)
}

type UnconnectedPong struct {
	BasePacket
	Timestamp  int64
	PongID     int64
	Magic      bool
	Identifier identifier.Identifier
	Connection *raknet.ConnectionType
}

func (pk UnconnectedPong) ID() byte {
	return IDUnconnectedPong
}

func (pk *UnconnectedPong) Encode() error {
	err := pk.BasePacket.Encode(pk)
	if err != nil {
		return err
	}

	err = pk.PutLong(pk.Timestamp)
	if err != nil {
		return err
	}

	err = pk.PutLong(pk.PongID)
	if err != nil {
		return err
	}

	err = pk.PutMagic()
	if err != nil {
		return err
	}

	err = pk.PutString(pk.Identifier.Build())
	if err != nil {
		return err
	}

	err = pk.PutConnectionType(pk.Connection)
	if err != nil {
		return err
	}

	return nil
}

func (pk *UnconnectedPong) Decode() error {
	err := pk.BasePacket.Decode(pk)
	if err != nil {
		return err
	}

	pk.Timestamp, err = pk.Long()
	if err != nil {
		return err
	}

	pk.PongID, err = pk.Long()
	if err != nil {
		return err
	}

	pk.Magic = pk.CheckMagic()

	id, err := pk.String()
	if err != nil {
		return err
	}

	pk.Connection, err = pk.ConnectionType()
	if err != nil {
		return err
	}

	pk.Identifier = identifier.Base{
		Identifier: id,
		Connection: pk.Connection,
	}

	return nil
}

func (pk *UnconnectedPong) New() raknet.Packet {
	return new(UnconnectedPong)
}
