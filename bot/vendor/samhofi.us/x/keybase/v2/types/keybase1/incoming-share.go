// Auto-generated to Go types using avdl-compiler v1.4.6 (https://github.com/keybase/node-avdl-compiler)
//   Input file: ../client/protocol/avdl/keybase1/incoming-share.avdl

package keybase1

import (
	"fmt"
)

type IncomingShareType int

const (
	IncomingShareType_FILE  IncomingShareType = 0
	IncomingShareType_TEXT  IncomingShareType = 1
	IncomingShareType_IMAGE IncomingShareType = 2
	IncomingShareType_VIDEO IncomingShareType = 3
)

func (o IncomingShareType) DeepCopy() IncomingShareType { return o }

var IncomingShareTypeMap = map[string]IncomingShareType{
	"FILE":  0,
	"TEXT":  1,
	"IMAGE": 2,
	"VIDEO": 3,
}

var IncomingShareTypeRevMap = map[IncomingShareType]string{
	0: "FILE",
	1: "TEXT",
	2: "IMAGE",
	3: "VIDEO",
}

func (e IncomingShareType) String() string {
	if v, ok := IncomingShareTypeRevMap[e]; ok {
		return v
	}
	return fmt.Sprintf("%v", int(e))
}

type IncomingShareItem struct {
	Type          IncomingShareType `codec:"type" json:"type"`
	OriginalPath  string            `codec:"originalPath" json:"originalPath"`
	OriginalSize  int               `codec:"originalSize" json:"originalSize"`
	ScaledPath    *string           `codec:"scaledPath,omitempty" json:"scaledPath,omitempty"`
	ScaledSize    *int              `codec:"scaledSize,omitempty" json:"scaledSize,omitempty"`
	ThumbnailPath *string           `codec:"thumbnailPath,omitempty" json:"thumbnailPath,omitempty"`
	Content       *string           `codec:"content,omitempty" json:"content,omitempty"`
}

func (o IncomingShareItem) DeepCopy() IncomingShareItem {
	return IncomingShareItem{
		Type:         o.Type.DeepCopy(),
		OriginalPath: o.OriginalPath,
		OriginalSize: o.OriginalSize,
		ScaledPath: (func(x *string) *string {
			if x == nil {
				return nil
			}
			tmp := (*x)
			return &tmp
		})(o.ScaledPath),
		ScaledSize: (func(x *int) *int {
			if x == nil {
				return nil
			}
			tmp := (*x)
			return &tmp
		})(o.ScaledSize),
		ThumbnailPath: (func(x *string) *string {
			if x == nil {
				return nil
			}
			tmp := (*x)
			return &tmp
		})(o.ThumbnailPath),
		Content: (func(x *string) *string {
			if x == nil {
				return nil
			}
			tmp := (*x)
			return &tmp
		})(o.Content),
	}
}
