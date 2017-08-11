// Copyright ©2017 The ezgliding Authors.
//
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
//

package igc

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const (
	// TimeFormat is the golang time.Parse format for IGC time.
	TimeFormat = "150405"
	// DateFormat is the golang time.Parse format for IGC time.
	DateFormat = "020106"
)

// ParseLocation returns a Track object corresponding to the given file.
//
// It calls Parse internatlly, so the file content should be in IGC format.
func ParseLocation(location string) (Track, error) {
	var content []byte
	resp, err := http.Get(location)
	// case http
	if err == nil {
		defer resp.Body.Close()
		content, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			return Track{}, err
		}
	} else { // case file
		resp, err := ioutil.ReadFile(location)
		if err != nil {
			return Track{}, err
		}
		content = resp
	}

	return Parse(string(content))
}

// Parse returns a Track object corresponding to the given content.
//
// The value of content should be a text string with all the flight data
// in the IGC format.
func Parse(content string) (Track, error) {
	f := NewTrack()
	var err error
	p := parser{}
	lines := strings.Split(content, "\n")
	for i := range lines {
		line := strings.TrimSpace(lines[i])
		// ignore empty lines
		if len(strings.Trim(line, " ")) < 1 {
			continue
		}
		switch line[0] {
		case 'A':
			err = p.parseA(line, &f)
		case 'B':
			err = p.parseB(line, &f)
		case 'C':
			if !p.taskDone {
				err = p.parseC(lines[i:], &f)
			}
		case 'D':
			err = p.parseD(line, &f)
		case 'E':
			err = p.parseE(line, &f)
		case 'F':
			err = p.parseF(line, &f)
		case 'G':
			err = p.parseG(line, &f)
		case 'H':
			err = p.parseH(line, &f)
		case 'I':
			err = p.parseI(line)
		case 'J':
			err = p.parseJ(line)
		case 'K':
			err = p.parseK(line, &f)
		case 'L':
			err = p.parseL(line, &f)
		default:
			err = fmt.Errorf("invalid record :: %v", line)
		}
		if err != nil {
			return f, err
		}
	}

	return f, nil
}

type field struct {
	start int64
	end   int64
	tlc   string
}

type parser struct {
	IFields  []field
	JFields  []field
	taskDone bool
	numSat   int
}

func (p *parser) parseA(line string, f *Track) error {
	if len(line) < 7 {
		return fmt.Errorf("line too short :: %v", line)
	}
	f.Manufacturer = line[1:4]
	f.UniqueID = line[4:7]
	f.AdditionalData = line[7:]
	return nil
}

func (p *parser) parseB(line string, f *Track) error {
	if len(line) < 35 {
		return fmt.Errorf("line too short :: %v", line)
	}
	pt := NewPointFromDMD(
		line[7:15], line[15:24])

	var err error
	pt.Time, err = time.Parse(TimeFormat, line[1:7])
	if err != nil {
		return err
	}
	if line[24] == 'A' || line[24] == 'V' {
		pt.FixValidity = line[24]
	} else {
		return fmt.Errorf("invalid fix validity :: %v", line[24])
	}
	pt.PressureAltitude, err = strconv.ParseInt(line[25:30], 10, 64)
	if err != nil {
		return err
	}
	pt.GNSSAltitude, err = strconv.ParseInt(line[30:35], 10, 64)
	if err != nil {
		return err
	}
	for _, f := range p.IFields {
		pt.IData[f.tlc] = line[f.start-1 : f.end]
	}
	pt.NumSatellites = p.numSat
	f.Points = append(f.Points, pt)
	return nil
}

func (p *parser) parseC(lines []string, f *Track) error {
	line := lines[0]
	if len(line) < 25 {
		return fmt.Errorf("wrong line size :: %v", line)
	}
	var err error
	var nTP int
	if nTP, err = strconv.Atoi(line[23:25]); err != nil {
		return fmt.Errorf("invalid number of turnpoints :: %v", line)
	}
	if len(lines) < 5+nTP {
		return fmt.Errorf("invalid number of C record lines :: %v", lines)
	}
	if f.Task.DeclarationDate, err = time.Parse(DateFormat+TimeFormat, lines[0][1:13]); err != nil {
		f.Task.DeclarationDate = time.Time{}
	}
	if f.Task.Date, err = time.Parse(DateFormat, lines[0][13:19]); err != nil {
		f.Task.Date = time.Time{}
	}
	if f.Task.Number, err = strconv.Atoi(line[19:23]); err != nil {
		return err
	}
	f.Task.Description = line[25:]
	if f.Task.Takeoff, err = p.taskPoint(lines[1]); err != nil {
		return err
	}
	if f.Task.Start, err = p.taskPoint(lines[2]); err != nil {
		return err
	}
	for i := 0; i < nTP; i++ {
		var tp Point
		if tp, err = p.taskPoint(lines[3+i]); err != nil {
			return err
		}
		f.Task.Turnpoints = append(f.Task.Turnpoints, tp)
	}
	if f.Task.Finish, err = p.taskPoint(lines[3+nTP]); err != nil {
		return err
	}
	if f.Task.Landing, err = p.taskPoint(lines[4+nTP]); err != nil {
		return err
	}
	p.taskDone = true
	return nil
}

func (p *parser) taskPoint(line string) (Point, error) {
	if len(line) < 18 {
		return Point{}, fmt.Errorf("line too short :: %v", line)
	}
	pt := NewPointFromDMD(
		line[1:9], line[9:18])
	pt.Description = line[18:]
	return pt, nil
}

