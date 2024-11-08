// Package character provides functionality for parsing character cards stored in PNG format
package character

import (
	"bytes"
	"encoding/base64"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
)

// Common errors
var (
	ErrNotPNG = errors.New("not a PNG file")
	ErrNoCharacterData = errors.New("no character data found in PNG")
	ErrIncompletePNGChunk = errors.New("incomplete PNG chunk")
	ErrIncompletePNGChunkData = errors.New("incomplete PNG chunk data")
	"bytes"
	"encoding/base64"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
)

// CharacterMetadata contains all the metadata fields for a character
type CharacterMetadata struct {
	AlternateGreetings []interface{} `json:"alternate_greetings"`
	Avatar             string        `json:"avatar"`
	CharacterBook      interface{}   `json:"character_book"`
	CharacterVersion   string        `json:"character_version"`
	Chat              string        `json:"chat"`
	CreateDate        string        `json:"create_date"`
	Creator           string        `json:"creator"`
	CreatorNotes      string        `json:"creator_notes"`
	Description       string        `json:"description"`
	Extensions        struct {
		Chub struct {
			Expressions      interface{} `json:"expressions"`
			FullPath        string      `json:"full_path"`
			ID              int         `json:"id"`
			RelatedLorebooks []interface{} `json:"related_lorebooks"`
		} `json:"chub"`
		Fav           bool   `json:"fav"`
		Talkativeness string `json:"talkativeness"`
	} `json:"extensions"`
	FirstMes                 string   `json:"first_mes"`
	MesExample              string   `json:"mes_example"`
	Name                    string   `json:"name"`
	Personality             string   `json:"personality"`
	PostHistoryInstructions string   `json:"post_history_instructions"`
	Scenario                string   `json:"scenario"`
	SystemPrompt            string   `json:"system_prompt"`
	Tags                    []string `json:"tags"`
	CharGreeting           string   `json:"char_greeting"`
	ExampleDialogue        string   `json:"example_dialogue"`
	WorldScenario          string   `json:"world_scenario"`
	CharPersona            string   `json:"char_persona"`
	CharName              string   `json:"char_name"`
}

type Character struct {
	Metadata       CharacterMetadata
	fallbackAvatar string
}

func (c *Character) Avatar() string {
	if c.Metadata.Avatar != "" && c.Metadata.Avatar != "none" {
		return c.Metadata.Avatar
	}
	return c.fallbackAvatar
}

func (c *Character) Description() string {
	return c.Metadata.Description
}

func (c *Character) Name() string {
	if c.Metadata.Name != "" {
		return c.Metadata.Name
	}
	return c.Metadata.CharName
}

type pngChunk struct {
	Length uint32
	Type   [4]byte
	Data   []byte
	CRC    uint32
}

func FromFile(data []byte) (*Character, error) {
	// Check PNG signature
	if len(data) < 8 || !bytes.Equal(data[:8], []byte{137, 80, 78, 71, 13, 10, 26, 10}) {
		return nil, errors.New("not a PNG file")
	}

	// Parse PNG chunks
	chunks, err := extractPNGChunks(data)
	if err != nil {
		return nil, err
	}

	// Find tEXt chunk with character data
	var charData []byte
	var imageChunks []pngChunk
	
	for _, chunk := range chunks {
		if string(chunk.Type[:]) == "tEXt" {
			keyword, text, found := bytes.Cut(chunk.Data, []byte{0})
			if found && string(keyword) == "chara" {
				charData = text
			}
		} else {
			imageChunks = append(imageChunks, chunk)
		}
	}

	if charData == nil {
		return nil, errors.New("no character data found in PNG")
	}

	// Decode base64 character data
	jsonData, err := base64.StdEncoding.DecodeString(string(charData))
	if err != nil {
		return nil, fmt.Errorf("failed to decode base64 character data: %w", err)
	}

	// Try parsing as V2 format first
	var v2Format struct {
		SpecVersion string           `json:"spec_version"`
		Data        CharacterMetadata `json:"data"`
	}
	if err := json.Unmarshal(jsonData, &v2Format); err == nil && v2Format.SpecVersion == "2.0" {
		// Encode remaining chunks back to PNG for avatar
		avatarData, err := encodePNGChunks(imageChunks)
		if err != nil {
			return nil, err
		}
		avatarBase64 := "data:image/png;base64," + base64.StdEncoding.EncodeToString(avatarData)
		
		return &Character{
			Metadata:       v2Format.Data,
			fallbackAvatar: avatarBase64,
		}, nil
	}

	// Try V1 format
	var metadata CharacterMetadata
	if err := json.Unmarshal(jsonData, &metadata); err != nil {
		return nil, fmt.Errorf("failed to parse character data: %w", err)
	}

	// Encode remaining chunks back to PNG for avatar
	avatarData, err := encodePNGChunks(imageChunks)
	if err != nil {
		return nil, err
	}
	avatarBase64 := "data:image/png;base64," + base64.StdEncoding.EncodeToString(avatarData)

	return &Character{
		Metadata:       metadata,
		fallbackAvatar: avatarBase64,
	}, nil
}

func extractPNGChunks(data []byte) ([]pngChunk, error) {
	var chunks []pngChunk
	pos := 8 // Skip PNG signature

	for pos < len(data) {
		if pos+12 > len(data) {
			return nil, errors.New("incomplete PNG chunk")
		}

		chunk := pngChunk{}
		chunk.Length = binary.BigEndian.Uint32(data[pos:])
		copy(chunk.Type[:], data[pos+4:pos+8])
		
		dataStart := pos + 8
		dataEnd := dataStart + int(chunk.Length)
		if dataEnd+4 > len(data) {
			return nil, errors.New("incomplete PNG chunk data")
		}
		
		chunk.Data = data[dataStart:dataEnd]
		chunk.CRC = binary.BigEndian.Uint32(data[dataEnd:dataEnd+4])
		
		chunks = append(chunks, chunk)
		pos = dataEnd + 4

		if string(chunk.Type[:]) == "IEND" {
			break
		}
	}

	return chunks, nil
}

func encodePNGChunks(chunks []pngChunk) ([]byte, error) {
	// Calculate total size
	size := 8 // PNG signature
	for _, chunk := range chunks {
		size += 12 + len(chunk.Data) // Length + Type + Data + CRC
	}

	// Create output buffer
	buf := make([]byte, 0, size)

	// Write PNG signature
	buf = append(buf, 137, 80, 78, 71, 13, 10, 26, 10)

	// Write chunks
	for _, chunk := range chunks {
		// Length
		lengthBuf := make([]byte, 4)
		binary.BigEndian.PutUint32(lengthBuf, uint32(len(chunk.Data)))
		buf = append(buf, lengthBuf...)

		// Type
		buf = append(buf, chunk.Type[:]...)

		// Data
		buf = append(buf, chunk.Data...)

		// CRC
		crcBuf := make([]byte, 4)
		binary.BigEndian.PutUint32(crcBuf, chunk.CRC)
		buf = append(buf, crcBuf...)
	}

	return buf, nil
}
