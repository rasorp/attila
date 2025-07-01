// Copyright (c) James Rasell
// SPDX-License-Identifier: Apache-2.0

package helper

import (
	"fmt"
	"net/url"
	"time"

	"github.com/ryanuber/columnize"

	"github.com/rasorp/attila/pkg/api"
)

func FormatList(in []string) string {
	columnConf := columnize.DefaultConfig()
	columnConf.Empty = "<none>"
	return columnize.Format(in, columnConf)
}

func FormatKV(in []string) string {
	columnConf := columnize.DefaultConfig()
	columnConf.Empty = "<none>"
	columnConf.Glue = " = "
	return columnize.Format(in, columnConf)
}

func FormatTime(t time.Time) string { return t.Format(time.RFC3339) }

func FormatError(cliMsg string, err error) string {

	var code int

	switch e := err.(type) {
	case *api.ResponseError:
		code = e.Code
	case *url.Error:
		code = 500
	default:
		code = 400
	}

	return FormatKV([]string{
		fmt.Sprintf("Description|%s", cliMsg),
		fmt.Sprintf("Error|%s", err),
		fmt.Sprintf("Code|%v", code),
	})
}
