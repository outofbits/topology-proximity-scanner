package csv

import (
    "bufio"
    "encoding/csv"
    "fmt"
    "github.com/outofbits/topology-proximity-scanner/scanner"
    "os"
    "strconv"
)

// writes the given scan results to the given file
func WriteScanResult(results []scanner.ScanResult, outputFile string) error {
    f, err := os.Create(outputFile)
    if err != nil {
        return err
    }
    defer f.Close()
    csvWriter := csv.NewWriter(bufio.NewWriter(f))
    err = csvWriter.Write([]string{"host", "port", "address", "average_rtt_ms", "successful_connects", "attempts"})
    if err != nil {
        return err
    }
    for _, entry := range results {
        node := entry.Node
        result := entry.Result
        err = csvWriter.Write([]string{node.HostAddress, strconv.Itoa(node.Port), entry.Address,
            fmt.Sprintf("%f", float64(result.Avg())/1e6), strconv.Itoa(result.SuccessCounter),
            strconv.Itoa(result.Counter)})
        if err != nil {
            return err
        }
    }
    csvWriter.Flush()
    return nil
}
