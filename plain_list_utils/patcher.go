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

package plainlistutils

import "strings"

func PatchPlainList(content *[]byte) []byte {
	plainList := string(*content)

	patchedList := ""
	for _, v := range strings.Split(plainList, "\n") {
		idx := strings.LastIndex(v, " ")
		packageName := v[idx + 1 :]

		if len(packageName) < 1 || packageName == "@system" {
			if (len(packageName) > 0) {
				patchedList += v
				patchedList += "\n"
			}
			continue
		}

		patchedList += v[: idx + 1]
		patchedList += "com.android.vending"
		patchedList += "\n"
	}

	return []byte(patchedList)
}
