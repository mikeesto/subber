package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
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
	transcript = regexp.MustCompile(`\s+`).ReplaceAllString(transcript, " ")
	transcript = strings.TrimSpace(transcript)

	// Delete the .wav file
	err = os.Remove(fileName)
	if err != nil {
		fmt.Printf("Error deleting wav file: %v\n", err)
		return
	}

	// If the format flag is set, format the transcript into paragraphs
	if formatFlag {
		paragraphs := SplitTranscriptIntoParagraphs(transcript, 3, 10, 0.3)

		file, err := os.OpenFile("formatted_transcript.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			fmt.Printf("Error opening formatted transcript file: %v\n", err)
			return
		}
		defer file.Close()

		for _, paragraph := range paragraphs {
			_, err = file.WriteString(paragraph + "\n\n")
			if err != nil {
				fmt.Printf("Error writing to formatted transcript file: %v\n", err)
				return
			}
		}
	}
}
