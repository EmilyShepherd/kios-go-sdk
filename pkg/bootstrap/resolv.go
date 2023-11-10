package bootstrap

import (
	"fmt"
	"os"
)

const ResolvFilePath = "/etc/resolv.conf"

// Sometimes nameservers can be backed into an image, or set by another
// process (eg DHCP). Other times we may want to pre-prime these via our
// bootstrap container.
func (b *Bootstrap) SaveNameservers() error {
	nameservers := b.Provider.GetNameservers()

	if len(nameservers) == 0 {
		return nil
	}

	f, err := os.OpenFile(ResolvFilePath, os.O_TRUNC|os.O_CREATE|os.O_WRONLY, 0)
	if err != nil {
		return err
	}

	for _, ns := range nameservers {
		f.WriteString(fmt.Sprintf("nameserver %s\n", ns))
	}

	return f.Close()
}
