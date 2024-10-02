from pow_consumer import PoWConsumer, Challenge
import time

def main():
    # PoW server URL
    server_url = "http://localhost:8080"

    # Initialize the consumer
    consumer = PoWConsumer(server_url)

    # Define the challenge parameters
    algorithm = "sha256"  # Can be 'sha256', 'sha1' or 'modexp'
    difficulty = 5        # Adjust as needed

    try:
        # Step 1: Get the challenge
        print("Requesting PoW challenge...")
        challenge: Challenge = consumer.fetch_challenge(algo=algorithm, difficulty=difficulty)
        print(f"Challenge received: ID={challenge.challenge_id}, Algorithm={challenge.algorithm}, Difficulty={challenge.difficulty}, Data={challenge.data}")
        start = time.time()
        
        # Step 2: Solve the challenge
        print("Solving the challenge...")
        solution = consumer.solve_challenge(challenge)
        if solution is None:
            print("Could not solve the challenge.")
            return
        print(f"Solution found: {solution}")

        # Step 3: Submit the solution for verification
        print("Submitting solution for verification...")
        valid = consumer.submit_solution(challenge.challenge_id, solution)
        if valid:
            print("Valid solution! Access granted.")
            end = time.time()
            print(f"Execution time: {end - start}")
        else:
            print("Invalid solution! Access denied.")

    except Exception as e:
        print(f"An error occurred: {e}")

if __name__ == "__main__":
    main()
