/*
 * Minio Client (C) 2016 Minio, Inc.
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
	"os"
	"sync"
	"time"

	"github.com/minio/minio/pkg/quick"
	"github.com/piensa/geo/pkg/probe"
)

/////////////////// Session V6 ///////////////////
// sessionV6Header for resumable sessions.
type sessionV6Header struct {
	Version            string            `json:"version"`
	When               time.Time         `json:"time"`
	RootPath           string            `json:"workingFolder"`
	GlobalBoolFlags    map[string]bool   `json:"globalBoolFlags"`
	GlobalIntFlags     map[string]int    `json:"globalIntFlags"`
	GlobalStringFlags  map[string]string `json:"globalStringFlags"`
	CommandType        string            `json:"commandType"`
	CommandArgs        []string          `json:"cmdArgs"`
	CommandBoolFlags   map[string]bool   `json:"cmdBoolFlags"`
	CommandIntFlags    map[string]int    `json:"cmdIntFlags"`
	CommandStringFlags map[string]string `json:"cmdStringFlags"`
	LastCopied         string            `json:"lastCopied"`
	TotalBytes         int64             `json:"totalBytes"`
	TotalObjects       int               `json:"totalObjects"`
}

func loadSessionV6Header(sid string) (*sessionV6Header, *probe.Error) {
	if !isSessionDirExists() {
		return nil, errInvalidArgument().Trace()
	}

	sessionFile, err := getSessionFile(sid)
	if err != nil {
		return nil, err.Trace(sid)
	}

	if _, e := os.Stat(sessionFile); e != nil {
		return nil, probe.NewError(e)
	}

	sV6Header := &sessionV6Header{}
	sV6Header.Version = "6"
	qs, e := quick.New(sV6Header)
	if e != nil {
		return nil, probe.NewError(e).Trace(sid, sV6Header.Version)
	}
	e = qs.Load(sessionFile)
	if e != nil {
		return nil, probe.NewError(e).Trace(sid, sV6Header.Version)
	}

	sV6Header = qs.Data().(*sessionV6Header)
	return sV6Header, nil
}

/////////////////// Session V7 ///////////////////
// RESERVED FOR FUTURE

// sessionV7Header for resumable sessions.
type sessionV7Header struct {
	Version            string            `json:"version"`
	When               time.Time         `json:"time"`
	RootPath           string            `json:"workingFolder"`
	GlobalBoolFlags    map[string]bool   `json:"globalBoolFlags"`
	GlobalIntFlags     map[string]int    `json:"globalIntFlags"`
	GlobalStringFlags  map[string]string `json:"globalStringFlags"`
	CommandType        string            `json:"commandType"`
	CommandArgs        []string          `json:"cmdArgs"`
	CommandBoolFlags   map[string]bool   `json:"cmdBoolFlags"`
	CommandIntFlags    map[string]int    `json:"cmdIntFlags"`
	CommandStringFlags map[string]string `json:"cmdStringFlags"`
	LastCopied         string            `json:"lastCopied"`
	LastRemoved        string            `json:"lastRemoved"`
	TotalBytes         int64             `json:"totalBytes"`
	TotalObjects       int               `json:"totalObjects"`
}

// sessionV7 resumable session container.
type sessionV7 struct {
	Header    *sessionV7Header
	SessionID string
	mutex     *sync.Mutex
	DataFP    *sessionDataFP
	sigCh     bool
}

// loadSessionV7 - reads session file if exists and re-initiates internal variables
func loadSessionV7(sid string) (*sessionV7, *probe.Error) {
	if !isSessionDirExists() {
		return nil, errInvalidArgument().Trace()
	}
	sessionFile, err := getSessionFile(sid)
	if err != nil {
		return nil, err.Trace(sid)
	}

	if _, e := os.Stat(sessionFile); e != nil {
		return nil, probe.NewError(e)
	}

	s := &sessionV7{}
	s.Header = &sessionV7Header{}
	s.SessionID = sid
	s.Header.Version = "7"
	qs, e := quick.New(s.Header)
	if e != nil {
		return nil, probe.NewError(e).Trace(sid, s.Header.Version)
	}
	e = qs.Load(sessionFile)
	if e != nil {
		return nil, probe.NewError(e).Trace(sid, s.Header.Version)
	}

	s.mutex = new(sync.Mutex)
	s.Header = qs.Data().(*sessionV7Header)

	sessionDataFile, err := getSessionDataFile(s.SessionID)
	if err != nil {
		return nil, err.Trace(sid, s.Header.Version)
	}

	dataFile, e := os.Open(sessionDataFile)
	if e != nil {
		return nil, probe.NewError(e)
	}
	s.DataFP = &sessionDataFP{false, dataFile}

	return s, nil
}
