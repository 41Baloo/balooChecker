package main

import (
	"balooChecker/proxy"
	"bufio"
	"fmt"
	"os"
	"strconv"
	"sync"
	"sync/atomic"
	"time"
)

var (
	sem           chan struct{}
	checking      int32
	caughtUp      bool
	mutex         sync.Mutex
	numProxies    int
	definedTimout time.Duration
	definedOutput string
)

// Usage ./main [timeout] [output] [threads]
func main() {

	sArgs := os.Args

	if len(sArgs) < 4 {
		fmt.Println("[ Usage ]: ./main [proxy timeout] [output file] [max threads]")
		os.Exit(0)
	}

	go stats()

	tTimeout, timeoutErr := strconv.Atoi(sArgs[1])
	if timeoutErr != nil {
		fmt.Println("[ Error ]: " + timeoutErr.Error())
		os.Exit(0)
	}

	definedTimout = time.Duration(tTimeout) * time.Second
	definedOutput = os.Args[2]

	tMaxThreads, threadsErr := strconv.Atoi(sArgs[3])
	if threadsErr != nil {
		fmt.Println("[ Error ]: " + threadsErr.Error())
		os.Exit(0)
	}
	sem = make(chan struct{}, tMaxThreads)

	writeToFile(definedOutput+".txt", "[ Output ]: "+time.Now().Format("15:04:05"))

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		caughtUp = false
		proxyAddr := scanner.Text()
		sem <- struct{}{} // Wait for an open spot
		go func(addr string) {
			atomic.AddInt32(&checking, 1)
			checkProxy(addr, definedTimout)
			<-sem // Signal we're done
			atomic.AddInt32(&checking, -1)
		}(proxyAddr)
	}

	caughtUp = true

	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "Error reading input:", err)
	}

	time.Sleep(1 * time.Second)
	for checking != 0 {
		time.Sleep(1 * time.Second)
	}

	fmt.Println("[ Found " + fmt.Sprint(numProxies) + " Proxies ]!")
}

// Debug
func stats() {
	for {
		fmt.Println("[ Proxies Found ]: " + fmt.Sprint(numProxies) + " [ HTTP ]: " + fmt.Sprint(proxy.HTTP_FOUND) + " [ HTTPS ]: " + fmt.Sprint(proxy.HTTPS_FOUND) + " [ SOCKS4 ]: " + fmt.Sprint(proxy.SOCKS4_FOUND) + " [ SOCKS5 ]: " + fmt.Sprint(proxy.SOCKS5_FOUND) + " [ Checking ]: " + fmt.Sprint(checking) + " [ Caught Up ]: " + fmt.Sprint(caughtUp))
		time.Sleep(1 * time.Minute)
	}
}

func checkProxy(proxyAddr string, proxyTimeout time.Duration) {
	respHTTP, errHTTP := proxy.ConnectHTTP(proxyAddr, proxyTimeout)
	if errHTTP == nil {
		if proxy.ValidateResponse(respHTTP) {
			fmt.Println("[ Found ]: " + proxyAddr + " ( HTTP )")
			addProxyToList(proxyAddr, proxy.PROXY_HTTP)
		}
	}

	respSOCKS5, errSOCKS5 := proxy.ConnectSOCKS5(proxyAddr, proxyTimeout)
	if errSOCKS5 == nil {
		if proxy.ValidateResponse(respSOCKS5) {
			fmt.Println("[ Found ]: " + proxyAddr + " ( SOCKS5 )")
			addProxyToList(proxyAddr, proxy.PROXY_SOCKS5)
		}
	}

	respSOCKS4, errSOCKS4 := proxy.ConnectSOCKS4(proxyAddr, proxyTimeout)
	if errSOCKS4 == nil {
		if proxy.ValidateResponse(respSOCKS4) {
			fmt.Println("[ Found ]: " + proxyAddr + " ( SOCKS4 )")
			addProxyToList(proxyAddr, proxy.PROXY_SOCKS4)
		}
	}

	respHTTPS, errHTTPS := proxy.ConnectHTTPS(proxyAddr, proxyTimeout)
	if errHTTPS == nil {
		if proxy.ValidateResponse(respHTTPS) {
			fmt.Println("[ Found ]: " + proxyAddr + " ( HTTPS )")
			addProxyToList(proxyAddr, proxy.PROXY_HTTPS)
		}
	}

}

func addProxyToList(proxyAddr string, proxyType string) {

	switch proxyType {
	case proxy.PROXY_HTTP:
		atomic.AddInt32(&proxy.HTTP_FOUND, 1)
		break
	case proxy.PROXY_HTTPS:
		atomic.AddInt32(&proxy.HTTPS_FOUND, 1)
		break
	case proxy.PROXY_SOCKS4:
		atomic.AddInt32(&proxy.SOCKS4_FOUND, 1)
		break
	case proxy.PROXY_SOCKS5:
		atomic.AddInt32(&proxy.SOCKS5_FOUND, 1)
		break
	}

	// Ensure only 1 access at a time
	mutex.Lock()
	numProxies++
	writeToFile(definedOutput+".txt", proxyAddr)
	writeToFile(definedOutput+"_"+proxyType+".txt", proxyAddr)
	mutex.Unlock()
}

func writeToFile(fName string, str string) {
	logger, logErr := os.OpenFile(fName, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if logErr != nil {
		fmt.Println("[ Error ]: " + logErr.Error())
	}
	_, err := logger.WriteString(str + "\n")
	if err != nil {
		fmt.Println("[ Error ]: " + err.Error())
	}
	defer logger.Close()
}
