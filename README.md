# Puppet monitor
Puppet report analyzer and node activity monitor

## Getting Started

This tool has two major roles:
* run as master process (no command line args)
* run as child process to get info or send a commands to master process (have some command line args)

### Workflow

Master process runs on any server wherever you want to, opens TCP ports on desired network interface and starts listening.
Puppet server should be configured to send node reports to this IP:PORT. Every report that comes to puppet monitoring
tool is analized and its results stored in memory. Any child process can establish connection to master process via RPC port
and get some piece of information, like print table with nodes and their status, print overall status (with or without errors) and so on...

Child process can be easily used to setup a rule in almost any monitoring system, which working with stdout.

So. Steps are

* Build
* Setup config and run
* Setup puppet server to sends reports into this tool
* Add to monitoring system (like nagios)
* ...
* Profit!

### monitor.conf:

File format: **json**

Default location: **/etc/puppetlabs/puppet/monitor.conf**

defaults is:
```javascript
{
  "port": 3840, // tcp port to recieve data on
  "ip": "127.0.0.1", // one of the local ip to recieve data on
  "pid": "/var/run/puppet_monitoring.pid", // pid file - just for information
  "rpc": 3841, // tcp port to open for master process and to communicate from child process
  "ctime": 35 // node report max delay, in minutes
}
```
you can omit any of this options since it has default values

### Command line options

```tex
Usage:
  puppet_monitoring [OPTIONS]

Application Options:
  -p, --print    print current environment collections
  -s, --status   print status
  -e, --error    print status with errors
  -v, --version  print version
      --stop     send selfkill signal to master process
  -r, --remove=  remove all data about specified host
      --rpc=     set rpc params for master process communication

Help Options:
  -h, --help     Show this help message
```

#### Examples:

##### Print actual info about all nodes
```
./puppet_monitoring --print --rpc=puppet-server.yourdomain.com:3841
or
./puppet_monitoring -p --rpc=puppet-server.yourdomain.com:3841
```

##### Print status
###### Without errors
```
./puppet_monitoring --status --rpc=puppet-server.yourdomain.com:3841
or
./puppet_monitoring -s --rpc=puppet-server.yourdomain.com:3841
```
###### With errors
```
./puppet_monitoring --status --error --rpc=puppet-server.yourdomain.com:3841
or
./puppet_monitoring -s -e --rpc=puppet-server.yourdomain.com:3841
```

To build and setup, follow the next instrustions

### Prerequisites
* [Go 1.7](https://golang.org/) - Build tool

### Building

* Clone this repository to your GO workspace
* Clone GO package [go-flags](https://godoc.org/github.com/jessevdk/go-flags)
```
go get github.com/jessevdk/go-flags
```
* CWD to cloned copy of this repository
* Build project with
```
go build
```
here you should get the binary file

### Getting binary to work

* Copy binary file to somewhere
* Create file /etc/puppetlabs/puppet/monitor.conf, if you don't like defaults, and setup'em
* Run binary file without args (or as system service)

Systemd service sample: *puppet_monitoring.service*
```
[Unit]
Description=Puppet node monitoring service

[Service]
WorkingDirectory=/usr/sbin/
ExecStart=/bin/sh -c "/usr/sbin/puppet_monitoring > /var/log/puppetlabs/puppet_monitor.log"
Restart=always

[Install]
WantedBy=multi-user.target
```

## Setup puppet server

for this, you'll need
[Logstash reporter](https://github.com/elastic/puppet-logstash-reporter) - just follow the instructions, but don't forget to point IP:PORT to this tool

