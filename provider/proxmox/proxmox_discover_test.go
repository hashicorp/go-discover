package proxmox

import (
	discover "github.com/hashicorp/go-discover"
)

var _ discover.Provider = (*Provider)(nil)
