/*
# Copyright 2022 Google Inc.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
*/

package main

import (
	"fmt"
	"strconv"
	"time"
)

func main() {
	// Setting default values
	if project_sync == "" {
		project_sync = "enable"
	}
	if project_sync_interval == "" {
		project_sync_interval = "30"
	}

	if project_sync == "enable" {
		for {
			interval, err := strconv.Atoi(project_sync_interval)
			checkError(err)
			fmt.Println("Executing projects sync job at: " + time.Now().String())
			projectIDConverter()
			time.Sleep(time.Duration(interval) * time.Minute)
		}
	}
}
