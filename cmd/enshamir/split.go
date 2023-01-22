package main

import (
	"encoding/base64"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/ycdesu/enshamir"
)

type splitOptions struct {
	SecretFilePath string
	Parts          int
	Threshold      int
	OutputDir      string
}

func (o splitOptions) validate() error {
	if o.SecretFilePath == "" {
		return fmt.Errorf("source file path is not specified")
	}
	if o.OutputDir == "" {
		return fmt.Errorf("output directory is not specified")
	}

	return nil
}

var splitOpts splitOptions

var splitCmd = &cobra.Command{
	Use:   "split",
	Short: "encrypt and split the secret into shares and salt",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := splitOpts.validate(); err != nil {
			return err
		}

		pwd, err := enshamir.AskPassword()
		if err != nil {
			return err
		}

		saltDir := splitOpts.OutputDir
		sharesDir := filepath.Join(splitOpts.OutputDir, "shares")
		if err := os.MkdirAll(sharesDir, 0750); err != nil && !errors.Is(err, fs.ErrExist) {
			return fmt.Errorf("unable to create output directory %s:  %w", splitOpts.OutputDir, err)
		}

		secret, err := os.ReadFile(splitOpts.SecretFilePath)
		if err != nil {
			return fmt.Errorf("unable to read secret file %s: %w", splitOpts.SecretFilePath, err)
		}

		salt, shares, err := enshamir.EncryptSplit(pwd, secret, splitOpts.Parts, splitOpts.Threshold)
		if err != nil {
			return err
		}
		encodedSalt := []byte(base64.StdEncoding.EncodeToString(salt))

		saltFilePath := filepath.Join(saltDir, saltFileName)
		if err := enshamir.WriteIfNotExisted(saltFilePath, encodedSalt, 0600); err != nil {
			return err
		}

		for i, share := range shares {
			sharePath := filepath.Join(sharesDir, fmt.Sprintf("%s-%d", shareFilePrefix, i+1))
			encodedShare := []byte(base64.StdEncoding.EncodeToString(share))

			if err := enshamir.WriteIfNotExisted(sharePath, encodedShare, 0600); err != nil {
				return err
			}
		}

		fmt.Printf("\nThe secret is splitted into %d shares. The %s and %d of shares are required to reconstruct the secret.\n",
			splitOpts.Parts, saltFileName, splitOpts.Threshold)

		return nil
	},
}

func init() {
	rootCmd.AddCommand(splitCmd)

	splitCmd.Flags().StringVar(&splitOpts.SecretFilePath, "secret-file", "", "Read secret from the file")
	splitCmd.Flags().IntVar(&splitOpts.Parts, "parts", 3, "The number of shares to generate")
	splitCmd.Flags().IntVar(&splitOpts.Threshold, "threshold", 2, "The minimum number of shares required to reconstruct the secret")
	splitCmd.Flags().StringVar(&splitOpts.OutputDir, "output-dir", "", "The directory to save the shares")
}
