# Subber

This is a Golang CLI tool for generating transcripts of video files. It uses ffmpeg, whisper and ollama.

## Usage

```bash
go run *.go [--format] <file>
```

### Options

- `--format`: Optional flag to trigger formatting the transcript into paragraphs using an LLM

### Example

```bash
go run *.go --format video.mp4
```

This command will:

1. Convert `video.mp4` to a WAV
2. Transcribe the audio using Whisper
3. Format the transcript into paragraphs using an LLM
4. Save the formatted transcript to `formatted_transcript.txt`

## Installation

I'm still working on the best way to distribute this. For now, you can clone the repository, ensure you have the required dependencies installed, and update the `WhisperPath` and `ModelPath` constants in the code to match your system.

To build the CLI:

```bash
go build -o subber
```

Then move the binary to somewhere on your PATH (e.g. `/usr/local/bin`).
