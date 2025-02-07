# WarpRelay
A simple reverse proxy that utilises YAML files intended for use with HTTP servers
## Why?
- Despite there being many reverse proxies like [Caddy]("https://caddyserver.com") which is based on go and offers much more functionality, i decided to make this
- Why? I am yet to know why 😅 
- Also, I got bored by CaddyFiles
## How to use
- Get the binary from the releases page
- You can then Execute the binary from terminal by running `./warp -s` such that the program can set up the required Config Directory
- From there create a file called **config.YAML** inside the **Config** directory based on the sample [here](#yaml-sample)
## YAML Sample
```yaml
proxyPort: "7080"
endpoints:
  - backendUrl : http://localhost:9080
    proxyEndpoint : /sample
  - backendUrl: https://localhost:9081
    proxyEndpoint: /another
```