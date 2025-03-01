package audio

import (
	"encoding/binary"
	"fmt"
	"log/slog"
	"os"

	"github.com/gordonklaus/portaudio"
)

// Audio settings
const (
	SampleRate = 44100 // CD Quality
	Channels   = 1     //  Mono-> 1
	BitDepth   = 16    // 	 16-Bit PCM
	BufferSize = 4096  // Number of samples per buffer
)

// Recorder struct
type Recorder struct {
	stream *portaudio.Stream
	file   *os.File
}

// New Recorder inits the Recorder

func NewRecorder(filename string) (*Recorder, error) {
	if err := portaudio.Initialize(); err != nil {
		return nil, fmt.Errorf("failed to initialize the PortAudio recorder: ", err)
	}

	file, err := os.Create(filename)
	if err != nil {
		portaudio.Terminate()
		return nil, fmt.Errorf("Failed to create the audio file :", err)
	}

	// empty WAV(Waveform audio file format)
	if err := writeWAVHeader(file); err != nil {
		file.Close()
		portaudio.Terminate()
		return nil, fmt.Errorf("failed to write WAV Header: ", err)
	}

	// create Audio stream
	rec := &Recorder{file: file}
	buffer := make([]int16, BufferSize)

	// open stream for recording
	stream, err := portaudio.OpenDefaultStream(1, 0, SampleRate, BufferSize, func(in []int16) {
		if err := binary.Write(rec.file, binary.LittleEndian, in); err != nil {
			slog.Error("failed to write audio data", "error", err)
		}
	})

	if err != nil {
		file.Close()
		portaudio.Terminate()
		return nil, fmt.Errorf("failed to open audio stream: ", err)

	}

}
