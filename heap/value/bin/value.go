// Licensed to the Apache Software Foundation (ASF) under one or more
// contributor license agreements.  See the NOTICE file distributed with
// this work for additional information regarding copyright ownership.
// The ASF licenses this file to You under the Apache License, Version 2.0
// (the "License"); you may not use this file except in compliance with
// the License.  You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"flag"
	"fmt"
)

var arr1 []struct1
var arr2 []*struct1
var n = 1000000

type struct1 struct {
	f1 int
	f2 string
}

func funcValue(s struct1) {
	arr1 = append(arr1, s)
}

func funcPointer(s *struct1) {
	arr2 = append(arr2, s)
}

func TestValue() {
	fmt.Println("run TestValue")
	for i := 0; i < n; i++ {
		obj := struct1{f1: i, f2: "abc"}
		funcValue(obj)
	}
}

func TestPointer() {
	fmt.Println("run TestPointer")
	for i := 0; i < n; i++ {
		obj := &struct1{f1: i, f2: "abc"}
		funcPointer(obj)
	}
}

func main() {
	flag.Parse()
	switch flag.Arg(0) {
	case "value":
		arr1 = make([]struct1, 0, n)
		TestValue()
	default:
		arr2 = make([]*struct1, 0, n)
		TestPointer()
	}
}
