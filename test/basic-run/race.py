import subprocess
import re
import csv
import time

# Configurations
NUM_EXECUTIONS = 5  # Number of executions for each program
LOG_PY = "log_python.txt"
LOG_GO = "log_go.txt"
DATA_CSV = "race_data.csv"

def execute_command(command, log_file):
    """Executes a command in the shell and saves the output to the log file."""
    with open(log_file, "w") as log:
        process = subprocess.run(command, shell=True, stdout=log, stderr=subprocess.PIPE, text=True)
        if process.returncode != 0:
            print(f"Error executing {command}: {process.stderr}")

def extract_time(log_file):
    """Extracts the execution time from the log file."""
    time = None
    with open(log_file, "r") as file:
        for line in file:
            if line.startswith("Execution time: "):
                match = re.search(r"Execution time:\s*([\d.]+)", line)
                if match:
                    time = float(match.group(1))
                    break
    return time

def main():
    data = {"Python": [], "Go": []}

    print("Starting race...")

    for i in range(1, NUM_EXECUTIONS + 1):
        print(f"Execution {i} of {NUM_EXECUTIONS} for Python...")
        execute_command("python3 main.py", LOG_PY)
        time_py = extract_time(LOG_PY)
        if time_py is not None:
            data["Python"].append(time_py)
            print(f"Python time: {time_py} seconds")
        else:
            print("Execution time for Python not found.")

        print(f"Execution {i} of {NUM_EXECUTIONS} for Go...")
        execute_command("go run main.go", LOG_GO)
        time_go = extract_time(LOG_GO)
        if time_go is not None:
            data["Go"].append(time_go)
            print(f"Go time: {time_go} seconds")
        else:
            print("Execution time for Go not found.")

        # Optional: wait a bit between executions
        time.sleep(1)

    # Save the data to a CSV file
    with open(DATA_CSV, "w", newline='') as csvfile:
        fieldnames = ['Execution', 'Python (s)', 'Go (s)']
        writer = csv.DictWriter(csvfile, fieldnames=fieldnames)
        writer.writeheader()
        for i in range(NUM_EXECUTIONS):
            writer.writerow({
                'Execution': i + 1,
                'Python (s)': data["Python"][i],
                'Go (s)': data["Go"][i]
            })

    print(f"Race completed. Data saved in {DATA_CSV}.")

if __name__ == "__main__":
    main()
