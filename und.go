package main

import (
	"flag"
	"log"
	"time"
)

var (
	nmKey      = flag.String("key", "12345", " To be replaced by your unique API key. Visit the API Manager page within your account for details.")
	nmDomain   = flag.String("domain", "namesilo.com", "The domain associated with the DNS resource record to modify")
	nmHost     = flag.String("host", "www", "The hostname to use (there is no need to include the \".DOMAIN\")")
	nmInterval = flag.Duration("interval", 300*time.Second, "The seconds of updating interval")
)

func main() {
	flag.Parse()
	updateDNSLoop()
}

func updateDNSLoop() {
	tick := time.NewTicker(*nmInterval)
	defer tick.Stop()
	for {
		select {
		case <-tick.C:
			err := doUpdateDNS(*nmDomain, *nmHost, *nmKey)
			if err != nil {
				log.Printf("[%v] update DNS record failed, domain:%s, host:%s, error:%v", time.Now().Format(time.RFC3339), *nmDomain, *nmHost, err)
			}
		}
	}
}

func doUpdateDNS(domain, host, key string) error {
	listResp, err := dnsList(domain, key)
	if err != nil {
		return err
	}

	// find the one need to be updated
	fullHost := host + "." + domain
	for _, item := range listResp.ListReply.DNSRecords {
		if item.Host == fullHost {
			if item.Value != listResp.Request.IP {
				// update record
				log.Printf("[%v] host: %s, update old IP: %s to new IP: %s", time.Now().Format(time.RFC3339), item.Host, item.Value, listResp.Request.IP)
				err = dnsUpdate(key, domain, item.RecordID, host, listResp.Request.IP)
				if err != nil {
					return err
				}
			} else {
				log.Printf("[%v] host: %s, old IP: %s, new IP: %s are same, nothing to do", time.Now().Format(time.RFC3339), item.Host, item.Value, listResp.Request.IP)
			}
			return nil
		}
	}
	return nil
}
