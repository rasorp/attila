// Copyright (c) James Rasell
// SPDX-License-Identifier: Apache-2.0

package file

import (
	"github.com/hashicorp/go-cty-funcs/cidr"
	"github.com/hashicorp/go-cty-funcs/crypto"
	"github.com/hashicorp/go-cty-funcs/encoding"
	"github.com/hashicorp/go-cty-funcs/filesystem"
	"github.com/hashicorp/go-cty-funcs/uuid"
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/ext/tryfunc"
	"github.com/hashicorp/hcl/v2/ext/typeexpr"
	ctyyaml "github.com/zclconf/go-cty-yaml"
	"github.com/zclconf/go-cty/cty/function"
	"github.com/zclconf/go-cty/cty/function/stdlib"
)

func hclEvalCtx(path string) *hcl.EvalContext {
	return &hcl.EvalContext{
		Functions: hclFunctions(path),
	}
}

func hclFunctions(path string) map[string]function.Function {
	return map[string]function.Function{
		"abs":             stdlib.AbsoluteFunc,
		"abspath":         filesystem.AbsPathFunc,
		"base64decode":    encoding.Base64DecodeFunc,
		"base64encode":    encoding.Base64EncodeFunc,
		"basename":        filesystem.BasenameFunc,
		"bcrypt":          crypto.BcryptFunc,
		"can":             tryfunc.CanFunc,
		"ceil":            stdlib.CeilFunc,
		"chomp":           stdlib.ChompFunc,
		"chunklist":       stdlib.ChunklistFunc,
		"cidrhost":        cidr.HostFunc,
		"cidrnetmask":     cidr.NetmaskFunc,
		"cidrsubnet":      cidr.SubnetFunc,
		"cidrsubnets":     cidr.SubnetsFunc,
		"coalesce":        stdlib.CoalesceFunc,
		"coalescelist":    stdlib.CoalesceListFunc,
		"compact":         stdlib.CompactFunc,
		"concat":          stdlib.ConcatFunc,
		"contains":        stdlib.ContainsFunc,
		"convert":         typeexpr.ConvertFunc,
		"csvdecode":       stdlib.CSVDecodeFunc,
		"dirname":         filesystem.DirnameFunc,
		"distinct":        stdlib.DistinctFunc,
		"element":         stdlib.ElementFunc,
		"flatten":         stdlib.FlattenFunc,
		"floor":           stdlib.FloorFunc,
		"format":          stdlib.FormatFunc,
		"formatdate":      stdlib.FormatDateFunc,
		"formatlist":      stdlib.FormatListFunc,
		"indent":          stdlib.IndentFunc,
		"index":           stdlib.IndexFunc,
		"join":            stdlib.JoinFunc,
		"jsondecode":      stdlib.JSONDecodeFunc,
		"jsonencode":      stdlib.JSONEncodeFunc,
		"keys":            stdlib.KeysFunc,
		"length":          stdlib.LengthFunc,
		"log":             stdlib.LogFunc,
		"lookup":          stdlib.LookupFunc,
		"lower":           stdlib.LowerFunc,
		"max":             stdlib.MaxFunc,
		"md5":             crypto.Md5Func,
		"merge":           stdlib.MergeFunc,
		"min":             stdlib.MinFunc,
		"pathexpand":      filesystem.PathExpandFunc,
		"parseint":        stdlib.ParseIntFunc,
		"pow":             stdlib.PowFunc,
		"range":           stdlib.RangeFunc,
		"reverse":         stdlib.ReverseFunc,
		"replace":         stdlib.ReplaceFunc,
		"regex_replace":   stdlib.RegexReplaceFunc,
		"rsadecrypt":      crypto.RsaDecryptFunc,
		"setintersection": stdlib.SetIntersectionFunc,
		"setproduct":      stdlib.SetProductFunc,
		"setunion":        stdlib.SetUnionFunc,
		"sha1":            crypto.Sha1Func,
		"sha256":          crypto.Sha256Func,
		"sha512":          crypto.Sha512Func,
		"signum":          stdlib.SignumFunc,
		"slice":           stdlib.SliceFunc,
		"sort":            stdlib.SortFunc,
		"split":           stdlib.SplitFunc,
		"strlen":          stdlib.StrlenFunc,
		"strrev":          stdlib.ReverseFunc,
		"substr":          stdlib.SubstrFunc,
		"timeadd":         stdlib.TimeAddFunc,
		"title":           stdlib.TitleFunc,
		"trim":            stdlib.TrimFunc,
		"trimprefix":      stdlib.TrimPrefixFunc,
		"trimspace":       stdlib.TrimSpaceFunc,
		"trimsuffix":      stdlib.TrimSuffixFunc,
		"try":             tryfunc.TryFunc,
		"upper":           stdlib.UpperFunc,
		"urlencode":       encoding.URLEncodeFunc,
		"uuidv4":          uuid.V4Func,
		"uuidv5":          uuid.V5Func,
		"values":          stdlib.ValuesFunc,
		"yamldecode":      ctyyaml.YAMLDecodeFunc,
		"yamlencode":      ctyyaml.YAMLEncodeFunc,
		"zipmap":          stdlib.ZipmapFunc,

		// filesystem calls
		"file":       filesystem.MakeFileFunc(path, false),
		"filebase64": filesystem.MakeFileFunc(path, true),
		"fileexists": filesystem.MakeFileExistsFunc(path),
		"fileset":    filesystem.MakeFileSetFunc(path),
	}
}
