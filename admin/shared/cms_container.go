package shared

import (
	"fmt"
	"net/http"

	"github.com/dracory/hb"
)

func CMSContainer(r *http.Request, content string) string {
	return hb.Div().
		Class("cms").
		Style(fmt.Sprintf(
			"padding: %dpx %dpx %dpx %dpx;",
			PaddingTopPx(r),
			PaddingRightPx(r),
			PaddingBottomPx(r),
			PaddingLeftPx(r),
		)).
		Child(hb.NewHTML(content)).
		ToHTML()
}
