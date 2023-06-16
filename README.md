# Shiro WAF シロワフ शिरो वाफ 

A white list based firewall. Why? Because I have seen how much effort and time is gone into setup black list rules, improving them, taking care of bypasses. Oh Im tired of the bypasses.

Also I wanted to learn GO.

## Todo

- [x] test web server 
- [x] proxy server 
- [x] connect both servers 
- [x] create custom golang web server 
- [x] rules regex file format 
- [x] filter and block traffic based on rules
- [] Convert current List of rules to Dictionary for faster lookup (not needed for now)
- [] Improve Blocking based on:
    - [x] URI
    - [x] Body
    - parameters
    - headers (in progress)
- [] Introduce monitoring mode (to inspect requests and create rule file) and blocking mode 
- [] auto generate rules file (Not sure how I'll do it)
- [] expand from http protocol to multiple protocols (maybe a whole new application is a better idea)
- [] make the tool CI/CD friendly
- [] dockerize
- [] host the the firewall (ngrok / AWS)
- [] GUI using maybe dart/kotlin/javascript

### Kudos

1. Absolute **CHAD** of an article: https://www.codedodle.com/go-reverse-proxy-example.html
2. Chat GPT (I'm not kidding)
3. [Charm](https://charm.sh/) (That website tho o_o)
4. [Kenneth](https://kennethreitz.org/) for HTTPBin (which is my primary target for proxy)