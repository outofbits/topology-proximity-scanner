<p align="center">
  <h1 align="center">
    Topology Proximity Scanner
    <br/>
    <a href="https://github.com/godano/cardano-lib/blob/master/LICENSE" ><img alt="license" src="https://img.shields.io/badge/license-MIT%20License%202.0-E91E63.svg?style=flat-square" /></a>
  </h1>
</p>

This application takes the topology file for a stake pool and checks the proximity (i.e. round-trip-time) for all the
specified nodes. If a DNS name is specified for a node, this application resolves it to the registered A/AAA records and
checks the proximity for all entries.  

## Build

```
git clone https://github.com/outofbits/topology-proximity-scanner.git
cd topology-proximity-scanner
go mod vendor
go build
```

## Usage

```
topology-proximity-scanner -i <topology-file-path> [-n <number>][-o <csv-output-file>]
  -i string
        path to the topology file
  -n int
        number of data points to collect for each node (default 3)
  -o string
        optional path to a csv file to which results shall be written
```


Example:
```
topology-proximity-scanner -i topology.json -n 10 -o results.csv
```