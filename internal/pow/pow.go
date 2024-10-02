package pow

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"time"
)

// PoWSystem represents the Proof of Work system
type PoWSystem struct {
	Store ChallengeStore
}

// NewPoWSystem creates a new PoWSystem
func NewPoWSystem(store ChallengeStore) *PoWSystem {
	return &PoWSystem{Store: store}
}

// SupportedAlgorithms maps algorithm names to their implementations
var SupportedAlgorithms = map[string]PoWAlgorithm{
	"sha256": SHA256Algorithm{},
}

// GenerateChallenge generates and stores a new challenge
func (p *PoWSystem) GenerateChallenge(algo string, difficulty int) (*Challenge, error) {
	algorithm, ok := SupportedAlgorithms[algo]
	if !ok {
		return nil, errors.New("unsupported algorithm")
	}

	challengeData, err := algorithm.GenerateChallengeData(difficulty)
	if err != nil {
		return nil, err
	}

	challengeID, err := generateUniqueID()
	if err != nil {
		return nil, err
	}

	challenge := &Challenge{
		ID:         challengeID,
		Algorithm:  algo,
		Difficulty: difficulty,
		Data:       challengeData,
		CreatedAt:  time.Now(),
	}

	err = p.Store.SaveChallenge(challenge)
	if err != nil {
		return nil, err
	}

	return challenge, nil
}

// VerifySolution verifies the provided solution for a challenge
func (p *PoWSystem) VerifySolution(challengeID string, solution string) (bool, error) {
	challenge, err := p.Store.GetChallenge(challengeID)
	if err != nil {
		return false, err
	}

	algorithm, ok := SupportedAlgorithms[challenge.Algorithm]
	if !ok {
		return false, errors.New("unsupported algorithm")
	}

	isValid := algorithm.VerifySolution(challenge.Data, solution, challenge.Difficulty)
	if isValid {
		// Optionally, delete the challenge after successful verification
		p.Store.DeleteChallenge(challengeID)
	}

	return isValid, nil
}

// Helper function to generate a unique ID
func generateUniqueID() (string, error) {
	bytes := make([]byte, 16)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}
