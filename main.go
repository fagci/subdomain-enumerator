package main

import (
	"bufio"
	"flag"
	"fmt"
	"net"
	"os"
	"reflect"
	"runtime"
)

func check(hostname string, fake_ips []net.IP, ch chan string) {
	if ips, err := net.LookupIP(hostname); err == nil && !reflect.DeepEqual(fake_ips, ips) {
		ch <- hostname
	} else {
		ch <- ""
	}
}

func scan(dict *os.File, target string) {
	dict_scanner := bufio.NewScanner(dict)

	fake_sd := "f4k3sd"
	fake_ips := []net.IP{}

	if ips, err := net.LookupIP(fake_sd + "." + target); err == nil {
		fake_ips = ips
	}

	ch := make(chan string, 128)

	for dict_scanner.Scan() {
		hostname := dict_scanner.Text() + "." + target
		fmt.Printf("\r[*] %s\u001b[0J", hostname)
		go check(hostname, fake_ips, ch)
		if sd := <-ch; sd != "" {
			fmt.Printf("\r[+] %s\u001b[0J\n", sd)
		}
	}

    runtime.Goexit()
	fmt.Println("\rDone\u001b[0J")
}

func main() {
	dict_path := flag.String("d", "", "subdomains dictionary file")
	target := flag.String("t", "", "target domain")
	flag.Parse()

	if flag.NFlag() < 2 {
		flag.Usage()
		return
	}

	dict, err := os.Open(*dict_path)
	if err != nil {
		fmt.Printf("Cannot open dict. (%s)\n", err)
		return
	}
	defer dict.Close()

	scan(dict, *target)

}
