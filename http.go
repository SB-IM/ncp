package main

import (
//  "fmt"
  "crypto/sha1"
  "crypto/hmac"
  "encoding/base64"
  //"net/url"

  "fmt"
  "io/ioutil"
  "log"
  "net/http"
  "time"
  "strings"
  //"bufio"
  "bytes"
  "mime/multipart"
  "os"
  "io"
)

func main() {
  //download()
  pre_upload()

}


func Upload(client *http.Client, url string, values map[string]io.Reader) (err error) {
  // Prepare a form that you will submit to that URL.
  var b bytes.Buffer
  w := multipart.NewWriter(&b)
  for key, r := range values {
    var fw io.Writer
    if x, ok := r.(io.Closer); ok {
      defer x.Close()
    }
    // Add an image file
    if x, ok := r.(*os.File); ok {
      if fw, err = w.CreateFormFile(key, x.Name()); err != nil {
        return
      }
    } else {
      // Add other fields
      if fw, err = w.CreateFormField(key); err != nil {
        return
      }
    }
    if _, err = io.Copy(fw, r); err != nil {
      return err
    }

  }
  // Don't forget to close the multipart writer.
  // If you don't close it, your request will be missing the terminating boundary.
  w.Close()

  // Now that you have a form, you can submit it to your handler.
  //req, err := http.NewRequest("POST", url, &b)
  req, err := http.NewRequest("PATCH", url, &b)
  if err != nil {
    return
  }
  // Don't forget to set the content type, this will contain the boundary.
  req.Header.Set("Content-Type", w.FormDataContentType())

  //--------------

  date := strings.Replace(time.Now().UTC().Format(time.RFC1123), "UTC", "GMT", 1)
  req.Header.Add("DATE", date)

  access_id, secret_key := "2", "cJXWEPknyfzAZPOmQX6/mbpGqSuzxYD9aezcDHezy4tVK7U94vfxLEObcC9yD4TRfBVddg/ir7XzDDDTn7GFMA=="
  //fmt.Println(req.Header["Content-Type"][0])
  fmt.Println(req.Header)
  //hhh := req.Header
  //fmt.Println(hhh.get("Content-Type"))

  str := req.Method + "," + req.Header["Content-Type"][0] + ",," + req.URL.Path + "," + date

  req.Header.Add("Authorization", "APIAuth " + access_id + ":" + SignMAC([]byte(str), []byte(secret_key)))
  //---------------

  // Submit the request
  res, err := client.Do(req)
  if err != nil {
    return
  }

  // Check the response
  if res.StatusCode != http.StatusOK {
    err = fmt.Errorf("bad status: %s", res.Status)
  }
  return
}


func mustOpen(f string) *os.File {
  r, err := os.Open(f)
  if err != nil {
    panic(err)
  }
  return r
}

func pre_upload() {
  //var client *http.Client
  client := &http.Client{}

  values := map[string]io.Reader{
    "air_log":  mustOpen("rrr.txt"), // lets assume its this file
  }

  remoteURL := "http://localhost:3000/ncp/v1/plans/13/plan_logs/37"
  err := Upload(client, remoteURL, values)
  if err != nil {
    panic(err)
  }

}

