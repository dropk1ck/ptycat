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
    var listenMode bool

    flag.BoolVar(&listenMode, "l", false, "listen locally for connection")
    flag.IntVar(&listenPort, "p", 0, "port for local listener")
    flag.BoolVar(&verboseMode, "v", false, "turn up verbosity")
    flag.Parse()

    // default to info-level debugging, 'verbose' turns on debug output
    level := "info"
    if verboseMode {
        level = "debug"
    }
    logLevel, _ := logrus.ParseLevel(level)
    logrus.SetLevel(logLevel)

    args := flag.Args()

    // argument sanity checks
    if (listenMode && listenPort == 0) || (!listenMode && listenPort != 0) {
        flag.Usage()
        return
    }

    if listenMode && (len(args) == 2) {
        // choose a mode, not both
        fmt.Println("Must choose either listen mode or connect mode, not both")
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

// simple TCP connect to specified addr/port combo
func doConnect(args []string) net.Conn {
    conn, err := net.Dial("tcp", args[0]+":"+args[1])
    if err != nil {
        logrus.Fatal(err)
    }
    return conn
}

// listen on all interfaces on specified port
func doListen(listenPort int) net.Conn {
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
    return conn
}
