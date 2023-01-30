package data

import (
	"io"
	"os"

	log "github.com/sirupsen/logrus"
)

func InData(inFile *string) ([]byte, error) {
	l := log.WithFields(log.Fields{
		"fn": "InData",
	})
	l.Debug("reading input data")
	var d []byte
	if *inFile == "-" {
		l.Debug("reading from stdin")
		var err error
		d, err = io.ReadAll(os.Stdin)
		if err != nil {
			l.Error(err)
			return nil, err
		}
	} else if *inFile != "" {
		l.WithField("file", *inFile).Debug("reading from file")
		var err error
		d, err = os.ReadFile(*inFile)
		if err != nil {
			l.Error(err)
			return nil, err
		}
	}
	return d, nil
}

func OutData(d []byte, outFile *string) error {
	l := log.WithFields(log.Fields{
		"fn": "OutData",
	})
	l.Debug("writing output data")
	if *outFile == "-" {
		l.Debug("writing to stdout")
		_, err := os.Stdout.Write(d)
		if err != nil {
			l.Error(err)
			return err
		}
	} else if *outFile != "" {
		l.WithField("file", *outFile).Debug("writing to file")
		err := os.WriteFile(*outFile, d, 0644)
		if err != nil {
			l.Error(err)
			return err
		}
	}
	return nil
}
