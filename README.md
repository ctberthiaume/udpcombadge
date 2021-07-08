# udpcombadge

Simple Golang program to send and listen for UDP messages

This project is essentially a copy and paste of Vaughn Iverson's [https://github.com/vsivsi/udpsend](https://github.com/vsivsi/udpsend) and [https://github.com/vsivsi/udplisten](https://github.com/vsivsi/udplisten) into a single new CLI app, with modifications to use modern Go modules and to replace `github.com/jessevdk/go-flags` with `github.com/spf13/cobra`.

## Build

`go build`

## Run

### udpcombadge listen

```
Listen for UDP messages

Usage:
  udpcombadge listen [flags]

Flags:
  -b, --buffer uint   Max receive buffer size (default 1500)
  -f, --file string   Append received data to
  -h, --help          help for listen
  -H, --host string   Interface IP to bind to (default "0.0.0.0")
  -n, --newline       Add newline to end of each message
  -p, --port uint16   UDP port to bind to (default 1234)
  -q, --quiet         Suppress informational status on stderr
```

### udpcombadge send

```
Send UDP messages

Usage:
  udpcombadge send [flags]

Flags:
  -h, --help          help for send
  -H, --host string   IP destination address (default "255.255.255.255")
  -m, --msg string    Data to send
  -p, --port uint16   UDP destination port (default 1234)
```
