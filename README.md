# WarpRelay
A simple reverse proxy that utilises YAML files intended for use with HTTP servers
## Why?
- Despite there being many reverse proxies (like [Caddy]("https://caddyserver.com") which is based on go and offers much more functionality) I decided to make this
- Why? I wanted to know how
- Also, I got bored by CaddyFiles and I find YAML to be relaxing for my smooth brain 
## How to use
- Get the binary from the releases page
- You can then Execute the binary from terminal by running `./warp` to use Default options or check [options](#options) for other methods to initialize the server
- From there create a file called **config.YAML** inside the **Config** directory based on the sample [here](#yaml-sample)
## Options
The syntax for running warp from terminal is:
```
warp [flags] [...]
```
The available flags are as follows:
```
-v,                             Enable verbose output on default settings
-cv </path/to/config/file>,     Uses path as config file and enables verbose output
-h,                             Shows Help Text
-c </path/to/config/file>,      Uses path as config file
```
## YAML Sample
```yaml
proxyPort: "7080"   # Port at which proxy server will run
certPath : /path/to/server.crt    # Location of SSL cert file  
keyPath: /path/to/server.key    # Location of SSL key file
endpoints:
  - backendUrl : http://localhost:9080    # IP:Port of your webserver
    proxyEndpoint : /one    # Endpoint where above webserver is to be accessible from
  - backendUrl: https://localhost:9081
    proxyEndpoint: /two
```