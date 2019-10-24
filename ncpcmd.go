package main

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os/exec"
	"regexp"

	"light"
)

type NcpCmd struct {
	config Ncp
	webrtc light.Star
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

func (this *NcpCmd) Download(filename, source string) ([]byte, error) {
	if (*this).config.Download[filename] == "" {
		return []byte(""), errors.New("No " + filename + " config found")
	} else {
		return []byte(""), httpDownload((*this).config.Common.Id, (*this).config.Common.SecretKey, (*this).config.Download[filename], source)
	}
}

func (this *NcpCmd) Upload(filename, target string) ([]byte, error) {
	if (*this).config.Upload[filename] == "" {
		return []byte(""), errors.New("No " + filename + " config found")
	} else {
		return []byte(""), httpUpload((*this).config.Common.Id, (*this).config.Common.SecretKey, filename, (*this).config.Upload[filename], target)
	}
}

func (this *NcpCmd) Status() ([]byte, error) {
	return json.Marshal((*this).config.Status)
}

func (this *NcpCmd) Shell(command string) ([]byte, error) {
	c := (*this).config.Shell
	if c.Path == "" {
		return []byte(""), errors.New("Disable shell call")
	}

	_, err := exec.Command(c.Prefix, c.Path+command+c.Suffix).CombinedOutput()
	if err != nil {
		return []byte(""), err
	}
	return []byte(""), nil
}

func (this *NcpCmd) Webrtc(raw []byte) ([]byte, error) {
	return *(this.webrtc.Light(&(*this).config.Webrtc.Iceserver, (*this).config.Webrtc.Args, &raw)), nil
}
