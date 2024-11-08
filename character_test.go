package character

import (
	"errors"
	"os"
	"testing"
)

func TestFromFile(t *testing.T) {
	// Read test data
	validData, err := os.ReadFile("testdata/character.png")
	if err != nil {
		t.Fatalf("Failed to read test file: %v", err)
	}

	tests := []struct {
		name    string
		data    []byte
		wantErr error
	}{
		{
			name:    "invalid PNG signature",
			data:    []byte("not a PNG"),
			wantErr: ErrNotPNG,
		},
		{
			name:    "valid character file",
			data:    validData,
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			char, err := FromFile(tt.data)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("FromFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if err == nil {
				t.Logf("Character Name: %s", char.Name())
				t.Logf("Character Description: %s", char.Description())
			}
		})
	}
}
