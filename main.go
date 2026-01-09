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

// ------- GeoIP advanced lookup -------
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

// ------- Phone parsing lookup -------
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

func main() {
	target := flag.String("target", "localhost", "Target host or IP")
	port := flag.String("port", "80", "Port number")
	proto := flag.String("proto", "tcp", "Protocol: tcp | udp | http | https")
	phoneNum := flag.String("phone", "", "Phone number metadata lookup")
	jsonOut := flag.Bool("json", false, "JSON output")
	monitor := flag.Bool("monitor", false, "Continuously monitor target")
	interval := flag.Int("interval", 5, "Check interval in seconds")
	showVer := flag.Bool("version", false, "Show version info")
	geoip := flag.Bool("geoip", false, "GeoIP and ASN info for the target host/IP")
	flag.Parse()

	if *showVer {
		fmt.Printf("%s v%s by %s\n", version.Name, version.Version, version.Author)
		return
	}
	output.Banner()

	// GeoIP lookup if requested
	if *geoip {
		info, err := AdvancedGeoIPLookup(*target)
		if err != nil {
			fmt.Println("GeoIP lookup failed:", err)
			return
		}
		fmt.Printf("\nGeoIP/ASN info for %s\n", *target)
		fmt.Printf("  City:      %s\n", info.City)
		fmt.Printf("  Region:    %s\n", info.Region)
		fmt.Printf("  Country:   %s\n", info.Country)
		fmt.Printf("  Latitude:  %f\n", info.Lat)
		fmt.Printf("  Longitude: %f\n", info.Lon)
		fmt.Printf("  ASN:       %s\n", info.ASN)
		fmt.Printf("  Org:       %s\n", info.Org)
		fmt.Printf("  Hostname:  %s\n", info.Hostname)
		fmt.Println("")
		return
	}

	if *phoneNum != "" {
		LocalPhoneLookup(*phoneNum)
		return
	}

	var wasUp *bool
	var downSince time.Time

	for {
		up := false
		result := map[string]any{
			"target": *target,
			"proto":  *proto,
			"time":   time.Now().Format(time.RFC3339),
		}

		switch *proto {
		case "tcp":
			ok, latency := netcheck.TCP(*target, *port)
			up = ok
			result["up"] = ok
			result["latency_ms"] = latency

		case "udp":
			ok := netcheck.UDP(*target, *port)
			up = ok
			result["up"] = ok

		case "http", "https":
			code, latency := netcheck.HTTP(*proto+"://"+*target+":"+*port, true)
			up = code > 0
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
