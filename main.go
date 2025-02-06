/*
 *  This file is part of better_known_installed_go.
 *
 *  better_known_installed_go is free software: you can redistribute it and/or modify
 *  it under the terms of the GNU General Public License as published by
 *  the Free Software Foundation, either version 3 of the License, or
 *  (at your option) any later version.
 *
 *  better_known_installed_go is distributed in the hope that it will be useful,
 *  but WITHOUT ANY WARRANTY; without even the implied warranty of
 *  MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 *  GNU General Public License for more details.
 *
 *  You should have received a copy of the GNU General Public License
 *   along with better_known_installed_go.  If not, see <https://www.gnu.org/licenses/>.
 */

package main

import (
	"better_known_installed_go/abx_utils/common"
	"better_known_installed_go/abx_utils/decoder"
	"better_known_installed_go/abx_utils/encoder"
	plainlistutils "better_known_installed_go/plain_list_utils"
	"bytes"
	"os"
	"os/exec"
	"strings"
)

const prefix = "/data/system/"
// const prefix = ""

const packagesListFile = prefix + "packages.list"
const packagesListBakFile = packagesListFile + ".bak"
const newPackagesListFile = prefix + "packages_new.list"

const packagesXMLFile = prefix + "packages.xml"
const packagesXMLBakFile = packagesXMLFile + ".bak"
const newPackagesXMLFile = prefix + "packages_new.xml"

const packagesWarnsXMLFile = prefix + "packages-warnings.xml"
const packagesWarnsXMLBakFile = packagesWarnsXMLFile + ".bak"
const newPackagesWarnsXMLFile = prefix + "packages-warnings_new.xml"

var playStorePkgName = []byte("com.android.vending")

func patchPackagesList() {
	plBuffer, err := os.ReadFile(packagesListFile)
	if err != nil {
		panic(err)
	}

	err = os.WriteFile(packagesListBakFile, plBuffer, 0644)
	if err != nil {
		panic(err)
	}

	patchedContent := plainlistutils.PatchPlainList(&plBuffer)
	err = os.WriteFile(newPackagesListFile, patchedContent, 0644)
	if err != nil {
		panic(err)
	}

	err = os.Rename(newPackagesListFile, packagesListFile)
	if err != nil {
		panic(err)
	}
}

func patchPackagesXML() {
	buffer, err := os.ReadFile(packagesXMLFile)
	if err != nil {
		panic(err)
	}

	err = os.WriteFile(packagesXMLBakFile, buffer, 0644)
	if err != nil {
		panic(err)
	}

	abxDecoder := decoder.ABXDecoder{ Input: &buffer }
	output, err := abxDecoder.Parse()
	if err != nil {
		panic(err)
	}

	println("root tag:", string((*output).TagName), "len:", len((*output).SubElements))

	userPackages := make([]*common.XMLElement, 0)
	var playStorePackage *common.XMLElement = nil
	for _, v := range output.FindElementsByTagName([]byte("package")) {
		codePath := v.FindAttributeByName("codePath")
		packageName := v.FindAttributeByName("name")

		if packageName != nil && bytes.Equal(*packageName.Value, playStorePkgName) {
			println("Found Play Store", string(*codePath.Value))
			playStorePackage = v
		}

		if codePath != nil && strings.Index(string(*codePath.Value), "/data/app/") == 0 {
			userPackages = append(userPackages, v)
		}
	}

	if playStorePackage == nil {
		panic("No Play Store found")
	}

	psUidPtr := playStorePackage.FindAttributeByName("userId")
	psUid := make([]byte, len(*psUidPtr.Value))
	copy(psUid, *psUidPtr.Value)

	for _, v := range userPackages {
		packageName := v.FindAttributeByName("name")
		if packageName == nil {
			continue
		}

		currentPkgName := string(*packageName.Value)
		for k, attr := range v.Attributes {
			switch k {
			case "installerUid":
				// println(k, "found for", currentPkgName)
				if bytes.Equal(psUid, *attr.Value) {
					continue
				}

				println(k, "will be patched for", currentPkgName)
				v.AddAttribute(k, attr.DataType, psUid)
			case "installer", "installInitiator":
				// println(k, "found for", currentPkgName)
				if bytes.Equal(playStorePkgName, *attr.Value) {
					continue
				}

				println(k, "will be patched for", currentPkgName)
				v.AddAttribute(k, attr.DataType, playStorePkgName)
			case "installOriginator":
				println(k, "will be removed for", currentPkgName)

				v.RemoveAttribute(k)
			}
		}
	}

	abxEncoder := encoder.ABXEncoder{ Root: output }
	output2 := abxEncoder.Parse()

	err = os.WriteFile(newPackagesXMLFile, *output2, 0644)
	if err != nil {
		panic(err)
	}

	err = os.Rename(newPackagesXMLFile, packagesXMLFile)
	if err != nil {
		panic(err)
	}
}

func patchPackagesWarningsXML() {
	warnsBuffer, err := os.ReadFile(packagesWarnsXMLFile)
	if err != nil {
		panic(err)
	}

	err = os.WriteFile(packagesWarnsXMLBakFile, warnsBuffer, 0644)
	if err != nil {
		panic(err)
	}

	pkgWarnings := common.XMLElement{ TagName: []byte("packages") }
	pkgWarnEncoder := encoder.ABXEncoder{ Root: &pkgWarnings }
	output3 := pkgWarnEncoder.Parse()

	err = os.WriteFile(newPackagesWarnsXMLFile, *output3, 0644)
	if err != nil {
		panic(err)
	}

	err = os.Rename(newPackagesWarnsXMLFile, packagesWarnsXMLFile)
	if err != nil {
		panic(err)
	}
}

func fixPermissions() {
	chownCmd := exec.Command("chown", "system:system", packagesListFile, packagesXMLFile, packagesWarnsXMLFile)
	output, err := chownCmd.CombinedOutput()
	if err != nil {
		print(output)
		panic(err)
	}

	chmodCmd := exec.Command("chmod", "640", packagesListFile, packagesXMLFile, packagesWarnsXMLFile)
	output, err = chmodCmd.CombinedOutput()
	if err != nil {
		print(output)
		panic(err)
	}

	restoreconCmd := exec.Command("restorecon", packagesListFile, packagesXMLFile, packagesWarnsXMLFile)
	output, err = restoreconCmd.CombinedOutput()
	if err != nil {
		print(output)
		panic(err)
	}
}

func main() {
	patchPackagesList()
	patchPackagesXML()
	patchPackagesWarningsXML()
	fixPermissions()

	println("Done!")
	println("You should reboot your phone to use the new package list.")
	println("You can take a backup before reboot (it's recommended).")
	println("Check", prefix, "folder for files with .bak suffix.")
}
