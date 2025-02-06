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

var StartMagic = []byte{'A', 'B', 'X', 0}

const TOKEN_START_DOCUMENT = 0
const TOKEN_END_DOCUMENT = 1
const TOKEN_START_TAG = 2
const TOKEN_END_TAG = 3
const TOKEN_TEXT = 4
const TOKEN_CDSECT = 5
// const TOKEN_ENTITY_REF = 6
const TOKEN_IGNORABLE_WHITESPACE = 7
const TOKEN_PROCESSING_INSTRUCTION = 8
const TOKEN_COMMENT = 9
const TOKEN_DOCDECL = 10
const TOKEN_ATTRIBUTE = 15

const DATA_NULL = 1 << 4
const DATA_STRING = 2 << 4
const DATA_STRING_INTERNED = 3 << 4
const DATA_BYTES_HEX = 4 << 4
const DATA_BYTES_BASE64 = 5 << 4
const DATA_INT = 6 << 4
const DATA_INT_HEX = 7 << 4
const DATA_LONG = 8 << 4
const DATA_LONG_HEX = 9 << 4
const DATA_FLOAT = 10 << 4
const DATA_DOUBLE = 11 << 4
const DATA_BOOLEAN_TRUE = 12 << 4
const DATA_BOOLEAN_FALSE = 13 << 4
