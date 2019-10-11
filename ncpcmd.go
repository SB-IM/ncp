package main

import (
	"encoding/json"
	"errors"
	"os/exec"
)

type NcpCmd struct {
	config Ncp
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
