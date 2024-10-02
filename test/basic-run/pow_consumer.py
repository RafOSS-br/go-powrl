import requests
import hashlib
import random
import string
from dataclasses import dataclass
from typing import Optional, Dict
from sympy import mod_inverse

@dataclass
class Challenge:
    challenge_id: str
    algorithm: str
    difficulty: int
    data: str

class PoWConsumer:
    def __init__(self, server_url: str):
        """
        Initializes the PoW consumer.

        :param server_url: Base URL of the PoW server, e.g., 'http://localhost:8080'
        """
        self.server_url = server_url.rstrip('/')

    def fetch_challenge(self, algo: Optional[str] = None, difficulty: Optional[int] = None) -> Challenge:
        """
        Requests a PoW challenge from the server.

        :param algo: Desired algorithm ('sha256', 'sha1', 'modexp'). If None, uses the server default.
        :param difficulty: Difficulty level. If None, uses the server default.
        :return: Challenge object containing the challenge details.
        :raises: Exception if the request fails or the server returns an error.
        """
        params = {}
        if algo:
            params['algo'] = algo
        if difficulty:
            params['difficulty'] = str(difficulty)

        url = f"{self.server_url}/generate_challenge"
        response = requests.get(url, params=params)

        if response.status_code != 200:
            raise Exception(f"Failed to fetch challenge: {response.text}")

        data = response.json()
        return Challenge(
            challenge_id=data['challenge_id'],
            algorithm=data['algorithm'],
            difficulty=int(data['difficulty']),
            data=data['data']
        )

    def solve_challenge(self, challenge: Challenge) -> Optional[str]:
        """
        Solves the PoW challenge according to the specified algorithm.

        :param challenge: Challenge object to be solved.
        :return: Solution as a string or None if it cannot be solved.
        """
        if challenge.algorithm == 'sha256':
            return self._solve_sha256(challenge.data, challenge.difficulty)
        elif challenge.algorithm == 'sha1':
            return self._solve_sha1(challenge.data, challenge.difficulty)
        elif challenge.algorithm == 'modexp':
            return self._solve_modexp(challenge.data)
        else:
            print(f"Unknown algorithm: {challenge.algorithm}")
            return None

    def submit_solution(self, challenge_id: str, solution: str) -> bool:
        """
        Submits the challenge solution to the server for verification.

        :param challenge_id: Unique ID of the challenge.
        :param solution: Found solution.
        :return: True if the solution is valid, False otherwise.
        :raises: Exception if the request fails or the server returns an error.
        """
        url = f"{self.server_url}/verify_solution"
        payload = {
            "challenge_id": challenge_id,
            "solution": solution
        }
        headers = {'Content-Type': 'application/json'}

        response = requests.post(url, json=payload, headers=headers)

        if response.status_code != 200:
            raise Exception(f"Failed to submit solution: {response.text}")

        data = response.json()
        return data.get('valid', False)

    def _solve_sha256(self, data: str, difficulty: int) -> Optional[str]:
        """
        Solves a SHA-256 Hashcash challenge.

        :param data: Challenge data.
        :param difficulty: Number of leading zeros required in the hash.
        :return: Valid nonce as a string.
        """
        prefix = '0' * difficulty
        nonce = 0
        print(f"Solving SHA-256 with difficulty {difficulty}...")
        while True:
            solution = f"{nonce}"
            hash_input = f"{data}{solution}".encode('utf-8')
            hash_result = hashlib.sha256(hash_input).hexdigest()
            if hash_result.startswith(prefix):
                print(f"Solution found: {solution} (Hash: {hash_result})")
                return solution
            nonce += 1

    def _solve_sha1(self, data: str, difficulty: int) -> Optional[str]:
        """
        Solves a SHA-1 Hashcash challenge.

        :param data: Challenge data.
        :param difficulty: Number of leading zeros required in the hash.
        :return: Valid nonce as a string.
        """
        prefix = '0' * difficulty
        nonce = 0
        print(f"Solving SHA-1 with difficulty {difficulty}...")
        while True:
            solution = f"{nonce}"
            hash_input = f"{data}{solution}".encode('utf-8')
            hash_result = hashlib.sha1(hash_input).hexdigest()
            if hash_result.startswith(prefix):
                print(f"Solution found: {solution} (Hash: {hash_result})")
                return solution
            nonce += 1
            if nonce % 100000 == 0:
                print(f"Testing nonce: {nonce}")

    def _solve_modexp(self, data: str) -> Optional[str]:
        """
        Solves a modular exponentiation challenge.

        :param data: Challenge data in the format 'base|exponent|modulus' in hexadecimal.
        :return: Result of the modular exponentiation as a hexadecimal string.
        """
        try:
            parts = data.split('|')
            if len(parts) != 3:
                print("Invalid data format for modexp.")
                return None

            base_hex, exponent_hex, modulus_hex = parts
            base = int(base_hex, 16)
            exponent = int(exponent_hex, 16)
            modulus = int(modulus_hex, 16)

            print(f"Solving Modular Exponentiation: ({base} ^ {exponent}) % {modulus}")
            result = pow(base, exponent, modulus)
            result_hex = hex(result)[2:]
            print(f"Solution found: {result_hex}")
            return result_hex
        except Exception as e:
            print(f"Error solving modexp: {e}")
            return None
