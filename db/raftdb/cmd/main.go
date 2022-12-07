package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/ensn1to/experiment/tree/master/db/raftdb/pkg/store"
)

// command line parameters
var (
	httpAddr string
	raftAddr string
	joinAddr string
	nodeID   string
)

func init2() {
	flag.StringVar(&httpAddr, "http", "localhost:8091", "Set the http bind addr")
	flag.StringVar(&raftAddr, "http", "localhost:8089", "Set the  raft addr")
	flag.StringVar(&joinAddr, "join", "", "set join node addr")
	flag.StringVar(&nodeID, "id", "", "set node id")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "usage: %s [options]<raft-data-path>\n", os.Args[0])
		flag.PrintDefaults()
	}
}

func main() {
	flag.Parse()

	if flag.NArg() == 0 {
		fmt.Fprintf(os.Stderr, "no raft storage dir specifiled")
		os.Exit(1)
	}

	raftDir := flag.Arg(0)
	if raftDir == "" {
		fmt.Fprintf(os.Stderr, "no raft storage dir specifiled")
		os.Exit(1)
	}
	if err := os.MkdirAll(raftDir, 0o700); err != nil {
		fmt.Printf("failed to create raft storage dir: %s", err.Error())
	}

	s := store.New()
	s.RaftDir = raftDir
	s.RaftBind = raftAddr
	if err := s.Open(true, nodeID); err != nil {
		fmt.Printf("failed to open store: %s", err.Error())
	}
}
