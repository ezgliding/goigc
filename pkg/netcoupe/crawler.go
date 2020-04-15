// Copyright The ezgliding Authors.
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

package netcoupe

import (
	"errors"
	"fmt"
	"io"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/ezgliding/goigc/pkg/igc"
	"github.com/gocolly/colly"
	log "github.com/sirupsen/logrus"
)

// DailyUrlPattern is the main page to list netcoupe flights.
const DailyUrlPattern = "https://%v.netcoupe.net/Results/DailyResults.aspx"

// FlightBaseUrl is the base path to fetch flight details from a flight ID.
const FlightBaseUrlPattern = "https://%v.netcoupe.net/Results/FlightDetail.aspx?FlightID="

// TrackBaseUrlPattern is the base path to download the flight track from a track ID.
const TrackBaseUrlPattern = "https://%v.netcoupe.net/Download/DownloadIGC.aspx?FileID="

// This is a constant map.
var httpHeaders = map[string][]string{
	"Accept-Encoding":           []string{"gzip, deflate"},
	"Cache-Control":             []string{"max-age=0"},
	"Upgrade-Insecure-Requests": []string{"1"},
	"DNT":                       []string{"1"},
	"Origin":                    []string{"https://archive2019.netcoupe.net"},
	"User-Agent":                []string{"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Ubuntu Chromium/66.0.3359.181 Chrome/66.0.3359.181 Safari/537.36"},
	"Accept":                    []string{"text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8"},
	"Accept-Language":           []string{"en-US,en;q=0.9,de;q=0.8,fr;q=0.7,pt;q=0.6,es;q=0.5,it;q=0.4,ny;q=0.3"},
	"Connection":                []string{"keep-alive"},
	"Referer":                   []string{"https://archive2019.netcoupe.net/Results/DailyResults.aspx"}}

// Netcoupe implements a crawler for http://netcoupe.net.
type Netcoupe struct {
	collector *colly.Collector
	year      int
	baseUrl   string
}

func NewNetcoupeYear(year int) Netcoupe {
	n := Netcoupe{}
	n.year = year
	n.baseUrl = "www"
	n.collector = colly.NewCollector()
	n.collector.AllowURLRevisit = true
	n.collector.UserAgent = httpHeaders["User-Agent"][0]
	if year != 0 {
		n.baseUrl = fmt.Sprintf("archive%v", year)
	}
	return n
}

func NewNetcoupe() Netcoupe {
	return NewNetcoupeYear(0)
}

// Crawl checks for flights on netcoupe.net.
//
// It works by querying flights for specific days, making it easier to iterate
// through past data. The rules for flight submission are defined here:
// http://www.planeur.net/_download/netcoupe/2018_np4.2.pdf
// """
// Chaque performance doit être enregistrée dans un délai de 15 jours sur le
// site de la NetCoupe (www.netcoupe.net) par le commandant de bord ou par le
// Responsable de la NetCoupe de l’association avec l'accord du pilote.
// """
// Which means that it's only worth to crawl for new flights back to 2 weeks max.
func (n Netcoupe) Crawl(start time.Time, end time.Time) ([]igc.Flight, error) {
	var flights []igc.Flight

	// Do not allow start > end
	if end.Before(start) {
		return nil, errors.New("Invalid start end date pair")
	}

	r, _ := regexp.Compile(`.*DisplayFlightDetail\('(?P<ID>[0-9]+)'\).*`)

	c := n.newCollector()
	c.OnRequest(func(r *colly.Request) {
		log.WithFields(log.Fields{
			"url":     r.URL.String(),
			"headers": r.Headers}).Trace("Visiting flight list")
	})
	c.OnError(func(r *colly.Response, err error) {
		log.WithFields(log.Fields{
			"response": r,
			"error":    err}).Error("Failed to visit url")
		// Poor handling of http 420
		if strings.Contains(err.Error(), "ENHANCE_YOUR_CALM") {
			time.Sleep(10 * time.Second)
		}
	})

	d := n.newCollector()
	d.OnRequest(func(r *colly.Request) {
		log.WithFields(log.Fields{
			"url":     r.URL.String(),
			"headers": r.Headers}).Trace("Visiting flight details")
	})

	c.OnHTML("table tr td:nth-child(4) a[href]", func(e *colly.HTMLElement) {
		id := r.FindStringSubmatch(e.Attr("href"))
		log.WithFields(log.Fields{
			"flight_id": id[1]}).Trace("Scheduling visit to flight details")
		if len(id) == 2 {
			// TODO(rochaporto): handle error
			_ = d.Visit(fmt.Sprintf("%v%v", n.flightBaseUrl(), id[1]))
		}

	})

	d.OnHTML("div center table[width]", func(e *colly.HTMLElement) {
		f := igc.Flight{}
		f.URL = e.Request.URL.String()
		f.ID = e.Request.URL.Query()["FlightID"][0]
		f.Pilot = e.ChildText("tbody tr:nth-child(3) td:nth-child(2) a")
		f.Club = e.ChildText("tbody tr:nth-child(5) td:nth-child(2) a")
		f.Date, _ = time.Parse("02/01/2006", e.ChildText("tbody tr:nth-child(8) td:nth-child(2) div"))
		f.Takeoff = e.ChildText("tbody tr:nth-child(9) td:nth-child(2) div")
		f.Region = e.ChildText("tbody tr:nth-child(10) td:nth-child(2) div")
		f.Country = e.ChildText("tbody tr:nth-child(11) td:nth-child(2) div")
		f.Distance = parseFloat(e.ChildText("tbody tr:nth-child(12) td:nth-child(2) div"))
		f.Points = parseFloat(e.ChildText("tbody tr:nth-child(13) td:nth-child(2) div"))
		f.Glider = e.ChildText("tbody tr:nth-child(14) td:nth-child(2) div table tbody tr td")

		i := 0
		if strings.Contains(e.ChildText("tbody tr:nth-child(15) td:nth-child(1) div"), "Comp") {
			f.CompetitionURL = e.ChildText("tbody tr:nth-child(15) td:nth-child(2) div")
			i = 1
		}
		f.Type = e.ChildText(fmt.Sprintf("tbody tr:nth-child(%v) td:nth-child(2) div", 15+i))
		trackUrl, err := url.Parse(
			e.ChildAttr(fmt.Sprintf("tbody tr:nth-child(%v) td:nth-child(2) div a", 16+i), "href"))
		if err == nil && trackUrl.RawQuery != "" {
			f.TrackID = trackUrl.Query()["FileID"][0]
			f.TrackURL = fmt.Sprintf("%v%v", n.TrackBaseUrl(), f.TrackID)
		}
		f.Speed = parseFloat(e.ChildText(fmt.Sprintf("tbody tr:nth-child(%v) td:nth-child(2) div", 17+i)))
		f.Comments = e.ChildText(fmt.Sprintf("tbody tr:nth-child(%v) td:nth-child(2) div", 23+i))

		flights = append(flights, f)
	})

	current := time.Date(start.Year(), start.Month(), start.Day(), 12, 0, 0, 0, time.UTC)
	end = time.Date(end.Year(), end.Month(), end.Day(), 12, 0, 0, 0, time.UTC)
	for ; end.After(current.AddDate(0, 0, -1)); current = current.AddDate(0, 0, 1) {
		data := n.sessionHeaders(c)
		data["ddlDisplayRange"] = "0"
		data["ddlDisplayDate"] = current.Format("02/01/2006")
		data["rbgDisplayMode"] = "rbDisplayByDate"
		tmp := n.newCollector()
		tmp.OnHTML("input", func(e *colly.HTMLElement) {
			switch e.Attr("name") {
			case "__EVENTVALIDATION":
				data["__EVENTVALIDATION"] = e.Attr("value")
			case "__VIEWSTATE":
				data["__VIEWSTATE"] = e.Attr("value")
			case "__VIEWSTATEGENERATOR":
				data["__VIEWSTATEGENERATOR"] = e.Attr("value")
			}
		})
		n.post(tmp, n.dailyUrl(), data)
		// Set header for single page results (by default it pages)
		data["__EVENTTARGET"] = "dgDailyResults$ctl01$ctl01"
		n.post(c, n.dailyUrl(), data)
	}

	log.WithFields(log.Fields{
		"start":       start,
		"end":         end,
		"flights":     flights,
		"num_flights": len(flights),
	}).Trace("Finishing crawling flights")
	return flights, nil
}

