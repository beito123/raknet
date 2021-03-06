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
	"sort"

	"github.com/beito123/go-raknet"
	"github.com/beito123/go-raknet/binary"
)

type ACKType int

const (
	TypeACK ACKType = iota
	TypeNACK
)

type Acknowledge struct {
	BasePacket

	Type ACKType

	Records []*raknet.Record
}

func (ack *Acknowledge) ID() byte {
	switch ack.Type {
	case TypeACK:
		return IDACK
	case TypeNACK:
		return IDNACK
	}

	return 0xff
}

func (ack *Acknowledge) New() raknet.Packet {
	return &Acknowledge{
		Type: ack.Type,
	}
}

func (ack *Acknowledge) Encode() error {
	err := ack.BasePacket.Encode(ack)
	if err != nil {
		return err
	}

	ack.Records = CondenseRecords(ack.Records)

	err = ack.PutShort(uint16(len(ack.Records)))
	if err != nil {
		return err
	}

	for _, rec := range ack.Records {
		noRange := !rec.IsRanged() // 0 = ranged, 1 = no ranged

		err = ack.PutBool(noRange)
		if err != nil {
			return err
		}

		err = ack.PutLTriad(binary.Triad(rec.Index))
		if err != nil {
			return err
		}

		if !noRange { // ranged
			err = ack.PutLTriad(binary.Triad(rec.EndIndex))
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (ack *Acknowledge) Decode() error {
	err := ack.BasePacket.Decode(ack)
	if err != nil {
		return err
	}

	recLen, err := ack.Short()
	if err != nil {
		return err
	}

	ack.Records = []*raknet.Record{}
	for i := 0; i < int(recLen); i++ {
		noRange, err := ack.Bool()
		if err != nil {
			return err
		}

		index, err := ack.LTriad()
		if err != nil {
			return err
		}

		var endIndex binary.Triad
		if noRange {
			endIndex, err = ack.LTriad()
			if err != nil {
				return err
			}
		}

		ack.Records = append(ack.Records, &raknet.Record{
			Index:    int(index),
			EndIndex: int(endIndex),
		})
	}

	ack.Records = simplifyRecords(ack.Records)

	return nil
}

// CondenseRecords returns condensed records.
// For example (No need sort): [0, 2, 3, 5, 8, 9, 10, 15] -> [0, [2:3], 5, [8:10], 15]
func CondenseRecords(records []*raknet.Record) []*raknet.Record {
	var ids []int
	for _, record := range records {
		ids = append(ids, record.Numbers()...)
	}

	sort.Ints(ids)

	ln := len(ids)

	var nRecords []*raknet.Record
	for i := 0; i < ln; i++ {
		rec := ids[i]
		last := rec

		// find
		if i+1 < ln {
			for last+1 == ids[i+1] {
				last = ids[i+1]
				i++
				if i+1 >= ln {
					break
				}
			}
		}

		end := last

		if rec == end { // no ranged
			nRecords = append(nRecords, &raknet.Record{
				Index: rec,
			})
		} else { // ranged
			nRecords = append(nRecords, &raknet.Record{
				Index:    rec,
				EndIndex: end,
			})
		}
	}

	return nRecords
}

func simplifyRecords(records []*raknet.Record) []*raknet.Record {
	var ids []int
	for _, rec := range records {
		ids = append(ids, rec.Numbers()...)
	}

	recs := make([]*raknet.Record, len(ids))
	for i := 0; i < len(ids); i++ {
		recs[i] = &raknet.Record{
			Index: ids[i],
		}
	}

	return recs
}
