package es

import (
	"errors"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	log "github.com/sirupsen/logrus"
)

// envKeysFromIn takes a byte slice and returns a slice of strings
// it searches the input for environment variables (in the form $VARNAME or ${VARNAME})
// and returns a slice of strings containing the variable names
func EnvKeysFromIn(in []byte) []string {
	var keys []string
	re := regexp.MustCompile(`\$\{?(\w+)\}?`)
	matches := re.FindAllSubmatch(in, -1)
	for _, m := range matches {
		keys = append(keys, string(m[1]))
	}
	return removeDuplicates(keys)
}

func hasAllVars(keys []string) bool {
	log.Debug("checking for environment variables")
	for _, k := range keys {
		if os.Getenv(k) == "" {
			log.Errorf("missing variable %s", k)
			return false
		}
	}
	return true
}

func ProcessData(d []byte, require bool) ([]byte, error) {
	l := log.WithFields(log.Fields{
		"fn": "processData",
	})
	l.Debug("processing data")
	keys := EnvKeysFromIn(d)
	l.WithField("keys", keys).Debug("found keys")
	if require && !hasAllVars(keys) {
		return nil, errors.New("missing variables")
	}
	return []byte(os.ExpandEnv(string(d))), nil
}

func removeDuplicates(s []string) []string {
	l := log.WithFields(log.Fields{
		"fn": "removeDuplicates",
		"s":  s,
	})
	l.Debug("removing duplicates")
	keys := make(map[string]bool)
	var list []string
	for _, entry := range s {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}

func EnvKeysInDir(d string) ([]string, error) {
	l := log.WithFields(log.Fields{
		"fn": "envKeysInDir",
		"d":  d,
	})
	l.Debug("getting environment keys in directory")
	var keys []string
	err := filepath.Walk(d, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			l.Error(err)
			return err
		}
		l.WithField("path", path).Debug("processing file")
		// skip directories
		if info.IsDir() {
			l.Debug("skipping directory")
			return nil
		}
		// read the file to a byte slice
		bd, err := os.ReadFile(path)
		if err != nil {
			l.Error(err)
			return err
		}
		// get the keys from the byte slice
		keys = append(keys, EnvKeysFromIn(bd)...)
		return nil
	})
	if err != nil {
		l.Error(err)
		return nil, err
	}
	// remove duplicates
	keys = removeDuplicates(keys)
	l.WithField("keys", keys).Debug("found keys")
	return keys, nil
}

func ProcessDir(inDir *string, outDir *string, require bool) error {
	l := log.WithFields(log.Fields{
		"fn": "ProcessDir",
		"i":  *inDir,
		"o":  *outDir,
		"r":  require,
	})
	l.Debug("processing directory")
	// walk the input dir
	l.Debug("walking input dir")
	err := filepath.Walk(*inDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			l.Error(err)
			return err
		}
		l.WithField("path", path).Debug("processing file")
		// skip directories
		if info.IsDir() {
			l.Debug("skipping directory")
			return nil
		}
		// read the file to a byte slice
		bd, err := os.ReadFile(path)
		if err != nil {
			l.Error(err)
			return err
		}
		// process the data
		d, err := ProcessData(bd, require)
		if err != nil {
			return err
		}
		// write the data to the output dir
		opath := filepath.Join(*outDir, strings.TrimPrefix(path, *inDir))
		l.WithField("opath", opath).Debug("writing file")
		if err := os.MkdirAll(filepath.Dir(opath), 0755); err != nil {
			l.Error(err)
			return err
		}
		if err := os.WriteFile(opath, d, 0644); err != nil {
			l.Error(err)
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}
