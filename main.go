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

var dictPath, target string

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
	fakeHostname := "f4k3sd." + e.Target

	if ips, err := net.LookupIP(fakeHostname); err == nil {
		e.fakeIPs = ips
	}

	return e.fakeIPs
}

func (e *Enumerator) Check(hostname string) {
	defer e.waitGroup.Done()
	ips, err := net.LookupIP(hostname)
	e.ch <- Result{
		Hostname: hostname,
		Resolved: err == nil && !reflect.DeepEqual(e.fakeIPs, ips),
		IPs:      ips,
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

func (e *Enumerator) Results() <-chan Result {
	return e.ch
}

func NewEnumerator(target string, dictPath string) Enumerator {
	return Enumerator{
		Target:   target,
		DictPath: dictPath,
		ch:       make(chan Result),
	}
}

func init() {
	flag.StringVar(&dictPath, "d", "", "subdomains dictionary file")
	flag.StringVar(&target, "t", "", "target domain")
}

func main() {
	flag.Parse()

	if flag.NFlag() < 2 {
		flag.Usage()
		return
	}

	enumerator := NewEnumerator(target, dictPath)

	go enumerator.Scan()

	for result := range enumerator.Results() {
		fmt.Printf("\r[*] %s\u001b[0J", result.Hostname)
		if result.Resolved {
			fmt.Printf("\r[+] %s\u001b[0J\n", result.Hostname)
		}
	}

	fmt.Println("\rDone\u001b[0J")
}
