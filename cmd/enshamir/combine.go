package main

import (
	"encoding/base64"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/ycdesu/enshamir"
)

type combineOptions struct {
	SaltFilePath string
	SharesDir    string

	SecretFilePath string
}

var combineOpts combineOptions

var combineCmd = &cobra.Command{
	Use:   "combine",
	Short: "reconstruct and decrypt shares and salt to the secret",
	RunE: func(cmd *cobra.Command, args []string) error {
		password, err := enshamir.AskPassword()
		if err != nil {
			return err
		}

		encodedSalt, err := os.ReadFile(combineOpts.SaltFilePath)
		if err != nil {
			return fmt.Errorf("unable to read salt file %s: %w", combineOpts.SaltFilePath, err)
		}
		salt, err := base64.StdEncoding.DecodeString(string(encodedSalt))
		if err != nil {
			return fmt.Errorf("unable to decode salt file %s: %w", combineOpts.SaltFilePath, err)
		}

		sharesDir, err := os.ReadDir(combineOpts.SharesDir)
		if err != nil {
			return fmt.Errorf("unable to read shares directory %s: %w", combineOpts.SharesDir, err)
		}

		var shares [][]byte
		for _, s := range sharesDir {
			if s.Name() == saltFileName {
				continue
			}
			if s.IsDir() {
				continue
			}

			p := filepath.Join(combineOpts.SharesDir, s.Name())
			encodedShare, err := os.ReadFile(p)
			if err != nil {
				return fmt.Errorf("unable to read share file %s: %w", p, err)
			}
			s, err := base64.StdEncoding.DecodeString(string(encodedShare))
			if err != nil {
				return fmt.Errorf("unable to decode share file %s: %w", p, err)
			}
			shares = append(shares, s)
		}

		secret, err := enshamir.CombineDecrypt(password, salt, shares)
		if err != nil {
			return fmt.Errorf("unable to reconstruct and decrypt shares: %w", err)
		}

		if err := enshamir.WriteIfNotExisted(combineOpts.SecretFilePath, secret, 0600); err != nil {
			return fmt.Errorf("unable to write secret file %s: %w", combineOpts.SecretFilePath, err)
		}

		fmt.Println("\nsecret file is written to", combineOpts.SecretFilePath)

		return nil
	},
}

func init() {
	rootCmd.AddCommand(combineCmd)

	combineCmd.Flags().StringVar(&combineOpts.SaltFilePath, "salt-file", "", "salt file path")
	combineCmd.Flags().StringVar(&combineOpts.SharesDir, "shares-dir", "", "shares directory")
	combineCmd.Flags().StringVar(&combineOpts.SecretFilePath, "secret-file", "", "output secret file")
}
