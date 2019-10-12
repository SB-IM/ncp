package main

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os/exec"
	"regexp"
)

type NcpCmd struct {
	config Ncp
}

func (this *NcpCmd) Init() (*NcpCmd, error) {

	// If shell enable, Run init
	if c := (*this).config.Shell; c.Path != "" {
		files, err := ioutil.ReadDir(c.Path)
		if err != nil {
			return this, err
		}

		for _, file := range files {
			if regexp.MustCompile(`_init_.*`).MatchString(file.Name()) {
				_, err := exec.Command(c.Prefix, c.Path+file.Name()).CombinedOutput()
				if err != nil {
					return this, err
				}
			}
		}
	}
	return this, nil
}

func (this *NcpCmd) Download(filename, source string) error {
	if (*this).config.Download[filename] == "" {
		return errors.New("No " + filename + " config found")
	} else {
		return httpDownload((*this).config.Common.Id, (*this).config.Common.SecretKey, (*this).config.Download[filename], source)
	}
}

func (this *NcpCmd) Upload(filename, target string) error {
	if (*this).config.Upload[filename] == "" {
		return errors.New("No " + filename + " config found")
	} else {
		return httpUpload((*this).config.Common.Id, (*this).config.Common.SecretKey, filename, (*this).config.Upload[filename], target)
	}
}

func (this *NcpCmd) Status() string {
	r, _ := json.Marshal((*this).config.Status)
	return string(r)
}

func (this *NcpCmd) Shell(command string) error {
	c := (*this).config.Shell
	if c.Path == "" {
		return errors.New("Disable shell call")
	}

	_, err := exec.Command(c.Prefix, c.Path+command+c.Suffix).CombinedOutput()
	if err != nil {
		return err
	}
	return nil
}
