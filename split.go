package main

import (
	"math"
	"regexp"
	"strings"
)

type Sentence struct {
	Text   string
	Vector map[string]float64
}

func SplitTranscriptIntoParagraphs(transcript string, minSentences, maxSentences int, similarityThreshold float64) []string {
	sentences := splitIntoSentences(transcript)
	vectors := computeTFIDF(sentences)

	var paragraphs []string
	var currentParagraph []string
	var currentVector map[string]float64

	for i, sentence := range sentences {
		if len(currentParagraph) == 0 {
			currentParagraph = append(currentParagraph, sentence)
			currentVector = vectors[i]
		} else {
			similarity := cosineSimilarity(currentVector, vectors[i])

			if (similarity > similarityThreshold && len(currentParagraph) < maxSentences) || len(currentParagraph) < minSentences {
				currentParagraph = append(currentParagraph, sentence)
				currentVector = addVectors(currentVector, vectors[i])
			} else {
				paragraphs = append(paragraphs, formatParagraph(currentParagraph))
				currentParagraph = []string{sentence}
				currentVector = vectors[i]
			}
		}
	}

	if len(currentParagraph) > 0 {
		paragraphs = append(paragraphs, formatParagraph(currentParagraph))
	}

	return paragraphs
}

func formatParagraph(sentences []string) string {
	for i, sentence := range sentences {
		sentence = strings.TrimSpace(sentence)
		if !strings.HasSuffix(sentence, ".") && !strings.HasSuffix(sentence, "!") && !strings.HasSuffix(sentence, "?") {
			sentence += "."
		}
		sentences[i] = sentence
	}
	return strings.Join(sentences, " ")
}

func splitIntoSentences(text string) []string {
	re := regexp.MustCompile(`([.!?]+)\s+`)
	sentences := re.Split(text, -1)

	var result []string
	for i := 0; i < len(sentences); i++ {
		sentence := strings.TrimSpace(sentences[i])
		if sentence != "" {
			if i+1 < len(sentences) {
				sentence += re.FindString(text[len(strings.Join(sentences[:i+1], ""))-1:])
			}
			result = append(result, sentence)
		}
	}
	return result
}

func computeTFIDF(sentences []string) []map[string]float64 {
	docFreq := make(map[string]int)
	vectors := make([]map[string]float64, len(sentences))

	// First pass to calculate term frequency and document frequency
	for i, sentence := range sentences {
		words := strings.Fields(strings.ToLower(sentence))
		uniqueWords := make(map[string]bool)
		vectors[i] = make(map[string]float64)

		for _, word := range words {
			vectors[i][word]++
			uniqueWords[word] = true
		}

		for word := range uniqueWords {
			docFreq[word]++
		}
	}

	numDocs := float64(len(sentences))

	// Second pass to calculate TF-IDF
	for i := range vectors {
		totalTerms := float64(len(strings.Fields(sentences[i])))
		for word, count := range vectors[i] {
			tf := count / totalTerms
			idf := math.Log10(numDocs / float64(docFreq[word]))
			vectors[i][word] = tf * idf
		}
	}

	return vectors
}

func cosineSimilarity(v1, v2 map[string]float64) float64 {
	dotProduct := 0.0
	magnitude1 := 0.0
	magnitude2 := 0.0

	for word, val1 := range v1 {
		if val2, ok := v2[word]; ok {
			dotProduct += val1 * val2
		}
		magnitude1 += val1 * val1
	}

	for _, val2 := range v2 {
		magnitude2 += val2 * val2
	}

	return dotProduct / (math.Sqrt(magnitude1) * math.Sqrt(magnitude2))
}

func addVectors(v1, v2 map[string]float64) map[string]float64 {
	result := make(map[string]float64)
	for word, val := range v1 {
		result[word] = val
	}
	for word, val := range v2 {
		result[word] += val
	}
	return result
}
