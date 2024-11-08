package character

import (
	"testing"
)

func TestFromFile(t *testing.T) {
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
		// Add more test cases as needed
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := FromFile(tt.data)
			if err != tt.wantErr {
				t.Errorf("FromFile() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
