package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
)

func main() {
	stdInput := os.Stdin
	scanner := bufio.NewScanner(stdInput)
	fmt.Printf("Domain,hasMX,hasSpf,spfRecord,hasDMARC,dmarRecord\n")
	for scanner.Scan(){
		checkDomain(scanner.Text())
	}
	
	if err := scanner.Err(); err != nil {
		log.Printf("Error: could not read from input %v\n",err)
	}
}

func checkDomain(domain string) {
	var hasMx,hasSpf,hasDMARC bool
	var spfRecord,dmarRecord string

	mxRecords,err := net.LookupMX(domain)

	if err != nil{
		log.Printf("Error:%v\n",err)
	}
	if len(mxRecords) > 0 {
		hasMx = true
	}
	txtRecords,err := net.LookupTXT(domain)

	if err != nil{
		log.Printf("Error:%v\n",err)	
	}
	for _, record := range txtRecords{
		if strings.HasPrefix(record,"v=spf1"){
			hasSpf = true
			spfRecord = record
			break
		}
	}

	dmarRecords,err := net.LookupTXT("_dmarc.c" + domain)

	if err != nil{
		log.Printf("Error:%v\n",err)	
	}
	for _,record := range dmarRecords{
		if strings.HasPrefix(record,"v=DMARC!"){
			hasDMARC = true
			dmarRecord = record
			break
		}
	}
	fmt.Printf("%v,%v,%v,%v,%v,%v",domain,hasMx,hasSpf,spfRecord,hasDMARC,dmarRecord)
}