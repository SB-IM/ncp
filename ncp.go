package main

import (
  "fmt"
  "encoding/json"
  "path/filepath"
  "strings"
  "reflect"
  "runtime"
)

//type Ncp struct {
//}

func (ncp *Ncp) Method_status() string {
  b, _ := json.Marshal(ncp.Status)
  return string(b)
  //return ncp.Status
}

func (ncp *Ncp) Method_upload(file, path string) {
  fmt.Println(MethodName())
  fmt.Println(ncp.Upload)

  //return ""
}

func MethodName() string {
  pc, _, _, _ := runtime.Caller(1)
  //return runtime.FuncForPC(pc).Name()
  return strings.TrimPrefix(filepath.Ext(runtime.FuncForPC(pc).Name()), ".")
}


func CallObjectMethod(object interface{}, methodName string, args ...interface{}) interface{} {
  inputs := make([]reflect.Value, len(args))
  for i, _ := range args {
    inputs[i] = reflect.ValueOf(args[i])
  }
  return reflect.ValueOf(object).MethodByName(methodName).Call(inputs)
}

//func main() {
//  CallObjectMethod(new(Ncp), "Method_status")
//  CallObjectMethod(new(Ncp), "Method_upload", "map", "http")
//}


