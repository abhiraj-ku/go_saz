package audio

import (
	"fmt"
	"os"

	"github.com/go-audio/wav"
)

// Spectogram specifications (consts)
const (
	WindowSize    = 1024 // no of samples per fft Window
	Overlap       = 512  // 50% overlap for smooth transition
	FrequencyBins = WindowSize / 2
)

// Spectogram struct
type Spectogram struct {
	filepath   string
	data       []float64
	sampleRate int
}

// NewSpectogram initializes the Spectogram and reads wav file data

func NewSpectogramGenerator(filepath string) (*Spectogram, error) {
	// Open the WAV file and read PCM data (among the 44 bits of tha audio data)
	file, err := os.Open(filepath)
	if err != nil {
		return nil, fmt.Errorf("failed to open desired file: %w", err)
	}
	defer file.Close()

	// Using go-audio/wav read the wav file content
	decoder := wav.NewDecoder(file)
	if !decoder.IsValidFile() {
		return nil, fmt.Errorf("provided file is not a .wav file: %w", err)
	}
	// read data
	decoder.ReadInfo()
	if decoder.SampleRate != SampleRate {
		return nil, fmt.Errorf("unspported sample rate: %w", err)
	}

	// Read the raw PCM data
	buf, err := decoder.FullPCMBuffer()
	if err != nil {
		return nil, fmt.Errorf("failed to read the PCM data: %w", err)
	}

	// convert PCM data to float64 (normalize it to [-1,1])
	floatData := make([]float64, len(buf.Data))
	for i, sample := range buf.Data {
		floatData[i] = float64(sample) / float64(1<<15)
	}

	return &Spectogram{filepath, floatData, int(decoder.SampleRate)}, nil

}

// ComputeSpectogram performs SFTF and returns 2-D matrix of (time v/s frequency)

func (s *Spectogram) ComputeSpectogram() [][]float64 {
	numFrames := (len(s.data) - WindowSize) / Overlap
	spectoGram := make([][]float64, numFrames)

	// hann window to manage spectra leakage
	hannWindow := applyHannWindow(WindowSize)

	// outer loop for each audio frame of data
	for i := 0; i < numFrames; i++ {
		start := i * Overlap
		windowedData := make([]float64, WindowSize)

		//apply hann window
		for j := 0; j < WindowSize; j++ {
			windowedData[j] = s.data[start+j] * han
		}
	}
}
