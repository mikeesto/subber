package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

const (
	WhisperPath = "/Users/michael/Code/whisper.cpp/main"
	ModelPath   = "/Users/michael/Code/whisper.cpp/models/ggml-medium.en.bin"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Please provide a file name.")
		os.Exit(1)
	}

	fileName := os.Args[1]

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

	// Delete the .wav file
	err = os.Remove(fileName)
	if err != nil {
		fmt.Printf("Error deleting wav file: %v\n", err)
		return
	}
}
