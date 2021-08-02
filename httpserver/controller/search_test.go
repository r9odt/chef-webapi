package controller

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestConstructSearchORFilter(t *testing.T) {
	Convey("Give empty sign", t, func() {
		var names = []string{"a", "b"}
		result := ConstructSearchORFilter("", names)
		So(result, ShouldEqual, "")
	})
	
	Convey("Give empty names", t, func() {
		var names = []string{}
		result := ConstructSearchORFilter("sign", names)
		So(result, ShouldEqual, "")
	})
	
	Convey("Give nil names", t, func() {
		var names []string
		result := ConstructSearchORFilter("sign", names)
		So(result, ShouldEqual, "")
	})
	
	Convey("Give non-empty names", t, func() {
		var names = []string{"a", "b"}
		result := ConstructSearchORFilter("sign", names)
		So(result, ShouldEqual, "sign:a OR sign:b")
	})
	
	Convey("Give one name", t, func() {
		var names = []string{"a"}
		result := ConstructSearchORFilter("sign", names)
		So(result, ShouldEqual, "sign:a")
	})
	
	Convey("Give three names", t, func() {
		var names = []string{"a", "b", "c"}
		result := ConstructSearchORFilter("sign", names)
		So(result, ShouldEqual, "sign:a OR sign:b OR sign:c")
	})
}
