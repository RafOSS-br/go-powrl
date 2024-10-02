package pow

import (
	"time"
)

// Challenge represents a PoW challenge
type Challenge struct {
	ID         string    // Unique identifier
	Algorithm  string    // Algorithm used
	Difficulty int       // Difficulty level
	Data       string    // Random data for the challenge
	CreatedAt  time.Time // Timestamp
}
