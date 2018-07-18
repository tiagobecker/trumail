package verifier

import (
	"strings"
	"sync"
	"time"

	"github.com/sdwolfe32/httpclient"
)

// updateInterval is how often we should reach out to update
// the free address map
const updateIntervalforFree = 240 * time.Minute

// Free contains the map of known free email domains
type Free struct {
	client  *httpclient.Client
	freeMap *sync.Map
}

// NewFree creates a new Free and starts a domain farmer
// that retrieves all known free domains periodically
func NewFree(client *httpclient.Client) *Free {
	d := &Free{client, &sync.Map{}}
	go d.farmDomains(updateIntervalforFree)
	return d
}

// IsFree tests whether a string is among the known set of free
// mailbox domains. Returns true if the address is free
func (d *Free) IsFree(domain string) bool {
	_, ok := d.freeMap.Load(domain)
	return ok
}

// farmDomains retrieves new free domains every set interval
func (d *Free) farmDomains(interval time.Duration) error {
	for {
		for _, url := range listsfree {
			// Perform the request for the domain list
			body, err := d.client.GetString(url)
			if err != nil {
				continue
			}

			// Split
			for _, domain := range strings.Split(body, "\n") {
				d.freeMap.Store(strings.TrimSpace(domain), true)
			}
		}
		time.Sleep(interval)
	}
}

// list is a slice of free email address lists
var listsfree = []string{
	"https://raw.githubusercontent.com/tiagobecker/trumail/master/verifier/free.txt",	
}
