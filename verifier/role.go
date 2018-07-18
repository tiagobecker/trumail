package verifier

import (
	"strings"
	"sync"
	"time"

	"github.com/sdwolfe32/httpclient"
)

// updateInterval is how often we should reach out to update
// the role address map
const updateIntervalforRole = 240 * time.Minute

// Role contains the map of known role users emails
type Role struct {
	client  *httpclient.Client
	roleMap *sync.Map
}

// NewRole creates a new Role and starts a domain farmer
// that retrieves all known role users emails periodically
func NewRole(client *httpclient.Client) *Role {
	d := &Role{client, &sync.Map{}}
	go d.farmUserName(updateIntervalforRole)
	return d
}

// IsRole tests whether a string is among the known set of roles
// users emails. Returns true if the address is role
func (d *Role) IsRole(username string) bool {
	_, ok := d.roleMap.Load(username)
	return ok
}

// farmUserName retrieves new role users emails every set interval
func (d *Role) farmUserName(interval time.Duration) error {
	for {
		for _, url := range listsrole {
			// Perform the request for the user email list
			body, err := d.client.GetString(url)
			if err != nil {
				continue
			}

			// Split
			for _, username := range strings.Split(body, "\n") {
				d.roleMap.Store(strings.TrimSpace(username), true)
			}
		}
		time.Sleep(interval)
	}
}

// list is a slice of role user email address lists
var listsrole = []string{
	"https://raw.githubusercontent.com/tiagobecker/trumail/master/verifier/role.txt",	
}
