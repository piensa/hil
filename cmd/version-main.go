/*
 * Minio Client (C) 2014, 2015, 2016, 2017 Minio, Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/fatih/color"
	"github.com/minio/cli"
	"github.com/piensa/gnd/pkg/console"
	"github.com/piensa/gnd/pkg/probe"
)

// Print version.
var versionCmd = cli.Command{
	Name:   "version",
	Usage:  "Print version info.",
	Action: mainVersion,
	Before: setGlobalsFromContext,
	Flags: []cli.Flag{
		cli.BoolFlag{
			Name:  "quiet, q",
			Usage: "Suppress chatty console output.",
		},
		cli.BoolFlag{
			Name:  "json",
			Usage: "Enable JSON formatted output.",
		},
	},
	CustomHelpTemplate: `NAME:
  {{.HelpName}} - {{.Usage}}

USAGE:
  {{.HelpName}}{{if .VisibleFlags}} [FLAGS]{{end}}
{{if .VisibleFlags}}
FLAGS:
  {{range .VisibleFlags}}{{.}}
  {{end}}{{end}}
EXAMPLES:
   1. Prints the minio client version:
       $ {{.HelpName}}
`,
}

// Structured message depending on the type of console.
type versionMessage struct {
	Status  string `json:"status"`
	Version struct {
		Value  string `json:"value"`
		Format string `json:"format"`
	} `json:"version"`
	ReleaseTag string `json:"releaseTag"`
	CommitID   string `json:"commitID"`
}

// Colorized message for console printing.
func (v versionMessage) String() string {
	return console.Colorize("Version", fmt.Sprintf("Version: %s\n", v.Version.Value)) +
		console.Colorize("ReleaseTag", fmt.Sprintf("Release-tag: %s\n", v.ReleaseTag)) +
		console.Colorize("CommitID", fmt.Sprintf("Commit-id: %s", v.CommitID))
}

// JSON'ified message for scripting.
func (v versionMessage) JSON() string {
	v.Status = "success"
	msgBytes, e := json.Marshal(v)
	fatalIf(probe.NewError(e), "Unable to marshal into JSON.")
	return string(msgBytes)
}

func mainVersion(ctx *cli.Context) error {

	// Additional command speific theme customization.
	console.SetColor("Version", color.New(color.FgGreen, color.Bold))
	console.SetColor("ReleaseTag", color.New(color.FgGreen))
	console.SetColor("CommitID", color.New(color.FgGreen))

	verMsg := versionMessage{}
	verMsg.CommitID = CommitID
	verMsg.ReleaseTag = ReleaseTag
	verMsg.Version.Value = Version
	verMsg.Version.Format = "RFC3339"

	if !globalQuiet {
		printMsg(verMsg)
	}
	return nil
}
