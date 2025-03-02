package main

import (
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/abhiraj-ku/go_saz.git/pkg/audio"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: record <output.wav>")
		os.Exit(1)
	}

	outputFile := os.Args[1]
	duration := 5 * time.Second

	// initialize the recorder
	rec, err := audio.NewRecorder(outputFile)
	if err != nil {
		slog.Error("Failed to initialize recorder", "error", err)
		os.Exit(1)
	}

	defer rec.Close()

	// start recording
	if err := rec.Start(duration); err != nil {
		slog.Error("failed to start recording", "error", err)
		os.Exit(1)
	}
	slog.Info("Audio recording saved", "file", outputFile)

}
