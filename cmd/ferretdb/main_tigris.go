// Copyright 2021 FerretDB Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

//go:build ferretdb_tigris

package main

import (
	"github.com/alecthomas/kong"
)

// The tigrisFlags struct represents flags that
// are used specifically for a "tigris" handler.
var tigrisFlags struct {
	TigrisClientID     string `default:"" help:"Tigris Client ID."`
	TigrisClientSecret string `default:"" help:"Tigris Client secret."`
	TigrisToken        string `default:"" help:"Tigris token."`
	TigrisURL          string `default:"http://127.0.0.1:8081/" help:"Tigris URL."`
}

// init adds "tigris" handler flags when "ferretdb_tigris" build tag is provided.
func init() {
	cli.Plugins = kong.Plugins{&tigrisFlags}
}
