package tls_disselector

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io"

	"github.com/superstes/calamary/log"
)

const (
	RecordHeaderLen = 5
)

// record content type
const (
	ChangeCipherSpec = 0x14
	EncryptedAlert   = 0x15
	Handshake        = 0x16
	AppData          = 0x17
)

var (
	ErrBadType = errors.New("bad type")
)

type Version uint16

type Record struct {
	Type    uint8
	Version Version
	Opaque  []byte
}

func ReadRecord(r io.Reader) (*Record, error) {
	record := &Record{}
	if _, err := record.ReadFrom(r); err != nil {
		return nil, err
	}
	return record, nil
}

func (rec *Record) ReadFrom(r io.Reader) (n int64, err error) {
	b := make([]byte, RecordHeaderLen)
	nn, err := io.ReadFull(r, b)
	n += int64(nn)
	if err != nil {
		return
	}
	rec.Type = b[0]
	rec.Version = Version(binary.BigEndian.Uint16(b[1:3]))
	length := int(binary.BigEndian.Uint16(b[3:5]))
	rec.Opaque = make([]byte, length)
	nn, err = io.ReadFull(r, rec.Opaque)
	n += int64(nn)
	return
}

func (rec *Record) WriteTo(w io.Writer) (n int64, err error) {
	buf := &bytes.Buffer{}
	buf.WriteByte(rec.Type)
	err = binary.Write(buf, binary.BigEndian, rec.Version)
	if err != nil {
		log.Warn("proc-parse-tls", "Hello error")
	}
	err = binary.Write(buf, binary.BigEndian, uint16(len(rec.Opaque)))
	if err != nil {
		log.Warn("proc-parse-tls", "Hello error")
	}
	buf.Write(rec.Opaque)
	return buf.WriteTo(w)
}
