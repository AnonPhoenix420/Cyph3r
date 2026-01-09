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
)

type GeoResult struct {
	City     string  `json:"city"`
	Region   string  `json:"region"`
	Country  string  `json:"country"`
	Lat      float64 `json:"lat"`
	Lon      float64 `json:"lon"`
	ASN      string  `json:"as"`
	Org      string  `json:"org"`
	Hostname string  `json:"reverse"`
}

func geoIPLookup(target string) (*GeoResult, error) {
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
	var data struct {
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
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, err
	}
	if data.Status != "success" {
		return nil, fmt.Errorf("lookup failed: %s", data.Message)
	}
	return &GeoResult{
		City:     data.City,
		Region:   data.Region,
		Country:  data.Country,
		Lat:      data.Lat,
		Lon:      data.Lon,
		ASN:      data.ASN,
		Org:      data.Org,
		Hostname: data.Host,
	}, nil
}

func lookupPhone(number string) {
	num, err := phonenumbers.Parse(number, "")
	if err != nil {
		fmt.Println("Error parsing number:", err)
		return
	}
	fmt.Println("Raw:", number)
	fmt.Println("Valid:", phonenumbers.IsValidNumber(num))
	fmt.Println("Region:", phonenumbers.GetRegionCodeForNumber(num))
	fmt.Println("Type:", phonenumbers.GetNumberType(num))
}

func scanPorts(host string, start, end int) []int {
	open := []int{}
	timeout := 300 * time.Millisecond
	for port := start; port = 200 && code < 400
			details = fmt.Sprintf("Status code: %d, Error: %v", code, err)
		case "https":
			code, err := checkHTTP(*target, true)
			up = err == nil && code >= 200 && code < 400
			details = fmt.Sprintf("Status code: %d, Error: %v", code, err)
		case "dns":
			up = checkDNS(*target)
		default:
			fmt.Println("Unknown protocol:", *proto)
			return
		}

		result["up"] = up
		if details != "" { result["details"] = details }
		if *jsonOut {
			printJSON(result)
		} else {
			fmt.Println(result)
		}
		if !*monitor {
			break
		}
		time.Sleep(time.Duration(*interval) * time.Second)
	}
}
