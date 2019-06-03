package echo

import (
	"log"

	"github.com/ddosakura/starmap/gate/check"
	"github.com/labstack/echo"
)

func Example() {
	var c echo.Context
	if e := check.
		Build(c).
		Rules(check.M("a", "b").
			Or(check.M("c")),
			check.Rename("c", "C"),
			check.DefaultValue("d", "666"),
			check.Atoi("d"),
		).Load(nil); e != nil {
		log.Fatal(e)
	}
}
