package main

import (
	"fmt"
	"strings"
	"unicode"
)

func SplitTextIntoChunks(text string, maxChunkSize int) []string {
	var chunks []string
	var currentChunk strings.Builder
	sentences := splitIntoSentences(text)

	for _, sentence := range sentences {
		// If adding this sentence would exceed maxChunkSize and we already have content
		if currentChunk.Len() > 0 && currentChunk.Len()+len(sentence)+1 > maxChunkSize {
			// Add current chunk and start a new one
			chunks = append(chunks, strings.TrimSpace(currentChunk.String()))
			currentChunk.Reset()
		}

		// If the sentence itself is longer than maxChunkSize, it becomes its own chunk
		if len(sentence) > maxChunkSize {
			if currentChunk.Len() > 0 {
				chunks = append(chunks, strings.TrimSpace(currentChunk.String()))
				currentChunk.Reset()
			}
			chunks = append(chunks, strings.TrimSpace(sentence))
			fmt.Println("Warning: Sentence longer than maxChunkSize")
		} else {
			// Add the sentence to the current chunk
			currentChunk.WriteString(sentence)
			currentChunk.WriteString(" ")
		}
	}

	// Add the last chunk if it's not empty
	if currentChunk.Len() > 0 {
		chunks = append(chunks, strings.TrimSpace(currentChunk.String()))
	}

	return chunks
}

func splitIntoSentences(text string) []string {
	var sentences []string
	var currentSentence strings.Builder
	inQuote := false

	for i, r := range text {
		currentSentence.WriteRune(r)

		if r == '"' {
			inQuote = !inQuote
		}

		if !inQuote && (r == '.' || r == '!' || r == '?') {
			// Check if it's the end of the text or the next character is a space
			if i == len(text)-1 || unicode.IsSpace(rune(text[i+1])) {
				sentences = append(sentences, strings.TrimSpace(currentSentence.String()))
				currentSentence.Reset()
			}
		}
	}

	// Add any remaining text as a sentence
	if currentSentence.Len() > 0 {
		sentences = append(sentences, strings.TrimSpace(currentSentence.String()))
	}

	return sentences
}
