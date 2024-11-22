package helpers

import (
	"bytes"
	"context"
	"fmt"
	"mime/multipart"
	"os"
	"os/exec"
	"time"

	domainApp "github.com/aldinokemal/go-whatsapp-web-multidevice/domains/app"
	"go.mau.fi/whatsmeow"
)

func SetAutoConnectAfterBooting(service domainApp.IAppService) {
	time.Sleep(2 * time.Second)
	_ = service.Reconnect(context.Background())
}

func SetAutoReconnectChecking(cli *whatsmeow.Client) {
	// Run every 5 minutes to check if the connection is still alive, if not, reconnect
	go func() {
		for {
			time.Sleep(5 * time.Minute)
			if !cli.IsConnected() {
				_ = cli.Connect()
			}
		}
	}()
}

func MultipartFormFileHeaderToBytes(fileHeader *multipart.FileHeader) []byte {
	file, _ := fileHeader.Open()
	defer file.Close()

	fileBytes := make([]byte, fileHeader.Size)
	_, _ = file.Read(fileBytes)

	return fileBytes
}

// ConvertToOggOpus converte um arquivo de áudio para o formato OGG com codificação Opus
func ConvertToOggOpus(inputBytes []byte) ([]byte, error) {
	// Cria um arquivo temporário para o input
	inputFile, err := os.CreateTemp("", "input-*.mp3")
	if err != nil {
		return nil, fmt.Errorf("erro ao criar arquivo temporário de entrada: %v", err)
	}
	defer os.Remove(inputFile.Name())

	// Escreve os bytes de entrada no arquivo temporário
	if _, err := inputFile.Write(inputBytes); err != nil {
		return nil, fmt.Errorf("erro ao escrever no arquivo temporário de entrada: %v", err)
	}
	inputFile.Close()

	// Cria um arquivo temporário para o output
	outputFile, err := os.CreateTemp("", "output-*.ogg")
	if err != nil {
		return nil, fmt.Errorf("erro ao criar arquivo temporário de saída: %v", err)
	}
	defer os.Remove(outputFile.Name())
	outputFile.Close()

	// Executa o comando FFmpeg para converter o áudio
	cmd := exec.Command("ffmpeg", "-i", inputFile.Name(), "-c:a", "libopus", outputFile.Name())
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("erro ao converter áudio: %v, %s", err, stderr.String())
	}

	// Lê os bytes do arquivo de saída
	outputBytes, err := os.ReadFile(outputFile.Name())
	if err != nil {
		return nil, fmt.Errorf("erro ao ler arquivo de saída: %v", err)
	}

	return outputBytes, nil
}

// GetAudioDuration retorna a duração de um arquivo de áudio em segundos
func GetAudioDuration(inputBytes []byte, extension string) (uint32, error) {
	// Cria um arquivo temporário para o input
	inputFile, err := os.CreateTemp("", "input-*."+extension)
	if err != nil {
		return 0, fmt.Errorf("erro ao criar arquivo temporário de entrada: %v", err)
	}
	defer os.Remove(inputFile.Name())

	// Escreve os bytes de entrada no arquivo temporário
	if _, err := inputFile.Write(inputBytes); err != nil {
		return 0, fmt.Errorf("erro ao escrever no arquivo temporário de entrada: %v", err)
	}
	inputFile.Close()

	// Executa o comando FFmpeg para obter a duração do áudio
	cmd := exec.Command("ffprobe", "-v", "error", "-show_entries", "format=duration", "-of", "default=noprint_wrappers=1:nokey=1", inputFile.Name())
	var stdout bytes.Buffer
	cmd.Stdout = &stdout
	if err := cmd.Run(); err != nil {
		return 0, fmt.Errorf("erro ao obter duração do áudio: %v", err)
	}

	// Converte a saída para float64
	duration, err := time.ParseDuration(stdout.String() + "s")
	if err != nil {
		return 0, fmt.Errorf("erro ao converter duração do áudio: %v", err)
	}
	// retorna em Uint32
	return uint32(duration.Seconds()), nil

}

// GetAudioWaveform retorna o waveform de um arquivo de áudio
func GetAudioWaveform(inputBytes []byte, extension string) ([]int, error) {
	// Cria um arquivo temporário para o input
	inputFile, err := os.CreateTemp("", "input-*."+extension)
	if err != nil {
		return nil, fmt.Errorf("erro ao criar arquivo temporário de entrada: %v", err)
	}
	defer os.Remove(inputFile.Name())

	// Escreve os bytes de entrada no arquivo temporário
	if _, err := inputFile.Write(inputBytes); err != nil {
		return nil, fmt.Errorf("erro ao escrever no arquivo temporário de entrada: %v", err)
	}
	inputFile.Close()

	// Executa o comando FFmpeg para obter o waveform do áudio
	cmd := exec.Command("ffmpeg", "-i", inputFile.Name(), "-filter_complex", "showwavespic=s=640x120", "-frames:v", "1", "-f", "image2", "-")
	var stdout bytes.Buffer
	cmd.Stdout = &stdout
	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("erro ao obter waveform do áudio: %v", err)
	}

	// Converte a saída para bytes
	waveformBytes := stdout.Bytes()

	// Converte os bytes para inteiros
	waveform := make([]int, len(waveformBytes))
	for i, b := range waveformBytes {
		waveform[i] = int(b)
	}

	return waveform, nil
}
