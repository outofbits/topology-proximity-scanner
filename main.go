package main

import (
    "flag"
    "fmt"
    "github.com/outofbits/stakepool-ops-lib/topology"
    "github.com/outofbits/topology-proximity-scanner/csv"
    "github.com/outofbits/topology-proximity-scanner/scanner"
    "os"
)

const defaultAppName = "topology-proximity-scanner"
const appVersion = "1.1.1"

// checks whether a file exists at the given path.
func fExists(name string) bool {
    stat, err := os.Stat(name)
    if err != nil {
        return !os.IsNotExist(err)
    } else {
        return !stat.IsDir()
    }
}

// gets the name of this application as it should be
// displayed in help messages.
func appName() string {
    name := defaultAppName
    if os.Args != nil && len(os.Args) > 0 {
        name = os.Args[0]
    }
    return name
}

func main() {
    input := flag.String("i", "", "path to the topology file")
    output := flag.String("o", "", "optional path to a csv file to which results shall be written")
    dataPoints := flag.Int("n", 3, "number of data points to collect for each node")
    parallelPings := flag.Int("p", 10, "number of parallel pings")
    version := flag.Bool("v", false, "print the version of this application")
    flag.Parse()

    if *version {
        fmt.Printf("%s v%s\n", defaultAppName, appVersion)
        os.Exit(0)
    }

    if *input != "" {
        if fExists(*input) {
            if *dataPoints <= 0 {
                _, _ = fmt.Fprintf(os.Stderr, "The given number of data points must be greater than zero.\n")
                os.Exit(1)
            }
            if *parallelPings <= 0 {
                _, _ = fmt.Fprintf(os.Stderr, "The given number of parallel pings must be greater than zero.\n")
                os.Exit(1)
            }
            var topologyConfig, err = topology.ReadTopologyFile(*input)
            if err == nil {
                results := scanner.Scan(topologyConfig.Producers, *dataPoints, *parallelPings)
                if *output != "" {
                    err = csv.WriteScanResult(results, *output)
                    if err != nil {
                        _, _ = fmt.Fprintf(os.Stderr, "An error occured while writing scan results to \"%s\". %s.\n",
                            *output, err.Error())
                        os.Exit(1)
                    }
                }
            } else {
                _, _ = fmt.Fprintf(os.Stderr, "Could not parse the specified topology file \"%s\". %s.\n", *input,
                    err.Error())
                os.Exit(1)
            }
        } else {
            _, _ = fmt.Fprintf(os.Stderr, "No topology file can be found at \"%s\".\n", *input)
            os.Exit(1)
        }
    } else {
        _, _ = fmt.Fprint(os.Stderr, "You need to specify the path to a topology file.\n")
        fmt.Printf("%s -i <topology-file-path> [-n <number>][-o <csv-output-file>]\n", appName())
        flag.PrintDefaults()
        os.Exit(1)
    }
}
