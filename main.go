package main

import (
	"bufio"
	"flag"
	"fmt"
	"net"
	"os"
	"reflect"
	"sync"
)

func get_wildcard_ips(target string) []net.IP {
	fake_sd := "f4k3sd"
	fake_ips := []net.IP{}

	if ips, err := net.LookupIP(fake_sd + "." + target); err == nil {
		fake_ips = ips
	}
    return fake_ips
}

func check(hostname string, fake_ips []net.IP, wg *sync.WaitGroup) {
	defer wg.Done()
	fmt.Printf("\r[*] %s\u001b[0J", hostname)
	if ips, err := net.LookupIP(hostname); err == nil && !reflect.DeepEqual(fake_ips, ips) {
		fmt.Printf("\r[+] %s\u001b[0J\n", hostname)
	}
}

func scan(dict *os.File, target string) {
	dict_scanner := bufio.NewScanner(dict)

    fake_ips := get_wildcard_ips(target)

	var wg sync.WaitGroup

	for dict_scanner.Scan() {
		hostname := dict_scanner.Text() + "." + target
		wg.Add(1)
		go check(hostname, fake_ips, &wg)
	}

	wg.Wait()
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
