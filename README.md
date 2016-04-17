# lustre-graph-driver
This is a Docker graph driver plugin (experimental in 1.10) that implements a Docker storage driver using a Lustre filesystem mount. 


## How to build
Install dependencies with go get (TODO list dependencies here). You will know which dependencies are missing when there are errors in the command below

``` sh
go build -v
```

## How to run
NOTE THIS ISN'T WORKING CURRENTLY

``` sh
[vagrant@localhost lustre-graph-driver]$ sudo ./lustre-graph-driver -D -s lustre
INFO[0000] listening on /run/docker/plugins/lustre.sock
 
DEBU[0000] root group found. gid: 0
```