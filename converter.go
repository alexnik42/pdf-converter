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

func convertToPDF(c tele.Context) error {
	var (
		file           = c.Message().Document
		fileName       = strings.TrimSuffix(file.FileName, filepath.Ext(file.FileName))
		fileExtension  = filepath.Ext(file.FileName)
		uniqueFileName = generateUniqueToken()
	)

	if fileExtension == ".pdf" {
		c.Send("You have submitted pdf file, no conversion was made")
		return logErrorEvent(errors.New("uploaded file in pdf format"), c)
	}

	c.Send("Conversion in progress...")

	logInfoEvent("Downloading file", c)

	workingDirectory, err := os.Getwd()
	if err != nil {
		return logErrorEvent(err, c)
	}

	inputFilePath := fmt.Sprintf("%s/%s%s", workingDirectory, uniqueFileName, fileExtension)
	defer os.Remove(inputFilePath)
	err = c.Bot().Download(file.MediaFile(), fmt.Sprintf("%s%s", uniqueFileName, fileExtension))
	if err != nil {
		return logErrorEvent(err, c)
	}

	logInfoEvent("Converting file to pdf", c)

	outputFilePath := fmt.Sprintf("%s/%s.pdf", workingDirectory, uniqueFileName)

	cmd := exec.Command("unoconv", "--output", outputFilePath, inputFilePath)
	_, err = cmd.Output()
	defer os.Remove(outputFilePath)
	if err != nil {
		c.Send("Could't convert to PDF. Unsupported format or file's size is too big")
		return logErrorEvent(err, c)
	}

	logInfoEvent("Sending pdf file to chat", c)

	err = c.Send(&tele.Document{
		File:     tele.FromDisk(outputFilePath),
		Caption:  "File was successfully converted",
		FileName: fileName + ".pdf",
	})
	if err != nil {
		return logErrorEvent(err, c)
	}

	logInfoEvent("Finished conversion", c)
	return nil
}
