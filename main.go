package main

import (
	"flag"
	"fmt"
	"io"
	"logwatcher/client"
	"logwatcher/extractor"
	"os"
	"strings"
	"time"

	"golang.org/x/xerrors"
)

var (
	filePath           = flag.String("file", "", "File path")
	regexString        = flag.String("match", "", "Match string")
	outputFormat       = flag.String("outputFormat", "", "Output format")
	sendTarget         = flag.String("sendTarget", "discord", "Send target")
	discordWebhookURL  = flag.String("discordWebhookURL", "", "Discord webhook URL")
	intervalTimeString = flag.String("intervalTime", "1s", "Interval time")
)

func main() {
	flag.Parse()

	fmt.Println("File path: ", *filePath)
	fmt.Println("Match string: ", *regexString)
	fmt.Println("Output format: ", *outputFormat)

	file, err := os.Open(*filePath)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	file.Seek(0, 2)

	var handler func(string) error
	switch *sendTarget {
	case "discord":
		handler = client.SendToDiscordWebhookHandler(*discordWebhookURL)
	default:
		panic("Invalid send target")
	}

	intervalTime, err := time.ParseDuration(*intervalTimeString)
	if err != nil {
		panic(err)
	}

	err = watchFile(file, *regexString, *outputFormat, handler, intervalTime)
	if err != nil {
		fmt.Printf("Error watching file: %+v", err)
	}
}

func watchFile(file *os.File, regexString string, outputFormat string, handler func(string) error, intervalTime time.Duration) error {
	for {
		time.Sleep(intervalTime)

		bytes, err := io.ReadAll(file)
		if err != nil {
			return xerrors.Errorf("failed to read file: %w", err)
		}

		text := string(bytes)
		outputs, err := extractor.ExtractAndReplaceText(text, regexString, outputFormat)
		if err != nil {
			return xerrors.Errorf("failed to extract text: %w", err)
		}

		if len(outputs) != 0 {
			fmt.Printf("%d hits\n", len(outputs))
		}

		output := strings.Join(outputs, "\n")
		if output != "" {
			err = handler(output)
			if err != nil {
				fmt.Printf("Error handdling output: %+v", err)
			}
		}

		file.Seek(0, 2)
	}
}
