package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"

	tele "gopkg.in/telebot.v3"
)

var mu sync.Mutex

func execCommand(inputFilePath, outputFilePath string) error {
	mu.Lock()
	defer mu.Unlock()
	cmd := exec.Command("unoconv", "--output", outputFilePath, inputFilePath)
	_, err := cmd.Output()
	return err
}

func convertToPDF(c tele.Context) error {
	go func(c tele.Context) {
		var (
			file           = c.Message().Document
			fileFullName   = file.FileName
			fileShortName  = strings.TrimSuffix(fileFullName, filepath.Ext(fileFullName))
			fileExtension  = filepath.Ext(file.FileName)
			uniqueFileName = generateUniqueToken()
		)

		if fileExtension == ".pdf" {
			logErrorEvent(errors.New("uploaded file in PDF format"), c)
			msg := fmt.Sprintf("You have submitted PDF file [%s], no conversion was made", fileFullName)
			err := c.Send(msg)
			if err != nil {
				logErrorEvent(err, c)
				return
			}
			return
		}

		msg := fmt.Sprintf("Conversion of [%s] in progress...", fileFullName)
		err := c.Send(msg)
		if err != nil {
			logErrorEvent(err, c)
			return
		}

		logInfoEvent("downloading file", c)

		workingDirectory, err := os.Getwd()
		if err != nil {
			logErrorEvent(err, c)
			return
		}

		inputFilePath := fmt.Sprintf("%s/%s%s", workingDirectory, uniqueFileName, fileExtension)
		defer os.Remove(inputFilePath)
		err = c.Bot().Download(file.MediaFile(), fmt.Sprintf("%s%s", uniqueFileName, fileExtension))
		if err != nil {
			logErrorEvent(err, c)
			return
		}

		logInfoEvent("converting file to PDF", c)

		outputFilePath := fmt.Sprintf("%s/%s.pdf", workingDirectory, uniqueFileName)

		err = execCommand(inputFilePath, outputFilePath)
		if err != nil {
			logErrorEvent(err, c)
			msg := fmt.Sprintf("Could't convert [%s] to PDF. Unsupported format or file's size is too big", fileFullName)
			err = c.Send(msg)
			if err != nil {
				logErrorEvent(err, c)
				return
			}
			return
		}
		defer os.Remove(outputFilePath)

		logInfoEvent("sending PDF file to chat", c)

		err = c.Send(&tele.Document{
			File:     tele.FromDisk(outputFilePath),
			Caption:  fmt.Sprintf("File [%s] was successfully converted", fileFullName),
			FileName: fileShortName + ".pdf",
		})
		if err != nil {
			logErrorEvent(err, c)
			return
		}

		logInfoEvent("finished conversion", c)
	}(c)
	return nil
}
