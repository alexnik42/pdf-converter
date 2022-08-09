package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	tele "gopkg.in/telebot.v3"
)

func convertToPdf(c tele.Context) error {
	var (
		file                     = c.Message().Document
		fileNameWithExtension    = file.FileName
		fileNameWithoutExtension = strings.TrimSuffix(fileNameWithExtension, filepath.Ext(fileNameWithExtension))
	)

	c.Send("Conversion in progress...")

	logInfoEvent("Downloading file", c)

	workingDirectory, err := os.Getwd()
	if err != nil {
		return logErrorEvent(err, c)
	}

	inputFilePath := fmt.Sprintf("%s/%s", workingDirectory, fileNameWithExtension)
	defer os.Remove(inputFilePath)
	err = c.Bot().Download(file.MediaFile(), fileNameWithExtension)
	if err != nil {
		return logErrorEvent(err, c)
	}

	logInfoEvent("Converting file to pdf", c)

	outputFilePath := fmt.Sprintf("%s/%s.pdf", workingDirectory, fileNameWithoutExtension)
	if outputFilePath == inputFilePath {
		c.Send("You have submitted pdf file, no conversion was made")
		return logErrorEvent(errors.New("uploaded file in pdf format"), c)
	}

	cmd := exec.Command("unoconv", "--output", outputFilePath, inputFilePath)
	_, err = cmd.Output()
	defer os.Remove(outputFilePath)
	if err != nil {
		c.Send("Could't convert to pdf")
		return logErrorEvent(err, c)
	}

	logInfoEvent("Sending pdf file to chat", c)

	err = c.Send(&tele.Document{
		File:     tele.FromDisk(outputFilePath),
		Caption:  "File was successfully converted",
		FileName: fileNameWithoutExtension + ".pdf",
	})
	if err != nil {
		return logErrorEvent(err, c)
	}

	logInfoEvent("Finished conversion", c)
	return nil
}
