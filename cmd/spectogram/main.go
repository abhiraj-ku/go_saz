package main

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/abhiraj-ku/go_saz.git/pkg/audio"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: spectogram <input.wav>")
		os.Exit(1)
	}

	inputFile := os.Args[1]
	outputImage := "specto.png"

	gen, err := audio.NewSpectogramGenerator(inputFile)
	if err != nil {
		slog.Error("Failed to initialize spectogram generator", "error", err)
		os.Exit(1)
	}

	spectogram := gen.ComputeSpectogram()
	if err := gen.SaveSpectrogramImage(spectogram, outputImage); err != nil {
		slog.Error("Failed to save spectrogram image", "error", err)
		os.Exit(1)
	}
	slog.Info("Spectogram saved", "file", outputImage)
}
