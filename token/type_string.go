// Code generated by "stringer -type=Type -linecomment"; DO NOT EDIT.

package token

import "strconv"

const _Type_name = "TokIntegerTokIdentifier'=''==''<''>''+''-''*''/''if''else''while''('')''{''}'';''var'"

var _Type_index = [...]uint8{0, 10, 23, 26, 30, 33, 36, 39, 42, 45, 48, 52, 58, 65, 68, 71, 74, 77, 80, 85}

func (i Type) String() string {
	if i < 0 || i >= Type(len(_Type_index)-1) {
		return "Type(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _Type_name[_Type_index[i]:_Type_index[i+1]]
}
