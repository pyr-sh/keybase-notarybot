package alice

import "os"

func parseMask(mask string) os.FileMode {
	if len(mask) != 10 {
		return 0
	}

	var res os.FileMode
	if mask[0] == 'd' {
		res = res | os.ModeDir
	} else if mask[0] == 'l' {
		res = res | os.ModeSymlink
	}
	if mask[1] == 'r' {
		res = res | 0400 // owner read perm
	}
	if mask[2] == 'w' {
		res = res | 0200 // owner write perm
	}
	if mask[3] == 'x' {
		res = res | 0100 // owner execute perm
	}
	if mask[4] == 'r' {
		res = res | 0040 // group read perm
	}
	if mask[5] == 'w' {
		res = res | 0020 // group write perm
	}
	if mask[6] == 'x' {
		res = res | 0010 // group execute perm
	}
	if mask[7] == 'r' {
		res = res | 0004 // others read perm
	}
	if mask[8] == 'w' {
		res = res | 0002 // others write perm
	}
	if mask[9] == 'x' {
		res = res | 0001 // others execute perm
	}
	return res
}
