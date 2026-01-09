package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/nyaruka/phonenumbers"

	"cyph3r/internal/netcheck"
	"cyph3r/internal/output"
	"cyph3r/internal/version"
)

func printHelp() {
	fmt.Println(`
CYPH3R: Network Diagnostics Utility

Usage:
  cyph3r [options]

Options:
  --target        Target host or IP (default: localhost)
  --port          Port number (default: 80)
  --proto         Protocol: tcp | udp | http | https (default: tcp)
  --geoip         Lookup GeoIP and ASN info for the target
  --phone         Show phone number info (region, type, validity)
  --json          Print output in JSON
  --monitor       Continuously monitor the target
  --interval      Check interval in seconds (default: 5)
  --version       Show program version info
  --portscan      Scan ports on IP/localhost (default range 1-1024)
  --scanstart     Port scan range start (default: 1)
  --scanend       Port scan range end (default: 1024)
  -h, --help      Show this help and exit

Examples:
  cyph3r --target example.com --geoip
  cyph3r --target 8.8.8.8 --portscan
  cyph3r --phone +14155552671
  cyph3r --target localhost --proto tcp
`)
}

func containsHelpFlag() bool {
	for _, arg := range os.Args {
		if arg == "--help" || arg == "-h" {
			return true
		}
	}
	return false
}

// --- GeoIP advanced lookup ---
type GeoResult struct {
	City     string  `json:"city"`
	Region   string  `json:"regionName"`
	Country  string  `json:"country"`
	Lat      float64 `json:"lat"`
	Lon      float64 `json:"lon"`
	ASN      string  `json:"as"`
	Org      string  `json:"org"`
	Hostname string  `json:"reverse"`
}

func AdvancedGeoIPLookup(target string) (*GeoResult, error) {
	ip := target
	ips, err := net.LookupIP(target)
	if err == nil && len(ips) > 0 {
		ip = ips[0].String()
	}
	resp, err := http.Get("http://ip-api.com/json/" + ip + "?fields=status,message,country,regionName,city,lat,lon,org,as,reverse")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var r struct {
		Status  string  `json:"status"`
		Message string  `json:"message"`
		Country string  `json:"country"`
		Region  string  `json:"regionName"`
		City    string  `json:"city"`
		Lat     float64 `json:"lat"`
		Lon     float64 `json:"lon"`
		ASN     string  `json:"as"`
		Org     string  `json:"org"`
		Host    string  `json:"reverse"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
		return nil, err
	}
	if r.Status != "success" {
		return nil, fmt.Errorf("lookup failed: %s", r.Message)
	}
	return &GeoResult{
		City:     r.City,
		Region:   r.Region,
		Country:  r.Country,
		Lat:      r.Lat,
		Lon:      r.Lon,
		ASN:      r.ASN,
		Org:      r.Org,
		Hostname: r.Host,
	}, nil
}

// --- Phone parsing lookup ---
func LocalPhoneLookup(number string) {
	num, err := phonenumbers.Parse(number, "")
	if err != nil {
		fmt.Println("Error parsing number:", err)
		return
	}

	region := phonenumbers.GetRegionCodeForNumber(num)
	numType := phonenumbers.GetNumberType(num)
	typeStr := map[phonenumbers.PhoneNumberType]string{
		phonenumbers.FIXED_LINE:           "Fixed line",
		phonenumbers.MOBILE:               "Mobile",
		phonenumbers.FIXED_LINE_OR_MOBILE: "Fixed line or mobile",
		phonenumbers.TOLL_FREE:            "Toll free",
		phonenumbers.PREMIUM_RATE:         "Premium rate",
		phonenumbers.SHARED_COST:          "Shared cost",
		phonenumbers.VOIP:                 "VOIP",
		phonenumbers.PERSONAL_NUMBER:      "Personal number",
		phonenumbers.PAGER:                "Pager",
		phonenumbers.UAN:                  "UAN",
		phonenumbers.VOICEMAIL:            "Voicemail",
		phonenumbers.UNKNOWN:              "Unknown",
	}[numType]

	fmt.Print("\nPhone Number Info\n")
	fmt.Printf("  Raw input:  %s\n", number)
	fmt.Printf("  E.164:      %s\n", phonenumbers.Format(num, phonenumbers.E164))
	fmt.Printf("  Region:     %s\n", region)
	fmt.Printf("  Type:       %s\n", typeStr)
	fmt.Printf("  Valid:      %v\n", phonenumbers.IsValidNumber(num))
	fmt.Println("")
}

// --- Simple port scanner ---
func isIPv4(address string) bool {
	ip := net.ParseIP(address)
	return ip != nil && strings.Count(address, ":") < 1
}

func scanPorts(host string, start, end int) []int {
	open := []int{}
	timeout := 300 * time.Millisecond
	for port := start; port  0
			result["status"] = code
			result["latency_ms"] = latency
			result["up"] = up

		default:
			fmt.Println("Unknown protocol:", *proto)
			return
		}

		// --- State tracking ---
		if wasUp == nil {
			wasUp = new(bool)
			*wasUp = up

			if up {
				output.Up("Target is UP")
			} else {
				downSince = time.Now()
				output.Down("Target is DOWN")
			}
		} else {
			// UP → DOWN
			if *wasUp && !up {
				downSince = time.Now()
				output.Down("Target went DOWN")
			}

			// DOWN → UP
			if !*wasUp && up {
				downtime := time.Since(downSince).Round(time.Second)
				output.Up(fmt.Sprintf("Target is UP again (downtime: %s)", downtime))
				result["downtime"] = downtime.String()
			}
		}

		*wasUp = up

		if *jsonOut {
			output.PrintJSON(result)
		}

		if !*monitor {
			break
		}

		time.Sleep(time.Duration(*interval) * time.Second)
	}
}
