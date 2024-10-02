package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"
)

const serverURL = "http://localhost:8080"
const difficulty = "5"

type Challenge struct {
	ChallengeID string `json:"challenge_id"`
	Algorithm   string `json:"algorithm"`
	Difficulty  string `json:"difficulty"`
	Data        string `json:"data"`
}

type SolutionRequest struct {
	ChallengeID string `json:"challenge_id"`
	Solution    string `json:"solution"`
}

type SolutionResponse struct {
	Valid bool `json:"valid"`
}

func main() {
	// Step 1: Fetch the challenge from the server
	challenge, err := fetchChallenge()
	if err != nil {
		fmt.Printf("Error fetching challenge: %v\n", err)
		return
	}

	// Step 2: Find the solution using multiple workers
	solution := findSolution(challenge.Data, challenge.Difficulty)
	if solution == "" {
		fmt.Println("No solution found.")
		return
	}

	// Step 3: Submit the solution to the server for verification
	valid, err := submitSolution(challenge.ChallengeID, solution)
	if err != nil {
		fmt.Printf("Error submitting solution: %v\n", err)
		return
	}

	if valid {
		fmt.Printf("Solution is valid! Nonce: %s\n", solution)
	} else {
		fmt.Println("Solution is invalid.")
	}
}

// fetchChallenge requests a challenge from the server
func fetchChallenge() (*Challenge, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/generate_challenge", serverURL), nil)
	if err != nil {
		return nil, err
	}

	q := req.URL.Query()
	q.Add("difficulty", difficulty)
	q.Add("algo", "sha256")

	req.URL.RawQuery = q.Encode()

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch challenge: %s", resp.Status)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	fmt.Println(string(body))

	var challenge Challenge
	if err := json.Unmarshal(body, &challenge); err != nil {
		return nil, err
	}

	return &challenge, nil
}

// submitSolution sends the solution back to the server for verification
func submitSolution(challengeID, solution string) (bool, error) {
	solutionRequest := SolutionRequest{
		ChallengeID: challengeID,
		Solution:    solution,
	}

	requestBody, err := json.Marshal(solutionRequest)
	if err != nil {
		return false, err
	}

	resp, err := http.Post(fmt.Sprintf("%s/verify_solution", serverURL), "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return false, fmt.Errorf("failed to verify solution: %s", resp.Status)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return false, err
	}

	var solutionResponse SolutionResponse
	if err := json.Unmarshal(body, &solutionResponse); err != nil {
		return false, err
	}

	return solutionResponse.Valid, nil
}

// findSolution attempts to find a solution using multiple workers
func findSolution(challengeData string, difficulty string) string {
	difficultyInt, err := strconv.Atoi(difficulty)
	if err != nil {
		fmt.Printf("Invalid difficulty: %v\n", err)
		return ""
	}

	prefix := strings.Repeat("0", difficultyInt)
	numWorkers := runtime.NumCPU() // Use all available CPU cores
	var wg sync.WaitGroup
	solutionChan := make(chan string, 1)
	stopChan := make(chan struct{})
	start := time.Now()
	fmt.Printf("Starting %d workers...\n", numWorkers)

	// Seed the random number generator
	rand := rand.New(rand.NewSource(time.Now().UnixNano()))

	// Calculate the total nonce space and partition it among workers
	var maxUint64 uint64 = ^uint64(0)
	nonceRange := maxUint64 / uint64(numWorkers)

	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		// Random starting point within the assigned range
		startNonce := uint64(rand.Int63n(int64(nonceRange))) + nonceRange*uint64(i)
		go func(workerID int, startNonce uint64) {
			defer wg.Done()
			worker(challengeData, prefix, solutionChan, stopChan, workerID, startNonce, nonceRange)
		}(i, startNonce)
	}

	// Wait for one worker to find the solution
	solution := <-solutionChan
	fmt.Printf("Tempo de execução: %v\n", time.Since(start).Seconds())
	close(stopChan) // Signal all workers to stop
	wg.Wait()       // Wait for all workers to finish
	return solution
}

// worker runs in a goroutine and tries to find a valid nonce
func worker(challengeData string, prefix string, solutionChan chan<- string, stopChan <-chan struct{}, workerID int, startNonce uint64, nonceRange uint64) {
	nonce := startNonce
	hash := sha256.New()
	for {
		select {
		case <-stopChan:
			return
		default:
			// Check if we've exceeded our assigned nonce range
			if nonce >= startNonce+nonceRange {
				return
			}

			// Prepare input
			data := fmt.Sprintf("%s%d", challengeData, nonce)
			hash.Reset()
			hash.Write([]byte(data))
			hashSum := hash.Sum(nil)
			hashStr := hex.EncodeToString(hashSum)

			if strings.HasPrefix(hashStr, prefix) {
				select {
				case solutionChan <- fmt.Sprintf("%d", nonce):
					fmt.Printf("Worker %d found the solution: %d\n", workerID, nonce)
					return
				default:
					return
				}
			}

			nonce++
		}
	}
}
