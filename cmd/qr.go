/*
 * This file is part of PaperCrypt.
 *
 * PaperCrypt lets you prepare encrypted messages for printing on paper.
 * Copyright (C) 2023-2024 TMUniversal <me@tmuniversal.eu>.
 *
 * PaperCrypt is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Affero General Public License as published
 * by the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Affero General Public License for more details.
 *
 * You should have received a copy of the GNU Affero General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 */

package cmd

import (
	"encoding/json"
	"errors"
	"image"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/caarlos0/log"
	"github.com/makiuchi-d/gozxing"
	"github.com/makiuchi-d/gozxing/qrcode"
	"github.com/spf13/cobra"
	"github.com/tmuniversal/papercrypt/internal"
)

var (
	qrCmdFromJSON = false
	qrCmdToJSON   = false
)

type versionContainer struct {
	// Version should contain the semver version of PaperCrypt used to generate the document
	Version string `json:"Version"`
}

type dataContainer struct {
	// Data should contain the document data
	Data string `json:"Data"`
}

type dataContainerV1 struct {
	// Data should contain the document data, nested in a Data object (from crypto.PGPMessage),
	// this is used for backwards compatibility with PaperCrypt v1
	Data dataContainer `json:"Data"`
}

// qrCmd represents the data command
var qrCmd = &cobra.Command{
	Aliases:      []string{"q"},
	Args:         cobra.MaximumNArgs(1),
	SilenceUsage: true,
	Use:          "qr <input>",
	Short:        "Decode a document from a QR code.",
	Long: `Decode a document from a QR code.

This command allows you to decode data saved by PaperCrypt.
The QR code in a PaperCrypt document contains a JSON serialized object
that contains the encrypted data and the PaperCrypt metadata.

If you have trouble scanning the QR code with this command,
you may also try a QR code scanner app on your phone or tablet,
such as "Scandit" (https://apps.apple.com/de/app/scandit-barcode-scanner/id453880584
or https://play.google.com/store/apps/details?id=com.scandit.demoapp).
The resulting JSON data can be read by this command, by supplying the --json flag.
`,
	Example: `papercrypt qr ./qr.png | papercrypt decode -o ./out.json -P passphrase`,
	RunE: func(_ *cobra.Command, args []string) error {
		// 1. get data from either argument or inFileName
		if len(args) != 0 {
			inFileName = args[0]
		}

		inFile, err := internal.PrintInputAndGetReader(inFileName)
		if err != nil {
			return err
		}
		defer inFile.Close()

		var data []byte

		if qrCmdFromJSON {
			data, err = io.ReadAll(inFile)
			if err != nil && err != io.EOF {
				return errors.Join(errors.New("error reading input file"), err)
			}
		} else {
			img, _, err := image.Decode(inFile)
			if err != nil {
				return errors.Join(errors.New("error decoding image"), err)
			}

			if err := inFile.Close(); err != nil {
				return errors.Join(errors.New("error closing input file"), err)
			}

			bmp, err := gozxing.NewBinaryBitmapFromImage(img)
			if err != nil {
				return errors.Join(errors.New("error creating binary bitmap"), err)
			}

			qrReader := qrcode.NewQRCodeReader()
			result, err := qrReader.Decode(bmp, nil)
			if err != nil {
				return errors.Join(errors.New("error decoding QR code"), err)
			}

			data = []byte(result.GetText())
		}

		// 2. Open output file
		outFile, err := internal.GetFileHandleCarefully(outFileName, overrideOutFile)
		if err != nil {
			return err
		}
		defer func(file *os.File) {
			err := internal.CloseFileIfNotStd(file)
			if err != nil {
				log.WithError(err).Error("Error closing file")
			}
		}(outFile)

		if qrCmdToJSON {
			n, err := outFile.Write(data)
			if err != nil {
				return errors.Join(errors.New("error writing output"), err)
			}

			internal.PrintWrittenSize(n, outFile)
			return nil
		}

		// 3. Deserialize
		var output []byte
		var paperCryptMajorVersion uint32

		// decode version information or find .Data.Data (string)
		vc := versionContainer{}
		err = json.Unmarshal(data, &vc)
		if err != nil {
			dc := dataContainer{}
			err = json.Unmarshal(data, &dc)
			if err != nil {
				dcV1 := dataContainerV1{}
				err = json.Unmarshal(data, &dcV1)
				if err != nil {
					return errors.Join(errors.New("error deserializing data"), err)
				}
				paperCryptMajorVersion = 1
			} else {
				paperCryptMajorVersion = 2
			}
		} else {
			parseInt, err := strconv.ParseInt(strings.TrimPrefix(strings.Split(vc.Version, ".")[0], "v"), 10, 32)
			if err != nil {
				return errors.Join(errors.New("error parsing version"), err)
			}
			paperCryptMajorVersion = uint32(parseInt)
		}

		switch paperCryptMajorVersion {
		case 1:
			pc := internal.PaperCryptV1{}
			err = json.Unmarshal(data, &pc)
			if err != nil {
				return errors.Join(errors.New("error deserializing data"), err)
			}

			output, err = pc.GetText(false)
		case 2:
			pc := internal.PaperCrypt{}
			err = json.Unmarshal(data, &pc)
			if err != nil {
				return errors.Join(errors.New("error deserializing data"), err)
			}

			output, err = pc.GetText(false)
		}
		if err != nil {
			return errors.Join(errors.New("error deserializing data"), err)
		}

		// 6. Write to file
		n, err := outFile.Write(output)
		if err != nil {
			return errors.Join(errors.New("error writing output"), err)
		}

		internal.PrintWrittenSize(n, outFile)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(qrCmd)

	qrCmd.Flags().BoolVarP(&qrCmdFromJSON, "from-json", "j", false, "Read input from JSON instead of an image")
	qrCmd.Flags().BoolVarP(&qrCmdToJSON, "to-json", "J", false, "Write JSON output instead of plaintext, this cannot be used in the decode command (yet).")
}
