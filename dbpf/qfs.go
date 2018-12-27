package dbpf

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

// QFSHeader contains magic and size
type QFSHeader struct {
	Magic      [2]byte
	OutputSize [3]byte
}

// QFSDecompress decompresses QFS compressed data
// see http://www.wiki.sc4devotion.com/index.php?title=DBPF_Compression
func QFSDecompress(data []byte) (result []byte) {
	reader := bytes.NewReader(data)
	length := uint32(0)

	binary.Read(reader, binary.LittleEndian, &length)
	if len(data) != int(length) {
		fmt.Println("lengths mismatch")
	}

	header := QFSHeader{}

	binary.Read(reader, binary.LittleEndian, &header)

	outputSize := uint32(header.OutputSize[0])<<16 +
		uint32(header.OutputSize[1])<<8 +
		uint32(header.OutputSize[2])

	result = make([]byte, outputSize)

	outputPos := 0

	cc0 := 0
	cc1 := 0
	cc2 := 0
	cc3 := 0

	for {

		numPlain := 0
		numCopy := 0
		copyOffset := 0

		b, err := reader.ReadByte()

		cc0 = int(b)

		if err != nil {
			break
		}

		if cc0 >= 0xfc {
			// CC length:      1 byte
			// Num plain text: (byte0 & 0x03)
			// Num to copy:    0
			// Copy offset:    -

			numPlain = cc0 & 0x03
		} else if cc0 >= 0xe0 {
			// CC length:      1 byte
			// Num plain text: ((byte0 & 0x1F) < < 2 ) + 4
			// Num to copy:    0
			// Copy offset:    -

			numPlain = (cc0&0x1f)<<2 + 4
		} else if cc0 >= 0xc0 {
			// CC length:      4 bytes
			// Num plain text: byte0 & 0x03
			// Num to copy:    ( (byte0 & 0x0C) < < 6 )  + byte3 + 5
			// Copy offset:    ((byte0 & 0x10) < < 12 ) + (byte1 < < 8 ) + byte2 + 1

			b, _ := reader.ReadByte()
			cc1 = int(b)
			b, _ = reader.ReadByte()
			cc2 = int(b)
			b, _ = reader.ReadByte()
			cc3 = int(b)

			numPlain = cc0 & 0x03
			numCopy = (cc0&0x0c)<<6 + cc3 + 5
			copyOffset = (cc0&0x10)<<12 + (cc1 << 8) + cc2 + 1

		} else if cc0 >= 0x80 {
			// CC length:      3 bytes
			// Num plain text: ((byte1 & 0xC0) > > 6 ) & 0x03
			// Num to copy:    (byte0 & 0x3F) + 4
			// Copy offset:    ( (byte1 & 0x3F) < < 8 ) + byte2 + 1

			b, _ := reader.ReadByte()
			cc1 = int(b)
			b, _ = reader.ReadByte()
			cc2 = int(b)

			numPlain = (cc1 & 0xc0) >> 6 & 0x03
			numCopy = (cc0 & 0x3f) + 4
			copyOffset = ((cc1 & 0x3f) << 8) + cc2 + 1

		} else {
			// CC length:      2 bytes
			// Num plain text: byte0 & 0x03
			// Num to copy:    ( (byte0 & 0x1C) > > 2) + 3
			// Copy offset:    ( (byte0 & 0x60) < < 3) + byte1 + 1

			b, _ := reader.ReadByte()
			cc1 = int(b)

			numPlain = cc0 & 0x03
			numCopy = ((cc0 & 0x1c) >> 2) + 3
			copyOffset = ((cc0 & 0x60) << 3) + cc1 + 1
		}

		for num := 0; num < numPlain; num++ {
			b, _ := reader.ReadByte()
			result[outputPos] = b
			outputPos++
		}

		start := outputPos - copyOffset
		for num := 0; num < numCopy; num++ {
			result[outputPos] = result[start]
			outputPos++
			start++
		}

	}

	if outputPos != int(outputSize) {
		fmt.Printf("Expected %d, but only wrote %d bytes!", outputSize, outputPos)
	}
	return
}
