package bootstrap

import (
	"fmt"
	"io"
	"os"
)

const BinaryDstDir = "/usr/libexec/kubernetes/kubelet-plugins/credential-provider/exec/"

// Copies across the binaries defined in the Bootstrap.Binaries field.
// Each binary is expected to live in /bin and will be copied to the
// node's credential-provider exec directory. A channel is returned
// which can be used to monitor the progress of the copies.
func (b *Bootstrap) CopyBinaries() (chan error, error) {
	c := make(chan error)

	if err := os.MkdirAll(BinaryDstDir, 0755); err != nil {
		return nil, fmt.Errorf("Could not create binary directory: %s", err)
	}

	for _, bin := range b.Binaries {
		go func(binary string) {
			src, err := os.Open("/bin/" + binary)
			if err != nil {
				c <- fmt.Errorf("Could not open binary: %s", err)
				return
			}
			defer src.Close()

			dst, _ := os.Create(BinaryDstDir + binary)
			if err != nil {
				c <- fmt.Errorf("Could not create binary copy: %s", err)
				return
			}
			defer dst.Close()

			if _, err := io.Copy(dst, src); err != nil {
				c <- fmt.Errorf("Could not copy binary: %s", err)
				return
			}
			if err := dst.Chmod(0755); err != nil {
				c <- fmt.Errorf("Could not update permissions of binary: %s", err)
				return
			}
			c <- nil
		}(bin)
	}

	return c, nil
}

// Listens to the given channel and blocks until a message for each of
// the copied binaries has been received.
func (b *Bootstrap) WaitForBinaries(c chan error) {
	for range b.Binaries {
		if err := <-c; err != nil {
			fmt.Printf("%s\n", err)
			os.Exit(1)
		}
	}
}
