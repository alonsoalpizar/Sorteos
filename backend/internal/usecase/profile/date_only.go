package profile

import (
	"fmt"
	"strings"
	"time"
)

// DateOnly representa una fecha sin hora (YYYY-MM-DD)
type DateOnly struct {
	time.Time
}

// UnmarshalJSON parsea una fecha en formato YYYY-MM-DD
func (d *DateOnly) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), "\"")
	if s == "null" || s == "" {
		return nil
	}

	// Parsear formato YYYY-MM-DD
	t, err := time.Parse("2006-01-02", s)
	if err != nil {
		return fmt.Errorf("invalid date format, expected YYYY-MM-DD: %w", err)
	}

	d.Time = t
	return nil
}

// MarshalJSON serializa la fecha en formato YYYY-MM-DD
func (d DateOnly) MarshalJSON() ([]byte, error) {
	if d.Time.IsZero() {
		return []byte("null"), nil
	}
	return []byte(fmt.Sprintf("\"%s\"", d.Time.Format("2006-01-02"))), nil
}
