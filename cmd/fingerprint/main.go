package main

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/abhiraj-ku/go_saz.git/pkg/audio"
)

func main() {
	if len(os.Args) < 2 {
		slog.Info("Not enough arguements", "Usage: fingerprint <spectogram.png>", 1)
		os.Exit(1)
	}

	inputFile := os.Args[1]

	// Load spectrogram
	gen, err := audio.NewSpectogramGenerator(inputFile)
	if err != nil {
		slog.Error("Failed to load spectrogram", "error", err)
		os.Exit(1)
	}

	spectrogram := gen.ComputeSpectogram()

	// Extract Fingerprint
	fingerprints, err := audio.ExtractFingerprints(spectrogram)
	if err != nil {
		slog.Error("Failed to extract fingerprints", "error", err)
		os.Exit(1)
	}

	slog.Info("Fingerprinting completed", "count", len(fingerprints))
	for _, fp := range fingerprints[:10] { // print first 10 fingerprints
		fmt.Println(fp)
	}

}
