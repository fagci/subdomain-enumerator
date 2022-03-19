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

type Result struct {
	Hostname string
	Resolved bool
	IPs      []net.IP
}

type Enumerator struct {
	Target    string
	DictPath  string
	waitGroup sync.WaitGroup
	ch        chan Result
	fakeIPs   []net.IP
	dict      *os.File
}

func (e *Enumerator) GetWildcardIPs() []net.IP {
	fake_hostname := "f4k3sd." + e.Target

	if ips, err := net.LookupIP(fake_hostname); err == nil {
		e.fakeIPs = ips
	}

	return e.fakeIPs
}

func (e *Enumerator) Check(hostname string) {
	defer e.waitGroup.Done()
	if ips, err := net.LookupIP(hostname); err == nil && !reflect.DeepEqual(e.fakeIPs, ips) {
		e.ch <- Result{Hostname: hostname, Resolved: true, IPs: ips}
	} else {
		e.ch <- Result{Hostname: hostname, Resolved: false}
	}
}

func (e *Enumerator) OpenDict() {
	dict, err := os.Open(e.DictPath)
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		os.Exit(1)
	}
	e.dict = dict
}

func (e *Enumerator) Scan() {
	e.OpenDict()
	e.GetWildcardIPs()
	defer e.dict.Close()
	defer close(e.ch)

	dictScanner := bufio.NewScanner(e.dict)
	for dictScanner.Scan() {
		hostname := dictScanner.Text() + "." + e.Target
		e.waitGroup.Add(1)
		go e.Check(hostname)
	}

	e.waitGroup.Wait()
}

func (e *Enumerator) GetResults() <-chan Result {
	return e.ch
}

func NewEnumerator(target string, dict_path string) Enumerator {
	return Enumerator{
		Target:   target,
		DictPath: dict_path,
		ch:       make(chan Result, 16),
	}
}

func main() {
	dict_path := flag.String("d", "", "subdomains dictionary file")
	target := flag.String("t", "", "target domain")
	flag.Parse()

	if flag.NFlag() < 2 {
		flag.Usage()
		return
	}

	enumerator := NewEnumerator(*target, *dict_path)

	go enumerator.Scan()

	for result := range enumerator.GetResults() {
		fmt.Printf("\r[*] %s\u001b[0J", result.Hostname)
		if result.Resolved {
			fmt.Printf("\r[+] %s\u001b[0J\n", result.Hostname)
		}
	}

	fmt.Println("\rDone\u001b[0J")
}
