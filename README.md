# shiro_waf

A white list based firewall. Why? Because I have seen how much effort and time is gone into setup black list rules, improving them, taking care of bypasses. Oh Im tired of the bypasses.

Also I wanted to learn GO.

## Todo

- [x] test web server 
- [x] proxy server 
- [x] connect both servers 
- [x] create custom golang web server 
- [x] rules regex file format 
- [x] filter and block traffic based on rules
- [] Convert current List of rules to Dictionary for faster lookup
- [] Improve Blocking based on:
    - URI
    - Body
    - parameters
    - headers
- [] auto generate rules file
- [] expand from http protocol to multiple protocols
- [] make the tool CI/CD friendly
- [] dockerize
- [] host the the firewall (ngrok / AWS)
- [] GUI using maybe dart/kotlin/javascript

Absolute **CHAD** of an article: https://www.codedodle.com/go-reverse-proxy-example.html