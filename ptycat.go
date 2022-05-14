package main

import (
    "flag"
    "io"
    "net"
    "os"
    "strconv"
    "github.com/dropk1ck/pwnlog"
    "golang.org/x/crypto/ssh/terminal"
)

var logger *pwnlog.Logger


func main() {
    var listenPort int
    var verboseMode bool
    var listenMode bool

    flag.BoolVar(&listenMode, "l", false, "listen locally for connection")
    flag.IntVar(&listenPort, "p", 0, "port for local listener")
    flag.BoolVar(&verboseMode, "v", false, "turn up verbosity")
    flag.Parse()

    // default to info-level debugging, 'verbose' turns on debug output
    logLevel := pwnlog.InfoLevel 
    if verboseMode {
        logLevel = pwnlog.DebugLevel 
    }
    logger = pwnlog.New(logLevel)   
    args := flag.Args()

    // argument sanity checks
    if (listenMode && listenPort == 0) || (!listenMode && listenPort != 0) {
        flag.Usage()
        return
    }

    if listenMode && (len(args) == 2) {
        // choose a mode, not both
        logger.Error("Must choose either listen mode or connect mode, not both")
        flag.Usage()
        return
    }

    // either go in to listen mode or connect-in
    var conn net.Conn
    if listenMode {
        conn = doListen(listenPort)
    } else {
        conn = doConnect(args) 
    }

    // either way we're now connected, setup terminal
    logger.Debug("Creating raw terminal, happy hacking")
    oldState, err := terminal.MakeRaw(int(os.Stdin.Fd()))
    if err != nil {
        logger.Fatal(err.Error())
    }
   
    // automatically restore terminal settings on exit
    defer func() { _ = terminal.Restore(int(os.Stdin.Fd()), oldState) }()

    // we're off to the races
    go func() { _, _ = io.Copy(os.Stdout, conn) }()
    _, err = io.Copy(conn, os.Stdin)

    logger.Debug("Lost remote terminal")
}

// simple TCP connect to specified addr/port combo
func doConnect(args []string) net.Conn {
    conn, err := net.Dial("tcp", args[0]+":"+args[1])
    if err != nil {
        logger.Fatal(err.Error())
    }
    return conn
}

// listen on all interfaces on specified port
func doListen(listenPort int) net.Conn {
    listenPortStr := strconv.Itoa(listenPort)
    logger.Debug("Listening on port " + listenPortStr)
    ln, err := net.Listen("tcp", ":"+listenPortStr)
    if err != nil {
        logger.Fatal(err.Error())        
    }

    conn, err := ln.Accept()
    if err != nil {
        logger.Fatal(err.Error())
    }
    logger.Debug("Accepted connection")
    return conn
}
