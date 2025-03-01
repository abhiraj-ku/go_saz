package audio

import (
	"encoding/binary"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"time"

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
	rec.stream = stream
	return rec, nil

}

// Start --> initiate recording of audio from source with specified duration (5s)
func (r *Recorder) Start(duration time.Duration) error {
	if r.stream == nil {
		return errors.New("audio stream not initialized ")
	}

	slog.Info("Recording started", "duration", duration)

	if err := r.stream.Start(); err != nil {
		return fmt.Errorf("failed to start recording: %w", err)
	}

	time.Sleep(duration)

	if err := r.stream.Stop(); err != nil {
		return fmt.Errorf("failed to stop recording: %w", err)
	}

	slog.Info("Recording finish")
	return nil
}

// Stop --> closesand relases the resources
func (r *Recorder) Close() {
	if r.stream != nil {
		r.stream.Close()
	}
	if r.file != nil {
		updateWAVHeader(r.file)
		r.file.Close()
	}

	portaudio.Terminate()
}

// WAV Header function --> writeWAVHeader

func writeWAVHeader(f *os.File) error {
	header := make([]byte, 44)
	copy(header[:4], "RIFF")
	copy(header[8:12], "WAVE")
	copy(header[12:16], "fmt ")
	binary.LittleEndian.PutUint32(header[16:20], 16) // PCM format
	binary.LittleEndian.PutUint16(header[20:22], 1)  // Audio format (PCM)
	binary.LittleEndian.PutUint16(header[22:24], Channels)
	binary.LittleEndian.PutUint32(header[24:28], SampleRate)
	binary.LittleEndian.PutUint32(header[28:32], SampleRate*Channels*(BitDepth/8))
	binary.LittleEndian.PutUint16(header[32:34], Channels*(BitDepth/8))
	binary.LittleEndian.PutUint16(header[34:36], BitDepth)
	copy(header[36:40], "data")
	_, err := f.Write(header)
	return err
}

// update WAV Header with file BufferSize
func updateWAVHeader(f *os.File) {
	info, err := f.Stat()
	if err != nil {
		slog.Error("Failed to get file info: ", "error", err)
		return
	}
	fileSize := uint32(info.Size())
	dataSize := fileSize - 44

	// Update RIFF chunk size
	f.Seek(4, 0)
	binary.Write(f, binary.LittleEndian, fileSize-8)

	// Update data chunk size
	f.Seek(40, 0)
	binary.Write(f, binary.LittleEndian, dataSize)
}
