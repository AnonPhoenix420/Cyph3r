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

	"cyph3r/output"
	"github.com/nyaruka/phonenumbers"
)

/*
====================
ENTRYPOINT
====================
*/
func main() {
	// ===== Banner =====
	output.Banner()

	// ===== Color sanity check =====
	fmt.Println(output.PinkBoldText("✔ Port Scan Color OK"))
	fmt.Println(output.YellowBoldText("✔ GeoIP Color OK"))
	fmt.Println(output.BlueText("✔ Status Color OK"))
	fmt.Println()

	// ===== Run application =====
	run()
}

/*
====================
APPLICATION LOGIC
====================
*/
func run() {
	target := flag.String("target", "localhost", "Target host or IP")
	port := flag.Int("port", 80, "Port number")
	proto := flag.String("proto", "tcp", "Protocol: tcp|udp|http|https|dns")
	geoip := flag.Bool("geoip", false, "GeoIP lookup")
	phone := flag.String("phone", "", "Phone number for info")
	jsonOut := flag.Bool("json", false, "JSON output")
	monitor := flag.Bool("monitor", false, "Continuous monitor")
	interval := flag.Int("interval", 5, "Monitor interval (seconds)")
	portscan := flag.Bool("portscan", false, "Perform port scan")
	scanstart := flag.Int("scanstart", 1, "Port scan start")
	scanend := flag.Int("scanend", 1024, "Port scan end")
	flag.Parse()

	if *geoip {
		res, err := geoIPLookup(*target)
		if err != nil {
			fmt.Println(output.RedText("GeoIP lookup failed:"), err)
			os.Exit(1)
		}

		if *jsonOut {
			printJSON(res)
		} else {
			fmt.Println(output.GeoLabel("Country:"), res.Country)
			fmt.Println(output.GeoLabel("Region:"), res.Region)
			fmt.Println(output.GeoLabel("City:"), res.City)
			fmt.Println(output.GeoLabel("ASN:"), res.ASN)
			fmt.Println(output.GeoLabel("Org:"), res.Org)
			fmt.Println(output.GeoLabel("Hostname:"), res.Hostname)
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
			printJSON(map[string]interface{}{"open_ports": ports})
		} else {
			for _, p := range ports {
				fmt.Println(output.PortLabel(fmt.Sprintf("Port %d/tcp OPEN", p)))
			}
		}
		return
	}

	for {
		up := false
		details := ""

		switch strings.ToLower(*proto) {
		case "tcp":
			addr := fmt.Sprintf("%s:%d", *target, *port)
			conn, err := net.DialTimeout("tcp", addr, 2*time.Second)
			up = err == nil
			if up {
				conn.Close()
			}
			details = fmt.Sprintf("err=%v", err)

		case "udp":
			up = checkUDP(*target, *port)
			details = "udp probe"

		case "http":
			code, err := checkHTTP(*target, false)
			up = err == nil && code >= 200 && code < 400
			details = fmt.Sprintf("HTTP %d err=%v", code, err)

		case "https":
			code, err := checkHTTP(*target, true)
			up = err == nil && code >= 200 && code < 400
			details = fmt.Sprintf("HTTPS %d err=%v", code, err)

		case "dns":
			up = checkDNS(*target)
			details = "DNS lookup"

		default:
			fmt.Println("Unknown protocol")
			return
		}

		if *jsonOut {
			printJSON(map[string]interface{}{
				"target":  *target,
				"port":    *port,
				"proto":   *proto,
				"up":      up,
				"details": details,
				"time":    time.Now().Format(time.RFC3339),
			})
		} else {
			if up {
				output.Up(fmt.Sprintf("%s %s:%d %s", *proto, *target, *port, details))
			} else {
				output.Down(fmt.Sprintf("%s %s:%d %s", *proto, *target, *port, details))
			}
		}

		if !*monitor {
			break
		}
		time.Sleep(time.Duration(*interval) * time.Second)
	}
}

/*
====================
SUPPORT FUNCTIONS
====================
*/

// --- GEOIP ---
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
	if ips, err := net.LookupIP(target); err == nil && len(ips) > 0 {
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

// --- PHONE ---
func lookupPhone(number string) {
	num, err := phonenumbers.Parse(number, "")
	if err != nil {
		fmt.Println("Error parsing number:", err)
		return
	}
	info := map[string]interface{}{
		"Raw":    number,
		"Valid":  phonenumbers.IsValidNumber(num),
		"Region": phonenumbers.GetRegionCodeForNumber(num),
		"Type":   phonenumbers.GetNumberType(num),
	}
	printJSON(info)
}

// --- PORTSCAN ---
func scanPorts(host string, start, end int) []int {
	open := []int{}
	timeout := 300 * time.Millisecond
	for port := start; port <= end; port++ {
		addr := fmt.Sprintf("%s:%d", host, port)
		conn, err := net.DialTimeout("tcp", addr, timeout)
		if err == nil {
			open = append(open, port)
			conn.Close()
		}
	}
	return open
}

// --- HELPERS ---
func checkUDP(target string, port int) bool {
	conn, err := net.DialTimeout("udp", fmt.Sprintf("%s:%d", target, port), time.Second)
	if err != nil {
		return false
	}
	defer conn.Close()
	_, err = conn.Write([]byte("ping"))
	return err == nil
}

func checkDNS(target string) bool {
	_, err := net.LookupHost(target)
	return err == nil
}

func checkHTTP(target string, https bool) (int, error) {
	scheme := "http"
	if https {
		scheme = "https"
	}
	resp, err := http.Get(scheme + "://" + target)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()
	return resp.StatusCode, nil
}

func printJSON(data interface{}) {
	js, _ := json.MarshalIndent(data, "", "  ")
	fmt.Println(string(js))
}
