package main


import (
"flag"
"fmt"


"cyph3r/internal/geo"
"cyph3r/internal/netcheck"
"cyph3r/internal/output"
"cyph3r/internal/phone"
"cyph3r/internal/version"
)


func main() {
target := flag.String("target", "localhost", "Target host")
port := flag.String("port", "80", "Port")
proto := flag.String("proto", "tcp", "tcp|udp|http|https")
jsonOut := flag.Bool("json", false, "JSON output")
phoneNum := flag.String("phone", "", "Phone lookup")
showVer := flag.Bool("version", false, "Show version")
flag.Parse()


if *showVer {
fmt.Println(version.Name, version.Version, "by", version.Author)
return
}


geo.Lookup(*target)


result := map[string]any{"target": *target, "proto": *proto}


switch *proto {
case "tcp":
ok, lat := netcheck.TCP(*target, *port)
result["up"] = ok
result["latency_ms"] = lat
case "udp":
result["up"] = netcheck.UDP(*target, *port)
case "http", "https":
code, lat := netcheck.HTTP(*proto+"://"+*target+":"+*port, true)
result["status"] = code
result["latency_ms"] = lat
}


if *phoneNum != "" {
phone.Lookup(*phoneNum)
}


if *jsonOut {
output.PrintJSON(result)
} else {
fmt.Println(result)
}
}
