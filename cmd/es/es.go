package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/robertlestak/es/internal/data"
	"github.com/robertlestak/es/internal/env"
	"github.com/robertlestak/es/pkg/es"
	log "github.com/sirupsen/logrus"
)

func printVars(inFile *string, vals bool) {
	l := log.WithFields(log.Fields{
		"fn": "printVars",
	})
	l.Debug("printing variables")
	var keys []string
	var err error
	var isDir bool
	if *inFile != "-" {
		// check if the input file is a directory
		fi, err := os.Stat(*inFile)
		if err != nil {
			l.Error(err)
			os.Exit(1)
		}
		isDir = fi.IsDir()
	}
	if isDir {
		keys, err = es.EnvKeysInDir(*inFile)
		if err != nil {
			l.Error(err)
			os.Exit(1)
		}
	} else {
		d, err := data.InData(inFile)
		if err != nil {
			l.Error(err)
			os.Exit(1)
		}
		keys = es.EnvKeysFromIn(d)
	}
	for _, k := range keys {
		if vals {
			fmt.Fprintf(os.Stdout, "%s=%s\n", k, os.Getenv(k))
		} else {
			fmt.Fprintf(os.Stdout, "%s\n", k)
		}
	}
}

func main() {
	l := log.WithFields(log.Fields{
		"fn": "main",
	})
	l.Debug("starting")
	requireVars := flag.Bool("r", false, "require all variables to be set")
	inFile := flag.String("i", "-", "input file")
	outFile := flag.String("o", "-", "output file")
	var envFiles []string
	flag.Func("e", "environment file", func(s string) error {
		envFiles = append(envFiles, s)
		return nil
	})
	printVariables := flag.Bool("v", false, "print variables")
	printVariableVals := flag.Bool("vv", false, "print variables and values")
	logLevel := flag.String("l", "info", "log level")
	flag.Parse()
	ll, err := log.ParseLevel(*logLevel)
	if err != nil {
		l.Fatal(err)
	}
	log.SetLevel(ll)
	if len(envFiles) > 0 {
		if err := env.ReadEnvFiles(envFiles); err != nil {
			l.Fatal(err)
		}
	}
	// get first arg, if it exists
	if flag.NArg() > 0 {
		*inFile = flag.Arg(0)
	}
	if *printVariables || *printVariableVals {
		printVars(inFile, *printVariableVals)
		os.Exit(0)
	}
	var isDir bool
	if *inFile != "-" {
		// check if the input file is a directory
		fi, err := os.Stat(*inFile)
		if err != nil {
			l.Error(err)
			os.Exit(1)
		}
		isDir = fi.IsDir()
	}
	if isDir {
		// if outdir is empty, create a temp dir
		var createdDir bool
		if *outFile == "-" {
			l.Debug("creating temp dir")
			var err error
			*outFile, err = os.MkdirTemp("", "es")
			if err != nil {
				l.Fatal(err)
			}
			l = l.WithField("o", *outFile)
			l.Debug("created temp dir")
			createdDir = true
		}
		if *inFile == *outFile {
			l.Fatal("input and output directories cannot be the same")
		}
		err := es.ProcessDir(inFile, outFile, *requireVars)
		if err != nil {
			os.Exit(1)
		}
		if createdDir {
			fmt.Println(*outFile)
		}
	} else {
		d, err := data.InData(inFile)
		if err != nil {
			l.Error(err)
			os.Exit(1)
		}
		d, err = es.ProcessData(d, *requireVars)
		if err != nil {
			os.Exit(1)
		}
		err = data.OutData(d, outFile)
		if err != nil {
			l.Error(err)
			os.Exit(1)
		}
	}
}
