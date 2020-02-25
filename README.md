# request-forwarder

## Description

This application is used to intercept HTTP requests to an application running locally which only listens on the local address, typically to limit the ability external users from directly interacting with the service.  Although other methods can be used to accomplish this, this service is to be used in situations such as:
    - When the target HTTP endpoint doesn't support various authentication method 
    - You don't want to configure local firewall/iptables rules

## Usage

Command:
```
./request-forwarder -a <addr> -p <port> -m <metrics_port> -ra <remote_addr> -rp <remote_port> -w <allowed_methods>
```

Flags:
```
  -a string
    	The address to bind to. (default "127.0.0.1")
  -m int
    	The port on which to expose Prometheus metrics. (default 9555)
  -p int
    	The port to bind to. (default 8080)
  -ra string
    	The remoate address to bind to. (default "127.0.0.1")
  -rp int
    	The remote port to bind to. (default 8500)
  -w string
    	Comma separated list of allowed methods. Empty means all.
```

Example:
----
This will start the process to listen on *:8080 and rewrite requests to 127.0.0.1:8500
```
./request-forwarder -a 0.0.0.0 -ra localdev -rp 8500 -w get,head 
```