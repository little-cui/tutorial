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

type parent struct {
	c *child
}

func (s *parent) set(f *child) {
	s.c = f
}

type child struct {
}

func value(s parent) parent {
	return s
}

func Case1() {
	var s parent
	_ = value(s)
}

func pointer1(s *parent) *parent {
	return s
}

func Case2() {
	var s parent
	_ = *pointer1(&s)
}

func pointer2(s parent) *parent {
	return &s
}

func Case3() {
	var s parent
	_ = *pointer2(s)
}

func pointer3(in *child, out *parent) {
	out.c = in
}

func Case4() {
	var s parent
	var i child
	pointer3(&i, &s)
}

func maps1() map[int]int {
	return make(map[int]int, 10)
}

func CaseMap1() {
	_ = maps1()
}

func CaseMap2() {
	m := make(map[int]*int)
	var i int
	m[i] = &i
}

func CaseMap3() {
	m := make(map[int]int)
	var i int
	m[i] = i
}

func slice1() []int {
	var a [10]int
	return a[:]
}

func CaseSlice1() {
	_ = slice1()
}

func CaseSlice2() {
	var a []*int
	var i int
	a = append(a, &i)
}

func CaseSlice3() {
	var a []int
	var i int
	a = append(a, i)
}

func CaseChan1() {
	ch := make(chan parent, 1)
	ch <- parent{}
}

func CaseChan2() {
	ch := make(chan *child, 1)
	ch <- &child{}
}

type interface1 interface {
	set(s *child)
}

func set1(s *parent, i *child) {
	s.set(i)
}

func set2(s interface1, i *child) {
	s.set(i)
}

func CaseInterface1() {
	var s parent
	var i child
	set1(&s, &i)
}

func CaseInterface2() {
	var s parent
	var i child
	set2(&s, &i)
}

func main() {
	Case1()
	Case2()
	Case3()    // escape!!!
	Case4()    // escape!!!
	CaseMap1() // escape!!!
	CaseMap2() // escape!!!
	CaseMap3()
	CaseSlice1() // escape!!!
	CaseSlice2() // escape!!!
	CaseSlice3()
	CaseChan1()
	CaseChan2() // escape!!!
	CaseInterface1()
	CaseInterface2() // escape!!!
}
