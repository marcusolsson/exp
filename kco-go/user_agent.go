// Copyright 2016 Marcus Olsson
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

package kco

import (
	"fmt"
	"runtime"
	"strings"
)

func userAgent() string {
	fields := []string{
		newField("Library", "Klarna.ApiWrapper", "3.0.0"),
		newField("OS", runtime.GOOS, runtime.GOARCH),
		newField("Language", "go", runtime.Version()),
	}
	return strings.Join(fields, " ")
}

func newField(key, name, version string, opts ...string) string {
	result := fmt.Sprintf("%s/%s_%s", key, name, version)

	if len(opts) == 0 {
		return result
	}

	return result + " (" + strings.Join(opts, " ; ") + ")"
}
