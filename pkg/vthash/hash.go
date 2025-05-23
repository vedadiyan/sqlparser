/*
Copyright 2023 The Vitess Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package vthash

import (
	"github.com/vedadiyan/sqlparser/pkg/vthash/highway"
	"github.com/vedadiyan/sqlparser/pkg/vthash/metro"
)

type Hasher = metro.Metro128
type Hash = [16]byte

func New() Hasher {
	h := Hasher{}
	h.Reset()
	return h
}

type Hasher256 = highway.Digest
type Hash256 = [32]byte

var defaultHash256Key = [32]byte{}

func New256() *Hasher256 {
	return highway.New(defaultHash256Key)
}
