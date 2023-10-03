package collector

import (
	"context"
	"fmt"
	"go/ast"
	"strconv"
)

type ArrayType struct {
	Length *int
	Values Expr
}

func (at ArrayType) String() (s string) {
	if at.Length == nil {
		s = fmt.Sprintf("[]%s", at.Values)
		goto end
	}
	s = fmt.Sprintf("[%d]%s", *at.Length, at.Values)
end:
	return s
}

func (c *Collector) CollectArrayType(ctx context.Context, aat *ast.ArrayType) (at ArrayType, err error) {
	if aat.Len != nil {
		//Value Verify this type assertion is the correct one
		strLen, err := c.CollectExprString(ctx, aat.Len)
		if err != nil {
			var length int
			length, err = strconv.Atoi(strLen)
			if err != nil {
				err = fmt.Errorf("unable to perform type conversion '%s' to int; %s",
					strLen,
					err.Error(),
				)
				goto end
			}
			at.Length = &length
		}
	}
	at.Values, err = c.CollectExpr(ctx, aat.Elt)
	if err != nil {
		goto end
	}
end:
	return at, err
}
