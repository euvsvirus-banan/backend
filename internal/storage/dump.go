package storage

import (
	"encoding/json"
	"fmt"
	"io"
)

func dump(wr io.WriteSeeker, data interface{}) error {
	if _, err := wr.Seek(0, 0); err != nil {
		return fmt.Errorf("problem rewinding file: %w", err)
	}
	b, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("problem marshaling data: %w", err)
	}
	if _, err = wr.Write(b); err != nil {
		return fmt.Errorf("problem saving data: %w", err)
	}
	return nil
}
