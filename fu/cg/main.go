package main

import (
	"github.com/weiyinfu/fugo/fu/cg/common"
	"github.com/weiyinfu/fugo/fu/cg/golang"
	"github.com/weiyinfu/fugo/fu/cg/ts"
)

/**
API生成工具:for golang gin
*/

func RunGoTs(idlSrcFile, gofile, goPackage, tsFile string) {
	conf := &common.Conf{
		SrcFile:   idlSrcFile,
		GoFile:    gofile,
		GoPackage: goPackage,
		TsFile:    tsFile,
	}
	api := common.Parse(conf.SrcFile)
	golang.Handle(api, conf)
	ts.Handle(api, conf)
}
