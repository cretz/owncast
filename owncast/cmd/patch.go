package cmd

import (
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/cretz/owncast/owncast/chrome"
	"github.com/cretz/owncast/owncast/log"
	"github.com/spf13/cobra"
)

func init() {
	patchCmd := &cobra.Command{
		Use:  "patch [path to chrome parent dir]",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			certDirAbs, err := filepath.Abs(certDir)
			if err != nil {
				return fmt.Errorf("Invalid cert dir: %v", err)
			}
			// Load existing root CA
			existingCA := chrome.LoadExistingRootCADERBytes()
			// Grab or create bytes to replace
			certFile := filepath.Join(certDirAbs, "ca.crt")
			certBytes, err := ioutil.ReadFile(certFile)
			if err != nil {
				if !os.IsNotExist(err) {
					return fmt.Errorf("Failed reading CA cert: %v", err)
				}
				log.Infof("Creating new root CA cert and saving as ca.crt and ca.key in %v", certDirAbs)
				kp, err := chrome.GenerateReplacementRootCA(len(existingCA), nil, nil)
				if err != nil {
					return fmt.Errorf("Unable to gen replacement root CA: %v", err)
				}
				if err = kp.PersistToFiles(certFile, filepath.Join(certDirAbs, "ca.key")); err != nil {
					return fmt.Errorf("Unable to persist replacement root CA: %v", err)
				}
				certBytes = kp.EncodeCertPEM()
			}
			certByteBlock, _ := pem.Decode(certBytes)
			// Find lib and patch
			lib, err := chrome.FindPatchableLib(args[0], existingCA)
			if err != nil {
				return err
			}
			log.Infof("Patching library at: %v", lib.Path())
			return lib.Patch(certByteBlock.Bytes)
		},
	}
	rootCmd.AddCommand(patchCmd)
}
