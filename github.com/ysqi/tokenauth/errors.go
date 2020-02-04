// Copyright 2016 Author YuShuangqi. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package tokenauth

import (
	"fmt"
)

//Customer error.
type ValidationError struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

func (v ValidationError) Error() string {
	return fmt.Sprintf("%s:%s", v.Code, v.Msg)
}
