// Package dbpf reads a DBPF file. (Specifically a SimCity 4 savegame)
package dbpf

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"io/ioutil"
)

// DBPF is the top level struct for a DBPF file
type DBPF struct {
	Header       Header
	Directory    map[string]DirectoryEntry
	Index        map[string]IndexEntry
	FileContents map[string][]byte
}

// Header DBPF File Header
type Header struct {
	Magic         [4]byte
	Major         uint32
	Minor         uint32
	UserMajor     uint32
	UserMinor     uint32
	Flags         [4]byte
	Created       uint32
	Modified      uint32
	IndexMajor    uint32
	IndexCount    uint32
	IndexLocation uint32
	IndexSize     uint32
	HoleCount     uint32
	HoleLocation  uint32
	HoleSize      uint32
	IndexMinor    uint32
}

// IndexEntry is a single entry in the file index
type IndexEntry struct {
	TypeID       uint32
	GroupID      uint32
	InstanceID   uint32
	FileLocation uint32
	FileSize     uint32
}

// DirectoryEntry DBDF directory entry
type DirectoryEntry struct {
	TypeID     uint32
	GroupID    uint32
	InstanceID uint32
	FileSize   uint32
}

func getTGIString(typeID uint32, groupID uint32, instanceID uint32) (tgi string) {
	tgi = fmt.Sprintf("%X.%X.%X", typeID, groupID, instanceID)
	return
}

// GetFileByTGI gets file contents for a given TGI
func GetFileByTGI(dbpf *DBPF, typeID uint32, groupID uint32, instanceID uint32) (contents []byte, ok bool) {
	tgi := getTGIString(typeID, groupID, instanceID)
	contents, ok = dbpf.FileContents[tgi]
	return
}

// ReadDBPF reads a dbpf file
func ReadDBPF(path string) (city DBPF) {
	city = DBPF{}

	rawBytes, err := ioutil.ReadFile(path)

	if err != nil {
		panic(err)
	}

	r := bytes.NewReader(rawBytes)
	city.Header = Header{}
	city.Index = make(map[string]IndexEntry, 0)
	city.FileContents = make(map[string][]byte, 0)
	city.Directory = make(map[string]DirectoryEntry, 0)

	binary.Read(r, binary.LittleEndian, &city.Header)

	r.Seek(int64(city.Header.IndexLocation), io.SeekStart)
	for index := 0; index < int(city.Header.IndexCount); index++ {
		i := IndexEntry{}
		binary.Read(r, binary.LittleEndian, &i)

		tgi := getTGIString(i.TypeID, i.GroupID, i.InstanceID)
		city.Index[tgi] = i
		city.FileContents[tgi] = rawBytes[i.FileLocation : i.FileLocation+i.FileSize]
	}

	if dirBytes, ok := GetFileByTGI(&city, 0xE86B1EEF, 0xE86B1EEF, 0x286B1F03); ok {
		dirReader := bytes.NewReader(dirBytes)
		dir := DirectoryEntry{}

		for {
			err := binary.Read(dirReader, binary.LittleEndian, &dir)
			tgi := getTGIString(dir.TypeID, dir.GroupID, dir.InstanceID)

			if err == io.EOF {
				break
			}

			city.Directory[tgi] = dir
		}
	}

	return
}
