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
	for port := start; port <= end; port++ {
		address := fmt.Sprintf("%s:%d", host, port)
		conn, err := net.DialTimeout("tcp", address, timeout)
		if err == nil {
			open = append(open, port)
			conn.Close()
		}
	}
	return open
}

func checkTCP(host string, port int) bool {
	address := fmt.Sprintf("%s:%d", host, port)
	conn, err := net.DialTimeout("tcp", address, 1*time.Second)
	if err != nil {
		return false
	}
	conn.Close()
	return true
}

func checkUDP(host string, port int) bool {
	address := fmt.Sprintf("%s:%d", host, port)
	conn, err := net.DialTimeout("udp", address, 1*time.Second)
	if err != nil {
		return false
	}
	defer conn.Close()
	// send/recv dummy packet (will usually timeout, but port open check works for many cases)
	_, err = conn.Write([]byte{0x0})
	return err == nil
}

func checkHTTP(url string, useHTTPS bool) (int, error) {
	if useHTTPS && !strings.HasPrefix(url, "https://") {
		url = "https://" + url
	}
	if !useHTTPS && !strings.HasPrefix(url, "http://") {
		url = "http://" + url
	}
	client := http.Client{Timeout: 2 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()
	return resp.StatusCode, nil
}

func checkDNS(host string) bool {
	_, err := net.LookupHost(host)
	return err == nil
}

func printJSON(data interface{}) {
	b, _ := json.MarshalIndent(data, "", "  ")
	fmt.Println(string(b))
}

func main() {
	// Flags
	target := flag.String("target", "localhost", "Target host or IP")
	port := flag.Int("port", 80, "Port number")
	proto := flag.String("proto", "tcp", "Protocol: tcp/udp/http/https/dns")
	geoip := flag.Bool("geoip", false, "GeoIP lookup")
	phone := flag.String("phone", "", "Phone number info")
	jsonOut := flag.Bool("json", false, "JSON output")
	monitor := flag.Bool("monitor", false, "Continuous monitor")
	interval := flag.Int("interval", 5, "Monitor interval (seconds)")
	portscan := flag.Bool("portscan", false, "Scan ports")
	scanstart := flag.Int("scanstart", 1, "Port scan range start")
	scanend := flag.Int("scanend", 1024, "Port scan range end")
	flag.Parse()

	if *geoip {
		res, err := geoIPLookup(*target)
		if err != nil {
			fmt.Println("GeoIP error:", err)
			return
		}
		if *jsonOut {
			printJSON(res)
		} else {
			fmt.Printf("GeoIP: %+v\n", res)
		}
		return
	}
	if *phone != "" {
		lookupPhone(*phone)
		return
	}
	if *portscan {
		ports := scanPorts(*target, *scanstart, *scanend)
		if *jsonOut {
			printJSON(map[string]interface{}{"host": *target, "open_ports": ports})
		} else {
			fmt.Println("Open ports:", ports)
		}
		return
	}

	for {
		result := make(map[string]interface{})
		result["protocol"] = *proto
		result["target"] = *target
		result["port"] = *port
		result["time"] = time.Now().Format(time.RFC3339)
		up := false
		var details string

		switch *proto {
		case "tcp":
			up = checkTCP(*target, *port)
		case "udp":
			up = checkUDP(*target, *port)
		case "http":
			code, err := checkHTTP(*target, false)
			up = err == nil && code >= 200 && code < 400
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