//
//func upload() {
//  fmt.Println("upload")
//
//  //content, err := ioutil.ReadFile("rrr.txt")
//  //if err != nil {
//  //  log.Fatal(err)
//  //}
//
//  //fmt.Printf("File contents: %s", content)
//
//  client := &http.Client{}
//
//  //file, err := os.Open("rrr.txt")
//  file, _ := os.Open("rrr.txt")
//
//  content := map[string]io.Reader{
//    "air_log": bufio.NewReader(file), // lets assume its this file
//  }
//
//  //content := bufio.NewReader(file)
//
//  req, _ := http.NewRequest("PATCH", "http://localhost:3000/ncp/v1/plans/13/plan_logs/37", content)
//  date := strings.Replace(time.Now().UTC().Format(time.RFC1123), "UTC", "GMT", 1)
//  //fmt.Printf(strings.Replace(time.Now().UTC().Format(time.RFC1123), "UTC", "GMT", 1))
//  //date := time.Now().UTC().Format("Mon, 02 Jan 2006 15:04:05 GMT")
//  //date := "Mon, 17 Jun 2019 09:32:09 GMT"
//  fmt.Println(date)
//  fmt.Println("-----------")
//  req.Header.Add("DATE", date)
//
//  access_id, secret_key := "2", "cJXWEPknyfzAZPOmQX6/mbpGqSuzxYD9aezcDHezy4tVK7U94vfxLEObcC9yD4TRfBVddg/ir7XzDDDTn7GFMA=="
//  str := req.Method + ",,," + req.URL.Path + "," + date
//
//  req.Header.Add("Authorization", "APIAuth " + access_id + ":" + SignMAC([]byte(str), []byte(secret_key)))
//
//  fmt.Println(req.Method)
//  fmt.Println(req)
//  //res, err := client.Do(req)
//  res, _ := client.Do(req)
//
//  fmt.Println("==============")
//  fmt.Println(res)
//
//
//}
//
func download() {
  client := &http.Client{}

  req, err := http.NewRequest("GET", "http://localhost:3000/ncp/v1/plans/12/get_map", nil)
  date := strings.Replace(time.Now().UTC().Format(time.RFC1123), "UTC", "GMT", 1)
  //fmt.Printf(strings.Replace(time.Now().UTC().Format(time.RFC1123), "UTC", "GMT", 1))
  //date := time.Now().UTC().Format("Mon, 02 Jan 2006 15:04:05 GMT")
  //date := "Mon, 17 Jun 2019 09:32:09 GMT"
  fmt.Println(date)
  fmt.Println("-----------")
  req.Header.Add("DATE", date)

  access_id, secret_key := "2", "cJXWEPknyfzAZPOmQX6/mbpGqSuzxYD9aezcDHezy4tVK7U94vfxLEObcC9yD4TRfBVddg/ir7XzDDDTn7GFMA=="
  str := req.Method + ",,," + req.URL.Path + "," + date

  req.Header.Add("Authorization", "APIAuth " + access_id + ":" + SignMAC([]byte(str), []byte(secret_key)))

  fmt.Println(req.Method)
  fmt.Println(req.URL.Path)
  res, err := client.Do(req)

  //fmt.Println(res)

  if err != nil {
    log.Fatal(err)
  }
  robots, err := ioutil.ReadAll(res.Body)
  err = ioutil.WriteFile("rrr.txt", robots, 0644)
  res.Body.Close()
  if err != nil {
    log.Fatal(err)
  }
  //fmt.Printf("%s", robots)
}

//
//func main() {
//  access_id, secret_key := "2", "cJXWEPknyfzAZPOmQX6/mbpGqSuzxYD9aezcDHezy4tVK7U94vfxLEObcC9yD4TRfBVddg/ir7XzDDDTn7GFMA=="
//  date := "Mon, 17 Jun 2019 03:10:03 GMT"
//  str := "GET,,,/ncp/v1/plans/12/get_map,Mon, 17 Jun 2019 03:10:03 GMT"
//
//  authorization := "APIAuth 2:SuX8+f6zt5TxYoTe0cuo9PALBbY="
//  fmt.Println(access_id)
//  fmt.Println(secret_key)
//  fmt.Println(date)
//  fmt.Println(str)
//  fmt.Println(authorization)
//
//
//  fmt.Println("23333333333")
//  v, _ := base64.StdEncoding.DecodeString("SuX8+f6zt5TxYoTe0cuo9PALBbY=")
//  //a := ValidMAC([]byte(str), []byte(v), []byte(secret_key))
//  a := ValidMAC([]byte(str), []byte(secret_key), "SuX8+f6zt5TxYoTe0cuo9PALBbY=")
//  fmt.Println(a)
//}
//

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

