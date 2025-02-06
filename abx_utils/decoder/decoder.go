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

package decoder

import (
	"better_known_installed_go/abx_utils/common"
	"bytes"
	"fmt"
)

type ABXDecoder struct {
	Input			*[]byte
	cursorPos		int
	internedStrings [][]byte
	elementStack	[]*common.XMLElement
}

func (d *ABXDecoder) IsABX() bool {
	d.cursorPos = 0
	
	readVal := d.readFromCurPos(4)
	if (bytes.Equal(readVal, common.StartMagic)) {
		// println("ABX v0 file found!")
		return true
	} else {
		return false
	}
}

func (d *ABXDecoder) Parse() (*common.XMLElement, error) {
	if (!d.IsABX()) {
		return nil, fmt.Errorf("this file is not ABX")
	}

	for {
		event := d.readFromCurPos(1)[0]
		token := event & 0x0f
		tType := event & 0xf0

		// println("Current token type", token, (tType >> 4))

		switch token {
		case common.TOKEN_ATTRIBUTE:
			attrName := d.readInternedString()
			// println("Start attribute", string(attrName))

			value := make([]byte, 0)
			// convertNumber := false

			switch tType {
			case common.DATA_NULL:
			case common.DATA_BOOLEAN_TRUE:
			case common.DATA_BOOLEAN_FALSE:
				// nop
			case common.DATA_STRING, common.DATA_BYTES_HEX, common.DATA_BYTES_BASE64:
				value = d.readString()
			case common.DATA_STRING_INTERNED:
				value = d.readInternedString()
			case common.DATA_INT, common.DATA_INT_HEX, common.DATA_FLOAT:
				value = d.readFromCurPos(4)
				// convertNumber = true
			case common.DATA_LONG, common.DATA_LONG_HEX, common.DATA_DOUBLE:
				value = d.readFromCurPos(8)
				// convertNumber = true
			default:
				return nil, fmt.Errorf("unexpected data type %v", tType >> 4)
			}

			/*
			if convertNumber {
				println("Attribute value", d.bytesToNumber(value), "len", len(value))
			} else {
				println("Attribute value", string(value), "len", len(value))
			}
			*/
			
			d.elementStack[len(d.elementStack) - 1].AddAttribute(string(attrName), tType, value)
		case common.TOKEN_START_DOCUMENT:
			// println("Start document")
		case common.TOKEN_END_DOCUMENT:
			// println("End document")
		case common.TOKEN_START_TAG:
			tagName := d.readInternedString()
			// println("Start tag", string(tagName))

			d.addElementToStack(&common.XMLElement{
				TagName: tagName,
			})
		case common.TOKEN_END_TAG:
			tagName := d.readInternedString()
			// println("End tag", string(tagName))
			lastTagName := d.elementStack[len(d.elementStack)-1].TagName

			// it shouldn't be happen
			if ! bytes.Equal(tagName, lastTagName) {
				println("Mismatching tags", string(tagName), "-", string(lastTagName))
			}

			if len(d.elementStack) == 1 {
				defer func () { d.elementStack = d.elementStack[: len(d.elementStack) - 1] }()
				return d.elementStack[0], nil
			}

			d.elementStack = d.elementStack[: len(d.elementStack) - 1]
		case common.TOKEN_TEXT, common.TOKEN_CDSECT, common.TOKEN_PROCESSING_INSTRUCTION, common.TOKEN_COMMENT, common.TOKEN_DOCDECL, common.TOKEN_IGNORABLE_WHITESPACE:
			readVal := d.readString()
			lastElement := d.elementStack[len(d.elementStack) - 1]
			lastElement.TextSections = append(
				lastElement.TextSections,
				&common.XMLTextSection{
					DataType: token,
					Text: &readVal,
				},
			)
		default:
			return nil, fmt.Errorf("unimplemented type %v %v", token, tType >> 4)
		}
	}
}

func (d *ABXDecoder) readShort() int16 {
	byteRead := d.readFromCurPos(2)

	result := int16(0)
    for i := 0; i < 2; i++ {
        result = result << 8
        result += int16(byteRead[i])
    }

    return result
}

/*
func (d *ABXDecoder) bytesToNumber(input []byte) int {
	result := int(0)
    for i := 0; i < len(input); i++ {
        result = result << 8
        result += int(input[i])
    }

    return result
}
*/

func (d *ABXDecoder) readString() []byte {
	length := d.readShort()
	return d.readFromCurPos(int(length))
}

func (d *ABXDecoder) readInternedString() []byte {
	index := d.readShort()
	if index < 0 {
		readStr := d.readString()
		d.internedStrings = append(d.internedStrings, readStr)
		return readStr
	}

	return d.internedStrings[index]
}

func (d *ABXDecoder) addElementToStack(element *common.XMLElement) {
	if len(d.elementStack) > 0 {
		lastElement := d.elementStack[len(d.elementStack) - 1]
		lastElement.SubElements = append(lastElement.SubElements, element)
	}

	d.elementStack = append(d.elementStack, element)
}

func (d *ABXDecoder) readFromCurPos(length int) []byte {
	defer func () { d.cursorPos += length }()
	return (*d.Input)[d.cursorPos : d.cursorPos + length]
}
