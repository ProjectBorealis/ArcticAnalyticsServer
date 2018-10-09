# ArcticAnalytics-Server

Basic server for saving performance analytics to a file.

### Building

##### With Go tools
`$ go build`

##### With Docker / Makefile
`$ make build`

### Running

##### With binary
```
# only use --behind-proxy if you trust X-Forwarded-For headers.
$ ./aa-server --behind-proxy --addr :80
```

##### With docker
```
$ make run
```