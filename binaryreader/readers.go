package binaryreader

import (
	"encoding/binary"
	"io"
)

// ReadByte reads a single byte from reader, and stores it in value
func ReadByte(reader io.Reader, value *uint8) {
	binary.Read(reader, binary.LittleEndian, value)
}

// ReadWord reads a uint16 from reader, and stores it in value
func ReadWord(reader io.Reader, value *uint16) {
	binary.Read(reader, binary.LittleEndian, value)
}

// ReadDWord reads a uint32 from reader, and stores it in value
func ReadDWord(reader io.Reader, value *uint32) {
	binary.Read(reader, binary.LittleEndian, value)
}

// ReadFloat32 reads a float32 from reader, and stores it in value
func ReadFloat32(reader io.Reader, value *float32) {
	binary.Read(reader, binary.LittleEndian, value)
}

// ReadBytes reads length bytes from reader, and stores it in value
func ReadBytes(reader io.Reader, value *[]byte, length uint32) {
	*value = make([]byte, length)
	binary.Read(reader, binary.LittleEndian, value)
}

// ReadString reads length bytes from reader, converts them to a string, and stores it in value
func ReadString(reader io.Reader, value *string, length uint32) {
	bytes := make([]byte, length)
	binary.Read(reader, binary.LittleEndian, &bytes)
	*value = string(bytes)
}
