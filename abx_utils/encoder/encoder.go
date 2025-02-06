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

package encoder

import (
	"better_known_installed_go/abx_utils/common"
	"bytes"
)

type ABXEncoder struct {
	Root			*common.XMLElement
	internedStrings [][]byte
	output			[]byte
}

func (e *ABXEncoder) Parse() *[]byte {
	e.output = append(e.output, common.StartMagic...)
	e.output = append(e.output, common.TOKEN_START_DOCUMENT | common.DATA_NULL)
	e.parseElement(e.Root)
	e.output = append(e.output, common.TOKEN_END_DOCUMENT | common.DATA_NULL)
	return &e.output
}

func (e *ABXEncoder) parseElement(element *common.XMLElement) {
	e.output = append(e.output, common.TOKEN_START_TAG | common.DATA_STRING_INTERNED)
	e.writeInternedString(&element.TagName)

	if len(element.Attributes) > 0 {
		for k, v := range element.Attributes {
			e.output = append(e.output, common.TOKEN_ATTRIBUTE | v.DataType)
			kAsByte := []byte(k)
			e.writeInternedString(&kAsByte)

			switch v.DataType {
			case common.DATA_NULL, common.DATA_BOOLEAN_TRUE, common.DATA_BOOLEAN_FALSE:
				continue
			case common.DATA_STRING, common.DATA_BYTES_HEX, common.DATA_BYTES_BASE64:
				e.writeString(v.Value)
			case common.DATA_STRING_INTERNED:
				e.writeInternedString(v.Value)
			default:
				e.output = append(e.output, *v.Value...)
			}
		}
	}

	if len(element.TextSections) > 0 {
		for _, v := range element.TextSections {
			e.output = append(e.output, v.DataType | common.DATA_STRING)
			e.writeString(v.Text)
		}
	}

	if len(element.SubElements) > 0 {
		for _, v := range element.SubElements {
			e.parseElement(v)
		}
	}

	e.output = append(e.output, common.TOKEN_END_TAG | common.DATA_STRING_INTERNED)
	e.writeInternedString(&element.TagName)
}

func (e *ABXEncoder) writeShort(num int16) {
	e.output = append(e.output, uint8(num >> 8))
	e.output = append(e.output, uint8(num & 0xff))
}

func (e *ABXEncoder) writeString(str *[]byte) {
	e.writeShort(int16(len(*str)))
	e.output = append(e.output, *str...)
}

func (e *ABXEncoder) writeInternedString(str *[]byte) {
	intStrIndex := e.getInternedStringIndex(str)
	e.writeShort(intStrIndex)
	if intStrIndex < 0 {
		e.writeString(str)
	}
}

func (e *ABXEncoder) getInternedStringIndex(str *[]byte) int16 {
	for i, v := range e.internedStrings {
		if bytes.Equal(v, *str) {
			return int16(i)
		}
	}

	return -1
}
