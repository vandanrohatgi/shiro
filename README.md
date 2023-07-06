# shiro

A white list based web application firewall. Why? Because I have seen how much effort and time is gone into setup of black list rules, improving them, taking care of bypasses. Oh Im tired of the bypasses.

Also I wanted to learn GO.


## Usage 

1. Download a binary for your system.
2. Running the program right and browsing to [http://localhost:8080](http://localhost:8080) right away will block all requests because there are no rules right now. You can either copy over the [rules.yaml.sample](rules.yaml.sample) file and remove ".sample" from it's name or enter monitoring mode and the program will create some rules for you.
3. Monitoring mode will not block anything but will create patterns from you browsing and store them in rules.yaml.
4. After you have a rules.yaml file you can start browsing and observe blocking by browsing a url that either does not have a rule for OR making a request to a URL which you did not interact with while generating the rules.

## Help

```
./shiro --help

Usage of ./shiro:
  -monitor
        Monitor proxy traffic and generate rules automatically
  -path string
        path to the rules file (default "rules.yaml")
  -proxyPort string
        port to host the proxy (default "8080")
  -targetURL string
        URL to proxy (default "https://httpbin.org/")
  -timeout int
        Timeout for the proxy requests (default 10)
  -verbose
        Output all types of logs
```

### Kudos

1. Absolute **CHAD** of an article: https://www.codedodle.com/go-reverse-proxy-example.html
2. Chat GPT (I'm not kidding)
3. [Charm](https://charm.sh/) for their log library (That website tho o_o)
4. [Kenneth](https://kennethreitz.org/) for HTTPBin (which is my primary target for proxy testing)
5. [itchyny](https://github.com/itchyny/) for the rassemble-go library
