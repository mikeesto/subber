package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

const (
	WhisperPath = "/Users/michael/Code/whisper.cpp/main"
	ModelPath   = "/Users/michael/Code/whisper.cpp/models/ggml-small.en.bin"
)

func main() {
	var formatFlag bool
	flag.BoolVar(&formatFlag, "format", false, "Optional: Trigger formatting the transcript into paragraphs using an LLM.")
	flag.Parse()

	if len(flag.Args()) < 1 {
		fmt.Println("Please provide a file name.")
		os.Exit(1)
	}

	fileName := flag.Args()[0]

	// If the file is not a .wav, convert it to .wav
	if filepath.Ext(fileName) != ".wav" {
		fmt.Println("Converting to .wav file...")
		cmd := exec.Command("ffmpeg", "-i", fileName, "-ar", "16000", fileName[:len(fileName)-4]+".wav")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err := cmd.Run()

		if err != nil {
			fmt.Printf("Error converting file to .wav: %v\n", err)
			return
		}

		fileName = fileName[:len(fileName)-4] + ".wav"
	}

	// Run whisper on the .wav file
	fmt.Println("Running whisper...")
	cmd := exec.Command(WhisperPath, "-m", ModelPath, "-f", fileName, "-otxt", "-of", "transcript")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()

	if err != nil {
		fmt.Printf("Error running whisper: %v\n", err)
		return
	}

	transcriptBytes, err := os.ReadFile("transcript.txt")
	if err != nil {
		fmt.Printf("Error reading transcript file: %v\n", err)
		return
	}
	transcript := strings.ReplaceAll(string(transcriptBytes), "\n", " ")

	// Delete the .wav file
	err = os.Remove(fileName)
	if err != nil {
		fmt.Printf("Error deleting wav file: %v\n", err)
		return
	}

	// If the format flag is set, format the transcript with an LLM
	if formatFlag {
		// 1 token = ~4 characters. Allow for 3000 tokens per chunk
		chunks := SplitTextIntoChunks(transcript, 12000)

		file, err := os.OpenFile("formatted_transcript.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			fmt.Printf("Error opening formatted transcript file: %v\n", err)
			return
		}
		defer file.Close()

		for _, chunk := range chunks {
			payload := fmt.Sprintf(`{"model":"llama3","system":"Split this text into paragraphs. Reply with the paragraphs only. Do not change the text at all.","prompt":%q,"stream":false,"options":{"num_ctx":8000}}`, chunk)
			req, _ := http.NewRequest("POST", "http://localhost:11434/api/generate", strings.NewReader(payload))
			req.Header.Set("Content-Type", "application/json")

			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				fmt.Printf("Error making request: %v\n", err)
				return
			}
			defer resp.Body.Close()

			bodyBytes, err := io.ReadAll(resp.Body)
			if err != nil {
				fmt.Printf("Error reading response body: %v\n", err)
				return
			}
			bodyString := string(bodyBytes)

			if resp.StatusCode != http.StatusOK {
				fmt.Println("Response Body:", bodyString)
				fmt.Printf("Error: Ollama API unreachable or returned error. Status Code: %d\n", resp.StatusCode)
				return
			}

			var result map[string]interface{}
			if err := json.Unmarshal(bodyBytes, &result); err != nil {
				fmt.Printf("Error decoding JSON response: %v\n", err)
				return
			}

			if response, ok := result["response"].(string); ok {
				fmt.Println(response)

				// Append the formatted transcript to the file
				_, err = file.WriteString(response + "\n\n")
				if err != nil {
					fmt.Printf("Error writing to formatted transcript file: %v\n", err)
					return
				}
			}
		}
	}
}