func (p *parser) parseD(line string, f *Track) error {
	if len(line) < 6 {
		return fmt.Errorf("line too short :: %v", line)
	}
	if line[1] == '2' {
		f.DGPSStationID = line[2:6]
	}
	return nil
}

func (p *parser) parseE(line string, f *Track) error {
	if len(line) < 10 {
		return fmt.Errorf("line too short :: %v", line)
	}
	t, err := time.Parse(TimeFormat, line[1:7])
	if err != nil {
		return err
	}
	f.Events = append(f.Events, Event{Time: t, Type: line[7:10], Data: line[10:]})
	return nil
}

func (p *parser) parseF(line string, f *Track) error {
	if len(line) < 7 {
		return fmt.Errorf("line too short :: %v", line)
	}
	t, err := time.Parse(TimeFormat, line[1:7])
	if err != nil {
		return err
	}
	ids := []string{}
	for i := 7; i < len(line)-1; i = i + 2 {
		ids = append(ids, line[i:i+2])
	}
	f.Satellites = append(f.Satellites, Satellite{Time: t, Ids: ids})
	p.numSat = len(ids)
	return nil
}

func (p *parser) parseG(line string, f *Track) error {
	f.Signature = f.Signature + line[1:]
	return nil
}

func (p *parser) parseH(line string, f *Track) error {
	var err error
	if len(line) < 5 {
		return fmt.Errorf("line too short :: %v", line)
	}

	switch line[2:5] {
	case "DTE":
		if len(line) < 11 {
			return fmt.Errorf("line too short :: %v", line)
		}
		f.Date, err = time.Parse(DateFormat, line[5:11])
	case "FXA":
		if len(line) < 8 {
			return fmt.Errorf("line too short :: %v", line)
		}
		f.FixAccuracy, err = strconv.ParseInt(line[5:8], 10, 64)
	case "PLT":
		f.Pilot = stripUpTo(line[5:], ":")
	case "CM2":
		f.Crew = stripUpTo(line[5:], ":")
	case "GTY":
		f.GliderType = stripUpTo(line[5:], ":")
	case "GID":
		f.GliderID = stripUpTo(line[5:], ":")
	case "DTM":
		if len(line) < 8 {
			return fmt.Errorf("line too short :: %v", line)
		}
		f.GPSDatum = stripUpTo(line[5:], ":")
	case "RFW":
		f.FirmwareVersion = stripUpTo(line[5:], ":")
	case "RHW":
		f.HardwareVersion = stripUpTo(line[5:], ":")
	case "FTY":
		f.FlightRecorder = stripUpTo(line[5:], ":")
	case "GPS":
		f.GPS = line[5:]
	case "PRS":
		f.PressureSensor = stripUpTo(line[5:], ":")
	case "CID":
		f.CompetitionID = stripUpTo(line[5:], ":")
	case "CCL":
		f.CompetitionClass = stripUpTo(line[5:], ":")
	case "TZN":
		z, err := strconv.ParseFloat(stripUpTo(line[5:], ":"), 64)
		if err != nil {
			return err
		}
		f.Timezone = int(z)
	default:
		err = fmt.Errorf("unknown record :: %v", line)
	}

	return err
}

func (p *parser) parseI(line string) error {
	if len(line) < 3 {
		return fmt.Errorf("line too short :: %v", line)
	}
	n, err := strconv.ParseInt(line[1:3], 10, 0)
	if err != nil {
		return fmt.Errorf("invalid number of I fields :: %v", line)
	}
	if len(line) != int(n*7+3) {
		return fmt.Errorf("wrong line size :: %v", line)
	}
	for i := 0; i < int(n); i++ {
		s := i*7 + 3
		start, _ := strconv.ParseInt(line[s:s+2], 10, 0)
		end, _ := strconv.ParseInt(line[s+2:s+4], 10, 0)
		tlc := line[s+4 : s+7]
		p.IFields = append(p.IFields, field{start: start, end: end, tlc: tlc})
	}
	return nil
}

func (p *parser) parseJ(line string) error {
	if len(line) < 3 {
		return fmt.Errorf("line too short :: %v", line)
	}
	n, err := strconv.ParseInt(line[1:3], 10, 0)
	if err != nil {
		return fmt.Errorf("invalid number of J fields :: %v", line)
	}
	if len(line) != int(n*7+3) {
		return fmt.Errorf("wrong line size :: %v", line)
	}
	for i := 0; i < int(n); i++ {
		s := i*7 + 3
		start, _ := strconv.ParseInt(line[s:s+2], 10, 0)
		end, _ := strconv.ParseInt(line[s+2:s+4], 10, 0)
		tlc := line[s+4 : s+7]
		p.JFields = append(p.JFields, field{start: start, end: end, tlc: tlc})
	}
	return nil
}

func (p *parser) parseK(line string, f *Track) error {
	if len(line) < 7 {
		return fmt.Errorf("line too short :: %v", line)
	}
	t, err := time.Parse(TimeFormat, line[1:7])
	if err != nil {
		return err
	}
	fields := make(map[string]string)
	for _, f := range p.JFields {
		fields[f.tlc] = line[f.start-1 : f.end]
	}
	f.K = append(f.K, K{Time: t, Fields: fields})
	return nil
}

func (p *parser) parseL(line string, f *Track) error {
	if len(line) < 4 {
		return fmt.Errorf("line too short :: %v", line)
	}
	f.Logbook = append(f.Logbook, LogEntry{Type: line[1:4], Text: line[4:]})
	return nil
}

func stripUpTo(s string, sep string) string {
	i := strings.Index(s, sep)
	if i == -1 {
		return s
	}
	return s[i+1:]
}
