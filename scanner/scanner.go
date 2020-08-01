package scanner

import (
    "fmt"
    "github.com/cloverstd/tcping/ping"
    "net"
    "time"
)

type ScanResult struct {
    Node    NodeConfig
    Address string
    Result  *ping.Result
}

// scans the proximity to all given nodes using "n" data points. if
// the node specifies a DNS name, it resolves the DNS names to all
// address entries and checks the proximity to all entries.
func Scan(nodes []NodeConfig, nDataPoints int) []ScanResult {
    scanResults := make([]ScanResult, 0)
    for _, node := range nodes {
        addresses, err := net.LookupHost(node.HostAddress)
        for i, address := range addresses {
            fmt.Printf("[%s:%v#%d(%s)] Collecting %d data points.\n", node.HostAddress, node.Port, i+1,
                address, nDataPoints)
            if err == nil {
                pinger := ping.NewTCPing()
                pinger.SetTarget(&ping.Target{
                    Protocol: ping.TCP,
                    Host:     address,
                    Port:     node.Port,
                    Counter:  nDataPoints,
                    Interval: 1 * time.Second,
                    Timeout:  10 * time.Second,
                })
                done := pinger.Start()
                _ = <-done
                result := pinger.Result()
                scanResults = append(scanResults, ScanResult{
                    Node:   node,
                    Address: address,
                    Result: result,
                })
                fmt.Printf("[%s:%v] Average RTT=%v\n", node.HostAddress, node.Port, result.Avg())
            }
        }
    }
    return scanResults
}
