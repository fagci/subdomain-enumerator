# Subdomain enumerator

## Features

- deals with wildcard subdomains
- 

## Usage

```sh
go run ./main.go -t 'example.com' -d ./dict/sub-1000.txt
```

Outputs:
```log
[+] mail.example.com
[+] my.example.com
[+] www.example.com
[+] apps.example.com
[+] test1.example.com
[+] admin.example.com
[+] account.example.com
[+] fr.example.com
[+] site.example.com
[+] docs.example.com
[+] app.example.com
[+] ns2.example.com
[+] es.example.com
[+] ru.example.com
[+] api.example.com
[+] backend.example.com
Done
```