func (n Netcoupe) Get(url string) ([]byte, error) {
	var result []byte

	t := n.newCollector()
	t.OnRequest(func(r *colly.Request) {
		log.WithFields(log.Fields{
			"url":     r.URL.String(),
			"headers": r.Headers}).Trace("Visiting flight track")
	})
	t.OnResponse(func(r *colly.Response) {
		result = r.Body
	})
	err := t.Visit(url)

	return result, err
}

func (n Netcoupe) newCollector() *colly.Collector {
	return n.collector.Clone()
}

func (n Netcoupe) sessionHeaders(c *colly.Collector) map[string]string {
	headers := map[string]string{
		"__EVENTARGUMENT": "",
		"__LASTFOCUS":     "",
		"__EVENTTARGET":   "ddlDisplayDate",
	}

	t := c.Clone()
	t.OnRequest(func(r *colly.Request) {
		log.WithFields(log.Fields{
			"url":     r.URL.String(),
			"headers": r.Headers}).Trace("Visiting for session data collection")

	})
	t.OnHTML("input", func(e *colly.HTMLElement) {
		switch e.Attr("name") {
		case "__EVENTVALIDATION":
			headers["__EVENTVALIDATION"] = e.Attr("value")
		case "__VIEWSTATE":
			headers["__VIEWSTATE"] = e.Attr("value")
		case "__VIEWSTATEGENERATOR":
			headers["__VIEWSTATEGENERATOR"] = e.Attr("value")
		}
	})
	// TODO(rochaporto): handle error
	_ = t.Request("GET", n.dailyUrl(), nil, nil, httpHeaders)

	return headers
}

func (n Netcoupe) post(c *colly.Collector, url string, data map[string]string) {
	cookies := c.Cookies(url)
	// TODO(rochaporto): handle error
	_ = c.SetCookies(url, cookies)
	dur, _ := time.ParseDuration("1m")
	c.SetRequestTimeout(dur)
	log.WithFields(log.Fields{
		"url": url}).Trace("POST request")
	// TODO(rochaporto): handle error
	_ = c.Request("POST", url, createFormReader(data), nil, httpHeaders)
}

func (n Netcoupe) dailyUrl() string {
	return fmt.Sprintf(DailyUrlPattern, n.baseUrl)
}

func (n Netcoupe) flightBaseUrl() string {
	return fmt.Sprintf(FlightBaseUrlPattern, n.baseUrl)
}

func (n Netcoupe) TrackBaseUrl() string {
	return fmt.Sprintf(TrackBaseUrlPattern, n.baseUrl)
}

func parseFloat(s string) float64 {
	rs := strings.Replace(strings.TrimSpace(strings.Split(s, " ")[0]), ",", ".", -1)
	r, _ := strconv.ParseFloat(rs, 32)
	return r
}

func createFormReader(data map[string]string) io.Reader {
	form := url.Values{}
	for k, v := range data {
		form.Add(k, v)
	}
	return strings.NewReader(form.Encode())
}
