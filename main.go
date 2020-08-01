package main

import (
    "flag"
    "fmt"
    "github.com/outofbits/topology-proximity-scanner/csv"
    "github.com/outofbits/topology-proximity-scanner/scanner"
    "os"
)

const defaultAppName = "topology-proximity-scanner"
const appVersion = "1.0.0"

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
    version := flag.Bool("v", false, "print the version of this application")
    flag.Parse()

    if *version {
        fmt.Printf("%s v%s\n", defaultAppName, appVersion)
        os.Exit(0)
    }

    if *input != "" {
        if fExists(*input) {
            var topology, err = scanner.ReadTopologyFile(*input)
            if err == nil {
                results := scanner.Scan(topology.Producers, *dataPoints)
                if *output != "" {
                    err = csv.WriteScanResult(results, *output)
                    if err != nil {
                        _, _ = fmt.Fprintf(os.Stderr, "An error occured while writing scan results to \"%s\". %s.\n",
                            *output, err.Error())
                    }
                }
            } else {
                _, _ = fmt.Fprintf(os.Stderr, "Could not parse the specified topology file\"%s\". %s.\n", *input,
                    err.Error())
            }
        } else {
            _, _ = fmt.Fprintf(os.Stderr, "No topology file can be found at \"%s\".\n", *input)
        }
    } else {
        _, _ = fmt.Fprint(os.Stderr, "You need to specify the path to a topology file.\n")
        fmt.Printf("%s -i <topology-file-path> [-n <number>][-o <csv-output-file>]\n", appName())
        flag.PrintDefaults()
    }
}
