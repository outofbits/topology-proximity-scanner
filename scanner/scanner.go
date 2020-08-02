package scanner

import (
    "fmt"
    "github.com/cloverstd/tcping/ping"
    "github.com/outofbits/stakepool-ops-lib/topology"
    "net"
    "time"
)

type ScanResult struct {
    Node    topology.NodeConfig
    Address string
    Result  *ping.Result
}

func min(a, b int) int {
    if a < b {
        return a
    }
    return b
}

// scans the proximity to all given nodes using "n" data points. if
// the node specifies a DNS name, it resolves the DNS names to all
// address entries and checks the proximity to all entries.
func Scan(nodes []topology.NodeConfig, nDataPoints int, parallelPings int) []ScanResult {
    scanResults := make([]ScanResult, 0)
    resultChan := make(chan []ScanResult)
    chunks := len(nodes) / parallelPings
    for c := 0; c <= chunks; c++ {
        currentNodes := nodes[(c * parallelPings):(min((c+1)*parallelPings, len(nodes)))]
        for _, node := range currentNodes {
            go scanHost(node, nDataPoints, resultChan)
        }
        for i := 0; i < len(currentNodes); i++ {
            results := <-resultChan
            for _, result := range results {
                scanResults = append(scanResults, result)
            }
        }
    }
    return scanResults
}

func scanHost(node topology.NodeConfig, nDataPoints int, resultChan chan []ScanResult) {
    addresses, err := net.LookupHost(node.HostAddress)
    if err == nil {
        addressResultChan := make(chan ScanResult)
        for i, address := range addresses {
            fmt.Printf("[%s:%v#%d(%s)] Collecting %d data points.\n", node.HostAddress, node.Port, i+1,
                address, nDataPoints)
            go scanAddress(address, node, nDataPoints, addressResultChan)
        }
        results := make([]ScanResult, 0)
        for i := 0; i < len(addresses); i++ {
            result := <-addressResultChan
            results = append(results, result)
        }
        resultChan <- results
    }
}

func scanAddress(address string, node topology.NodeConfig, nDataPoints int, addressResultChan chan ScanResult) {
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
    addressResultChan <- ScanResult{
        Node:    node,
        Address: address,
        Result:  result,
    }
    fmt.Printf("[%s:%v#%s] Average RTT=%v\n", node.HostAddress, node.Port, address, result.Avg())
}
