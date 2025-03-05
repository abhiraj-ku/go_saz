package audio

import (
	"fmt"
	"image/color"
	"image/png"
	"math"
	"os"

	"image"

	"github.com/go-audio/wav"
	"github.com/mjibson/go-dsp/fft"
)

// Spectogram specifications (consts)
const (
	WindowSize    = 1024 // no of samples per fft Window
	Overlap       = 512  // 50% overlap for smooth transition
	FrequencyBins = WindowSize / 2
)

// Spectogram struct
type SpectogramGenrator struct {
	filepath   string
	data       []float64
	sampleRate int
}

// NewSpectogram initializes the Spectogram and reads wav file data

func NewSpectogramGenerator(filepath string) (*SpectogramGenrator, error) {
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

	return &SpectogramGenrator{filepath, floatData, int(decoder.SampleRate)}, nil

}

// ComputeSpectogram performs SFTF and returns 2-D matrix of (time v/s frequency)

func (s *SpectogramGenrator) ComputeSpectogram() [][]float64 {
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
			windowedData[j] = s.data[start+j] * hannWindow[j]
		}

		// compute FTT for each window
		fftResult := fft.FFTReal(windowedData)

		//compute magnittude spectrum
		magnitudes := make([]float64, FrequencyBins)
		for j := 0; j < FrequencyBins; j++ {
			magnitudes[j] = cmplxAbs(fftResult[j])
		}
		spectoGram[i] = magnitudes
	}
	return spectoGram
}

// Generate spectogram as PNG and save it
func (s *SpectogramGenrator) SaveSpectrogramImage(spectogram [][]float64, outputPath string) error {
	height := len(spectogram[0]) //FrequencyBins
	width := len(spectogram)

	img := image.NewGray(image.Rect(0, 0, width, height))

	maxVal := 0.0
	for _, row := range spectogram {
		for _, val := range row {
			if val > maxVal {
				maxVal = val
			}
		}
	}
	if maxVal == 0 {
		return fmt.Errorf("spectogram has only zero values, unable to normalize: ")
	}
	// Normalize values for visualization
	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			intensity := uint8((spectogram[x][height-y-1] / maxVal) * 255)
			img.SetGray(x, y, color.Gray{Y: intensity})
		}
	}

	// save the image to outputPath
	file, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("unable to save spectrum image: %w", err)
	}
	defer file.Close()
	return png.Encode(file, img)

}

// applyHannWindow generates a hann window function

func applyHannWindow(size int) []float64 {
	window := make([]float64, size)
	for i := range window {
		window[i] = 0.5 * (1 - math.Cos(2*math.Pi*float64(i)/float64(size-1)))
	}
	return window
}

// cmplxAbs compuytes mag of a complex number
func cmplxAbs(c complex128) float64 {
	return math.Sqrt(real(c)*real(c) + imag(c)*imag(c))
}
