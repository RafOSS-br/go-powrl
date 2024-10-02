package pow

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"strings"
)

// PoWAlgorithm defines methods for challenge generation and verification
type PoWAlgorithm interface {
	GenerateChallengeData(difficulty int) (string, error)
	VerifySolution(challengeData string, solution string, difficulty int) bool
}

// SHA256Algorithm implements PoWAlgorithm using SHA-256
type SHA256Algorithm struct{}

func (a SHA256Algorithm) GenerateChallengeData(difficulty int) (string, error) {
	randomBytes := make([]byte, 32)
	_, err := rand.Read(randomBytes)
	if err != nil {
		return "", err
	}
	data := hex.EncodeToString(randomBytes)
	return data, nil
}

func (a SHA256Algorithm) VerifySolution(challengeData string, solution string, difficulty int) bool {
	hash := sha256.Sum256([]byte(challengeData + solution))
	hashStr := hex.EncodeToString(hash[:])
	prefix := strings.Repeat("0", difficulty)
	return strings.HasPrefix(hashStr, prefix)
}
