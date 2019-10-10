package main

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"errors"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"strings"
	"time"
)

func httpDownload(access_id, secret_key, filepath, source string) error {
	req, err := http.NewRequest("GET", source, nil)
	signHeader(access_id, secret_key, req)

	res, err := (&http.Client{}).Do(req)
	defer res.Body.Close()
	if err != nil {
		return err
	}

	if res.StatusCode != http.StatusOK {
		return errors.New(res.Status)
	}

	robots, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(filepath, robots, 0644)

	if err != nil {
		return err
	}
	return nil
}

func httpUpload(access_id, secret_key, filekey, filepath, target string) error {
	file, err := os.Open(filepath)
	if err != nil {
		return err
	}

	values := map[string]io.Reader{
		filekey: file,
	}

	var b bytes.Buffer
	w := multipart.NewWriter(&b)
	for key, r := range values {
		var fw io.Writer
		err := errors.New("emit macho dwarf: elf header corrupted")
		if x, ok := r.(io.Closer); ok {
			defer x.Close()
		}
		// Add an file
		if x, ok := r.(*os.File); ok {
			if fw, err = w.CreateFormFile(key, x.Name()); err != nil {
				return err
			}
		} else {
			// Add other fields
			if fw, err = w.CreateFormField(key); err != nil {
				return err
			}
		}
		if _, err = io.Copy(fw, r); err != nil {
			return err
		}

	}
	// Don't forget to close the multipart writer.
	// If you don't close it, your request will be missing the terminating boundary.
	w.Close()

	req, err := http.NewRequest("PATCH", target, &b)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", w.FormDataContentType())

	signHeader(access_id, secret_key, req)

	res, err := (&http.Client{}).Do(req)
	if err != nil {
		return err
	}

	if res.StatusCode != http.StatusOK {
		return errors.New(res.Status)
	}
	return nil
}

func signHeader(access_id, secret_key string, req *http.Request) {
	req.Header.Add("DATE", strings.Replace(time.Now().UTC().Format(time.RFC1123), "UTC", "GMT", 1))

	// https://github.com/mgomes/api_auth/blob/master/lib/api_auth/base.rb#L115
	req.Header.Add("Authorization", "APIAuth "+access_id+":"+SignMAC([]byte(getCanonicalString(req)), []byte(secret_key)))
}

func getCanonicalString(req *http.Request) string {
	// https://github.com/mgomes/api_auth/blob/master/lib/api_auth/headers.rb#L62
	return strings.Join([]string{
		req.Method,
		req.Header.Get("Content-Type"),
		req.Header.Get("Content-Md5"),
		req.URL.Path,
		req.Header.Get("date"),
	}, ",")
}

func ValidMAC(message []byte, key []byte, bMessageMAC string) bool {
	messageMAC, _ := base64.StdEncoding.DecodeString(bMessageMAC)

	mac := hmac.New(sha1.New, key)
	mac.Write(message)
	expectedMAC := mac.Sum(nil)

	return hmac.Equal(messageMAC, expectedMAC)
}

func SignMAC(message, key []byte) string {
	mac := hmac.New(sha1.New, key)
	mac.Write(message)
	expectedMAC := mac.Sum(nil)

	return base64.StdEncoding.EncodeToString(expectedMAC)
}
