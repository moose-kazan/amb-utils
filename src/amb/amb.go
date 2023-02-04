package amb

/*
 * Copyright (C) 2023 Vadim Kalinnikov <moose@ylsoftware.com>
 * This code published under terms of LGPL-2.1
 * GIthub repo: https://github.com/moose-kazan/amb-utils
 */

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io/ioutil"
	"strings"
)

type AmbFile struct {
	fileName string
	entries  map[string][]byte
	charMap  [128]int16
}

type AmbFileInterface interface {
	LoadFile() error
	InitCharMap()
	HasEntry(name string) bool
	GetEntryRaw(name string) ([]byte, error)
	GetEntry(name string) ([]byte, error)
	BSDSum(data []byte) int
}

func New(fileName string) *AmbFile {
	var a AmbFile
	a.fileName = fileName
	// Default charmap value
	a.charMap = [128]int16{
		0x00C7, 0x00FC, 0x00E9, 0x00E2, 0x00E4, 0x00E0, 0x00E5, 0x00E7, 0x00EA, 0x00EB, 0x00E8, 0x00EF, 0x00EE, 0x00EC, 0x00C4, 0x00C5,
		0x00C9, 0x00E6, 0x00C6, 0x00F4, 0x00F6, 0x00F2, 0x00FB, 0x00F9, 0x00FF, 0x00D6, 0x00DC, 0x00A2, 0x00A3, 0x00A5, 0x20A7, 0x0192,
		0x00E1, 0x00ED, 0x00F3, 0x00FA, 0x00F1, 0x00D1, 0x00AA, 0x00BA, 0x00BF, 0x2310, 0x00AC, 0x00BD, 0x00BC, 0x00A1, 0x00AB, 0x00BB,
		0x2591, 0x2592, 0x2593, 0x2502, 0x2524, 0x2561, 0x2562, 0x2556, 0x2555, 0x2563, 0x2551, 0x2557, 0x255D, 0x255C, 0x255B, 0x2510,
		0x2514, 0x2534, 0x252C, 0x251C, 0x2500, 0x253C, 0x255E, 0x255F, 0x255A, 0x2554, 0x2569, 0x2566, 0x2560, 0x2550, 0x256C, 0x2567,
		0x2568, 0x2564, 0x2565, 0x2559, 0x2558, 0x2552, 0x2553, 0x256B, 0x256A, 0x2518, 0x250C, 0x2588, 0x2584, 0x258C, 0x2590, 0x2580,
		0x03B1, 0x00DF, 0x0393, 0x03C0, 0x03A3, 0x03C3, 0x00B5, 0x03C4, 0x03A6, 0x0398, 0x03A9, 0x03B4, 0x221E, 0x03C6, 0x03B5, 0x2229,
		0x2261, 0x00B1, 0x2265, 0x2264, 0x2320, 0x2321, 0x00F7, 0x2248, 0x00B0, 0x2219, 0x00B7, 0x221A, 0x207F, 0x00B2, 0x25A0, 0x00A0,
	}
	return &a
}

func (a *AmbFile) LoadFile() error {
	// Try to load file into memory
	data, err := ioutil.ReadFile(a.fileName)
	if err != nil {
		return err
	}

	// Check file signature
	if string(data[0:4]) != "AMB1" {
		return errors.New("Not an AMB file!")
	}

	// Load entries
	a.entries = make(map[string][]byte)
	var filesCount = int(binary.LittleEndian.Uint16(data[4:6]))
	for i := 0; i < filesCount; i++ {
		entryOffset := 6 + i*20
		// Entry name
		entryName := ""
		for j := entryOffset; j < entryOffset+12 && data[j] != 0; j++ {
			entryName += string(data[j])
		}
		entryName = strings.ToLower(entryName)

		// Entry data
		entryFileOffset := int(binary.LittleEndian.Uint32(data[entryOffset+12 : entryOffset+16]))
		entryFileLength := int(binary.LittleEndian.Uint16(data[entryOffset+16 : entryOffset+18]))
		entryCkSum := uint16(binary.LittleEndian.Uint16(data[entryOffset+18 : entryOffset+20]))
		entryData := data[entryFileOffset : entryFileOffset+entryFileLength]

		// Check BSD sum of entry
		if entryCkSum != a.BSDSum(entryData) {
			return errors.New(fmt.Sprintf("Bad checksum of entry \"%s\"!", entryName))
		}

		// Put entry into entries-list
		a.entries[entryName] = entryData

	}

	// Load charMap if exists
	if ambEntry, ok := a.entries["unicode.map"]; ok {
		for i := 0; i < 128; i++ {
			a.charMap[i] = int16(binary.LittleEndian.Uint16(ambEntry[i*2 : i*2+2]))
		}
	}

	return nil
}

func (a *AmbFile) ListNames() []string {
	var rv []string
	for k := range a.entries {
		rv = append(rv, k)
	}
	return rv
}

func (a *AmbFile) HasEntry(name string) bool {
	if _, ok := a.entries[name]; ok {
		return true
	}
	return false
}

func (a *AmbFile) GetEntryRaw(name string) ([]byte, error) {
	if entry, ok := a.entries[name]; ok {
		return entry, nil
	}
	return make([]byte, 0), errors.New("Entry not found")
}

func (a *AmbFile) GetEntry(name string) ([]byte, error) {
	if entry, ok := a.entries[name]; ok {
		var rv []rune
		for _, c := range entry {
			if c < 128 {
				rv = append(rv, rune(c))
			} else {
				rv = append(rv, rune(a.charMap[c-128]))
			}
		}
		return []byte(string(rv)), nil
	}
	return make([]byte, 0), errors.New("Entry not found")
}

// Implements BSD "sum" hash
func (a *AmbFile) BSDSum(data []byte) uint16 {
	var sum uint16
	for _, b := range data {
		sum = (sum >> 1) + (sum << 15) + uint16(b)
	}
	return sum
}
