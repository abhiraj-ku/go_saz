package audio

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"log/slog"
	"math"
)

// Constants for fingerprinting
const (
	MinPeakThreshold  = 0.5 // Min rel amplitude for peak consideration
	AnchorDistance    = 5   // Distance in time bins for anchor pairing
	FrequencyDistance = 10  // Max frequency sep for peak pairing

)

// Fingerprinting struct --> represent extracted audio fingerprints
type Fingerprint struct {
	Hash      string
	TimeStamp int
}

// ExtractFingerprints --> process the spectogram to gen unique fingerprints
func ExtractFingerprints(spectogram [][]float64) ([]Fingerprint, error) {
	if len(spectogram) == 0 || len(spectogram[0]) == 0 {
		return nil, fmt.Errorf("empty spectogram data")
	}

	var fingerprints []Fingerprint
	peaks := detectPeaks(spectogram)

	// Iterate through each time bins --> generate fingerprints
	for t, peakFreqs := range peaks {
		for _, anchorFreq := range peakFreqs {
			for futureT := t + 1; futureT < len(peaks) && futureT <= t+AnchorDistance; futureT++ {
				for _, peakFreq := range peaks[futureT] {
					if math.Abs(float64(peakFreq-anchorFreq)) <= FrequencyDistance {

						// generate hash for(anchor,peak,deltaTime)
						hash := generateHash(anchorFreq, peakFreq, futureT-t)
						fingerprints = append(fingerprints, Fingerprint{
							Hash:      hash,
							TimeStamp: t,
						})
					}
				}
			}
		}
	}
	slog.Info("FingerPrint extraction completed", "count", len(fingerprints))
	return fingerprints, nil
}

// function to detectPeaks
func detectPeaks(spectogram [][]float64) map[int][]int {
	peaks := make(map[int][]int)

	for t, frame := range spectogram {
		threshold := MinPeakThreshold * maxInSlice(frame)
		var peakFreqs []int

		// Identify peaks above threshold
		for f, magnitude := range frame {
			if magnitude > threshold && isLocalPeak(frame, f) {
				peakFreqs = append(peakFreqs, f)
			}
		}
		if len(peakFreqs) > 0 {
			peaks[t] = peakFreqs
		}
	}
	return peaks
}

// isLocalPeak checks if a frequency bin is local maxima
func isLocalPeak(frame []float64, index int) bool {
	left := math.Max(0, float64(index-1))
	right := math.Min(float64(len(frame)-1), float64(index+1))

	return frame[index] > frame[int(left)] && frame[index] > frame[int(right)]
}

// maxInSlice return the max value in the given slice
func maxInSlice(frame []float64) float64 {
	if len(frame) == 0 {
		return 0
	}
	maxVal := frame[0]
	for _, v := range frame {
		if v > maxVal {
			maxVal = v
		}
	}
	return maxVal
}

// generateHash create a unique SHA-1 hash for a frequency pair and delta time (time b/w two anchor points)
func generateHash(anchorFreq, peakFreq, deltaTime int) string {
	data := fmt.Sprintf("%d|%d|%d", anchorFreq, peakFreq, deltaTime)
	hash := sha1.Sum([]byte(data))

	return hex.EncodeToString(hash[:])
}
