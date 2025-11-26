package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"os"

	"github.com/spf13/cobra"
)

var (
	Name    string = "udpcombadge"
	Version string = "v0.1.1"

	ListenHost    string
	ListenPort    uint16
	ListenFile    string
	ListenBuffer  uint
	ListenNewline bool
	ListenQuiet   bool

	SendHost string
	SendPort uint16
	SendMsg  string
)

func main() {

	var cmdVersion = &cobra.Command{
		Use:   "version",
		Short: "Print version to stdout and exit",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(Version)
		},
	}

	var cmdListen = &cobra.Command{
		Use:   "listen",
		Short: "Listen for UDP messages",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("listen")
			addr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("%v:%d", ListenHost, ListenPort))
			if err != nil {
				log.Panic(err)
			}

			l, err := net.ListenUDP("udp", addr)
			if err != nil {
				log.Panic(err)
			}
			if !ListenQuiet {
				log.Printf("Starting udplistener at %v:%v", ListenHost, ListenPort)
			}
			for {
				handleClient(l, ListenFile)
			}
		},
	}
	cmdListen.Flags().StringVarP(&ListenHost, "host", "H", "0.0.0.0", "Interface IP to bind to")
	cmdListen.Flags().Uint16VarP(&ListenPort, "port", "p", 1234, "UDP port to bind to")
	cmdListen.Flags().StringVarP(&ListenFile, "file", "f", "", "Append received data to")
	cmdListen.Flags().UintVarP(&ListenBuffer, "buffer", "b", 1500, "Max receive buffer size")
	cmdListen.Flags().BoolVarP(&ListenNewline, "newline", "n", false, "Add newline to end of each message")
	cmdListen.Flags().BoolVarP(&ListenQuiet, "quiet", "q", false, "Suppress informational status on stderr")

	var cmdSend = &cobra.Command{
		Use:   "send",
		Short: "Send UDP messages",
		Run: func(cmd *cobra.Command, args []string) {
			c, err := net.Dial("udp", fmt.Sprintf("%v:%d", SendHost, SendPort))
			if err != nil {
				log.Panic(err)
			}

			defer c.Close()

			var output []byte

			if len(SendMsg) == 0 {
				output, err = ioutil.ReadAll(os.Stdin)
				if err != nil {
					log.Panic(err)
				}
			} else {
				output = []byte(SendMsg)
			}

			n, err := c.Write(output)
			if err != nil {
				log.Panic(err)
			} else {
				log.Printf("Wrote %d bytes\n", n)
			}
		},
	}
	cmdSend.Flags().StringVarP(&SendHost, "host", "H", "255.255.255.255", "IP destination address")
	cmdSend.Flags().Uint16VarP(&SendPort, "port", "p", 1234, "UDP destination port")
	cmdSend.Flags().StringVarP(&SendMsg, "msg", "m", "", "Data to send")

	var rootCmd = &cobra.Command{Use: Name}
	rootCmd.AddCommand(cmdVersion, cmdListen, cmdSend)
	rootCmd.CompletionOptions.DisableDefaultCmd = true
	rootCmd.Execute()
}

func handleClient(conn *net.UDPConn, fn string) {
	b := make([]byte, ListenBuffer)
	n, addr, e := conn.ReadFromUDP(b)
	if e != nil {
		log.Printf("Read from UDP failed, err: %v", e)
		return
	}
	if !ListenQuiet {
		log.Printf("Read from client(%v:%v), len: %v\n", addr.IP, addr.Port, n)
	}
	if ListenNewline && n < len(b) && b[n-1] != '\n' {
		b[n] = '\n'
		n += 1
	}
	if len(ListenFile) != 0 {
		f, err := os.OpenFile(ListenFile, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0600)
		if err != nil {
			log.Printf("Open file failed, err: %v", err)
			return
		}
		defer f.Close()
		if _, err = f.Write(b[:n]); err != nil {
			log.Printf("Write file failed, err: %v", err)
			return
		}
	} else { // Write to stdout
		if _, err := os.Stdout.Write(b[:n]); err != nil {
			log.Printf("Write stdout failed, err: %v", err)
			return
		}
	}
}
