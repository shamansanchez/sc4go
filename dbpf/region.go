package dbpf

import (
	"bytes"

	"github.com/shamansanchez/sc4go/binaryreader"
)

// RegionView contains the region view subfile
// see http://www.wiki.sc4devotion.com/index.php?title=Region_View_Subfiles
// TGI: 0xCA027EDB, 0xCA027EE1, 0x00000000
//
// INCOMPLETE: no occupant group or neighbor connection data yet
type RegionView struct {
	MajorVersion          uint16
	MinorVersion          uint16
	TileX                 uint32
	TileY                 uint32
	SizeX                 uint32
	SizeY                 uint32
	ResidentialPopulation uint32
	CommercialPopulation  uint32
	IndustrialPopulation  uint32
	Unknown               float32
	Rating                uint8
	Stars                 uint8
	Tutorial              uint8
	GUID                  uint32
	Unknown1              uint32
	Unknown2              uint32
	Unknown3              uint32
	Unknown4              uint32
	Unknown5              uint32
	God                   uint8
	NameLength            uint32
	Name                  string
	FormerNameLength      uint32
	FormerName            string
	MayorNameLength       uint32
	MayorName             string
	DescriptionLength     uint32
	Description           string
	JonasLength           uint32
	Jonas                 string
}

// ReadRegion populates a RegionView struct from the raw file
func ReadRegion(region []byte) (out RegionView) {
	out = RegionView{}

	reader := bytes.NewReader(region)
	// binary.Read(reader, binary.LittleEndian, &out)

	binaryreader.ReadWord(reader, &out.MajorVersion)
	binaryreader.ReadWord(reader, &out.MinorVersion)

	binaryreader.ReadDWord(reader, &out.TileX)
	binaryreader.ReadDWord(reader, &out.TileY)

	binaryreader.ReadDWord(reader, &out.SizeX)
	binaryreader.ReadDWord(reader, &out.SizeY)

	binaryreader.ReadDWord(reader, &out.ResidentialPopulation)
	binaryreader.ReadDWord(reader, &out.CommercialPopulation)
	binaryreader.ReadDWord(reader, &out.IndustrialPopulation)

	if out.MinorVersion > 9 {
		binaryreader.ReadFloat32(reader, &out.Unknown)
	} else {
		out.Unknown = -1
	}

	if out.MinorVersion > 10 {
		binaryreader.ReadByte(reader, &out.Rating)
	} else {
		out.Rating = 255
	}

	binaryreader.ReadByte(reader, &out.Stars)
	binaryreader.ReadByte(reader, &out.Tutorial)

	binaryreader.ReadDWord(reader, &out.GUID)

	binaryreader.ReadDWord(reader, &out.Unknown1)
	binaryreader.ReadDWord(reader, &out.Unknown2)
	binaryreader.ReadDWord(reader, &out.Unknown3)
	binaryreader.ReadDWord(reader, &out.Unknown4)
	binaryreader.ReadDWord(reader, &out.Unknown5)

	binaryreader.ReadByte(reader, &out.God)

	binaryreader.ReadDWord(reader, &out.NameLength)
	binaryreader.ReadString(reader, &out.Name, out.NameLength)

	binaryreader.ReadDWord(reader, &out.FormerNameLength)
	binaryreader.ReadString(reader, &out.FormerName, out.FormerNameLength)

	binaryreader.ReadDWord(reader, &out.MayorNameLength)
	binaryreader.ReadString(reader, &out.MayorName, out.MayorNameLength)

	binaryreader.ReadDWord(reader, &out.DescriptionLength)
	binaryreader.ReadString(reader, &out.Description, out.DescriptionLength)

	binaryreader.ReadDWord(reader, &out.JonasLength)
	binaryreader.ReadString(reader, &out.Jonas, out.JonasLength)

	return
}
