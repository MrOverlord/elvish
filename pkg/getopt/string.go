// Code generated by "stringer -type=Config,HasArg,ContextType -output=string.go"; DO NOT EDIT.

package getopt

import "strconv"

const (
	_Config_name_0 = "DoubleDashTerminatesOptionsFirstArgTerminatesOptions"
	_Config_name_1 = "LongOnly"
)

var (
	_Config_index_0 = [...]uint8{0, 27, 52}
)

func (i Config) String() string {
	switch {
	case 1 <= i && i <= 2:
		i -= 1
		return _Config_name_0[_Config_index_0[i]:_Config_index_0[i+1]]
	case i == 4:
		return _Config_name_1
	default:
		return "Config(" + strconv.FormatInt(int64(i), 10) + ")"
	}
}

const _HasArg_name = "NoArgumentRequiredArgumentOptionalArgument"

var _HasArg_index = [...]uint8{0, 10, 26, 42}

func (i HasArg) String() string {
	if i >= HasArg(len(_HasArg_index)-1) {
		return "HasArg(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _HasArg_name[_HasArg_index[i]:_HasArg_index[i+1]]
}

const _ContextType_name = "NewOptionOrArgumentNewOptionNewLongOptionLongOptionChainShortOptionOptionArgumentArgument"

var _ContextType_index = [...]uint8{0, 19, 28, 41, 51, 67, 81, 89}

func (i ContextType) String() string {
	if i >= ContextType(len(_ContextType_index)-1) {
		return "ContextType(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _ContextType_name[_ContextType_index[i]:_ContextType_index[i+1]]
}