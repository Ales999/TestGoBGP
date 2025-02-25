# TestGoBGP

Testing access API GoBGP

Простой клиент позволяющий получить AS-pATH по префиксу или по IP от локального или удаленного GoBGP v3.x по его API

Example Output:

```text
 .\TestGoBGP.exe
  attrs: [{Origin: i} {AsPath: 64503 64502 64500 64504} {Nexthop: 172.24.1.1}]
  64503 64502 64500 64504
```

Есть поддержка ключей командной строки:

```text
$ ./TestGoBGP  -h
flag needs an argument: -h
Usage of ./TestGoBGP:
  -h string
        GoBgp host:port, examole: 192.168.1.11:50051 (default "127.0.0.1:50051")
  -n string
        Neigbror, example: 172.24.1.1 (default "172.24.1.1")
  -t string
        find target, example: 104.0.0.1 (default "104.0.0.1")
```
