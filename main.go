package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
)

func main() {
	// Datei öffnen
	file, err := os.Open("messages.txt")
	if err != nil {
		log.Fatalf("Fehler beim Öffnen der Datei: %v", err)
	}

	// Channel holen
	lines := getLinesChannel(file)

	// Zeilen lesen und im gewünschten Format ausgeben
	for line := range lines {
		fmt.Printf("read: %s\n", line)
	}
}

func getLinesChannel(f io.ReadCloser) <-chan string {
	out := make(chan string)

	go func() {
		defer f.Close()
		defer close(out)

		buffer := make([]byte, 8)
		lineBuffer := make([]byte, 0)

		for {
			n, err := f.Read(buffer)
			if n > 0 {
				lineBuffer = append(lineBuffer, buffer[:n]...)

				for {
					index := bytes.IndexByte(lineBuffer, '\n')
					if index == -1 {
						break
					}

					line := string(lineBuffer[:index+1])
					out <- line[:len(line)-1] // '\n' am Ende entfernen
					lineBuffer = lineBuffer[index+1:]
				}
			}

			if err != nil {
				if err == io.EOF {
					if len(lineBuffer) > 0 {
						out <- string(lineBuffer)
					}
					break
				}
				// Bei echtem Fehler abbrechen
				log.Printf("Lesefehler: %v", err)
				break
			}
		}
	}()

	return out
}
