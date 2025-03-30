package main

import (
	"bufio"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/shashankraocoding/go-shankskit"
)

func main() {
	routes := map[string]http.HandlerFunc{
		"/":           index,
		"/transcribe": transcribeAudio,
	}
	shankskit.StartApp("GoApp", "8080", routes)
}

func index(w http.ResponseWriter, r *http.Request) {
	shankskit.Respond("templates/index.html", w, map[string]string{})
}

func transcribeAudio(w http.ResponseWriter, r *http.Request) {
	// Ensure the request method is POST
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// Get model size from form
	modelSize := r.FormValue("modelSize")
	if modelSize == "" {
		http.Error(w, "modelSize is required", http.StatusBadRequest)
		return
	}

	// Get uploaded file
	file, _, err := r.FormFile("audio")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer file.Close()

	wd, _ := os.Getwd()
	filepath := filepath.Join(wd, "temp.wav")

	// Create a temporary file to store the uploaded content
	tempFile, err := os.Create(filepath)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Copy file data to temp file
	if _, err := io.Copy(tempFile, file); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tempFile.Close()
	defer os.Remove(filepath) // Clean up after transcription

	// Start Python process with model size and file path arguments
	cmd := exec.Command("python310", "transcribe.py", modelSize, filepath)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Start the command
	if err := cmd.Start(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Read output from Python script
	var output string
	scanner := bufio.NewScanner(stdout)
	for scanner.Scan() {
		output += scanner.Text() + "\n"
	}

	if err := scanner.Err(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Wait for Python process to finish
	if err := cmd.Wait(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data := map[string]string{
		"text": output,
	}
	shankskit.Respond("templates/transcript.html", w, data)
}
