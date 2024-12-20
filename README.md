# Website for creating spectrograms.

This is a simple web application running on localhost with a basic frontend. It processes audio files to generate spectrograms using two separate backends: one written in Go and the other in Python.

## How to Run the Application

### 1. Run the Go Backend
1. Open a terminal.
2. Navigate to the `/back` directory:
   ```bash
   cd back
   ```
3. Start the Go backend by running:
   ```bash
   go run .
   ```

### 2. Run the Python Backend
1. Open another terminal.
2. Navigate to the `/pythonServ` directory:
   ```bash
   cd pythonServ
   ```
3. Start the Python backend by running:
   ```bash
   python run.py
   ```

## Notes
- The frontend is served by the Go backend and can be accessed via `http://localhost:8080` (or the port configured in your Go code).
- The Python backend handles audio processing and spectrogram generation.

## Dependencies

### Go Backend
- Go 1.20 or higher.

### Python Backend
- Python 3.9 or higher.
  
---

This project is designed for local use and demonstrates basic audio processing capabilities. Feedback and contributions are welcome!
