package main

import (
    "flag"
    "fmt"
    "io"
    "net"
    "os"
    "strconv"
    "github.com/sirupsen/logrus"
    "golang.org/x/crypto/ssh/terminal"
)


func main() {
    var listenPort int
    var verboseMode bool

    flag.IntVar(&listenPort, "p", 0, "port for local listener")
    flag.BoolVar(&verboseMode, "v", false, "turn up verbosity")
    flag.Parse()

    // default to info-level debugging, 'verbose' turns on debug output
    level := "info"
    if verboseMode {
        level = "debug"
    }
    logLevel, err := logrus.ParseLevel(level)
    if err != nil {
        logLevel = logrus.InfoLevel
    }
    logrus.SetLevel(logLevel)

    // ensure we got a port if we're listening
    if listenPort == 0 {
        fmt.Println("Port must be specified with -p")
        return
    }

    listenPortStr := strconv.Itoa(listenPort)
    logrus.Debug("Listening on port " + listenPortStr)
    ln, e := net.Listen("tcp", ":"+listenPortStr)
    if e != nil {
        logrus.Fatal(e)        
    }

    conn, e := ln.Accept()
    if e != nil {
        logrus.Fatal(e)
    }
    logrus.Debug("Accepted connection")
    logrus.Debug("Creating raw terminal, happy hacking")
    oldState, e := terminal.MakeRaw(int(os.Stdin.Fd()))
    if e != nil {
        logrus.Fatal(e)
    }
   
    // automatically restore terminal settings on exit
    defer func() { _ = terminal.Restore(int(os.Stdin.Fd()), oldState) }()

    // we're off to the races
    go func() { _, _ = io.Copy(os.Stdout, conn) }()
    _, e = io.Copy(conn, os.Stdin)

}
