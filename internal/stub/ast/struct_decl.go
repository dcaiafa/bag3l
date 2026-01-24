package ast

import (
	"github.com/dcaiafa/bag3l/internal/stub/analysis"
	"github.com/dcaiafa/bag3l/internal/token"
)

type StructDecl struct {
	Name   string
	Fields []*StructField
}

func (d *StructDecl) RunPass(ctx *Context, pass Pass) {
	RunPassChildren(ctx, d, d.Fields, pass)

	if pass == Check {
		structDecl := &analysis.Struct{
			Name: d.Name,
		}
		for _, fieldAST := range d.Fields {
			field := &analysis.StructField{
				Name: fieldAST.Name,
				Type: fieldAST.Type.Type,
			}
			structDecl.Fields = append(structDecl.Fields, field)
		}

		err := ctx.Analysis.AddStruct(structDecl)
		if err != nil {
			ctx.Failf(token.Pos{}, "%v", err)
			return
		}
	}
}
