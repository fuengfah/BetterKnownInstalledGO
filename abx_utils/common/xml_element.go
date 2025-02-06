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

package common

import "bytes"

type XMLElement struct {
	TagName			[]byte
	Attributes		map[string]*XMLAttribute
	TextSections	[]*XMLTextSection
	SubElements		[]*XMLElement
}

func (e *XMLElement) FindElementsByTagName(name []byte) []*XMLElement {
	outputList := make([]*XMLElement, 0)

	if bytes.Equal(e.TagName, name) {
		outputList = append(outputList, e)
	}

	if len(e.SubElements) > 0 {
		for _, v := range e.SubElements {
			outputList = append(outputList, v.FindElementsByTagName(name)...)
		}
	}

	return outputList
}

func (e *XMLElement) FindAttributeByName(name string) *XMLAttribute {
	for k, v := range e.Attributes {
		if k == name {
			return v;
		}
	}

	return nil
}

func (e *XMLElement) AddAttribute(name string, dataType byte, value []byte) {
	if e.Attributes == nil {
		e.Attributes = make(map[string]*XMLAttribute)
	}

	e.Attributes[name] = &XMLAttribute{
		DataType: dataType,
		Value: &value,
	}
}

func (e *XMLElement) RemoveAttribute(name string) {
	delete(e.Attributes, name)
}
