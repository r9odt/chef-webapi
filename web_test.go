package web

import (
	"fmt"
	"os"
	"testing"

	"github.com/JIexa24/chef-webapi/logging"

	. "github.com/smartystreets/goconvey/convey"
)

func TestConcatenateStringWithDelimeter(t *testing.T) {
	Convey("Give empty delimeter", t, func() {
		var d = ""
		var names = []string{"a", "b"}
		result := ConcatenateStringWithDelimeter(d, names)
		So(result, ShouldEqual, "ab")
	})

	Convey("Give empty delimeter and names", t, func() {
		var d = ""
		var names = []string{}
		result := ConcatenateStringWithDelimeter(d, names)
		So(result, ShouldEqual, "")
	})

	Convey("Give empty names", t, func() {
		var d = ","
		var names = []string{}
		result := ConcatenateStringWithDelimeter(d, names)
		So(result, ShouldEqual, "")
	})

	Convey("Give nil names", t, func() {
		var d = ","
		var names []string
		result := ConcatenateStringWithDelimeter(d, names)
		So(result, ShouldEqual, "")
	})

	Convey("Give non-empty names", t, func() {
		var d = ","
		var names = []string{"a", "b"}
		result := ConcatenateStringWithDelimeter(d, names)
		So(result, ShouldEqual, "a,b")
	})

	Convey("Give one name", t, func() {
		var d = ","
		var names = []string{"a"}
		result := ConcatenateStringWithDelimeter(d, names)
		So(result, ShouldEqual, "a")
	})

	Convey("Give three names", t, func() {
		var d = ","
		var names = []string{"a", "b", "c"}
		result := ConcatenateStringWithDelimeter(d, names)
		So(result, ShouldEqual, "a,b,c")
	})
}

func TestApp(t *testing.T) {
	Convey("Get new application", t, func() {
		NewApplication("development", nil)
		So(App, ShouldNotBeNil)
		So(App.Env, ShouldEqual, "development")
	})

	Convey("Configure LDAP", t, func() {
		App.ConfigureLDAP("testBase", "testBind", "prefix", "suffix")
		So(App, ShouldNotBeNil)
		So(App.LDAP, ShouldNotBeNil)
		So(App.LDAP.BaseDN, ShouldEqual, "testBase")
		So(App.LDAP.BindAddress, ShouldEqual, "testBind")
		So(App.LDAP.BindPrefix, ShouldEqual, "prefix")
		So(App.LDAP.BindSuffix, ShouldEqual, "suffix")
	})

	Convey("Configure logger", t, func() {
		l, err := logging.ConfigureLog("stdout", "info", "web")
		if err != nil {
			fmt.Fprintf(os.Stderr, "Can not configure log: %s\n", err.Error())
			os.Exit(1)
		}
		App.ConfigureLogger(l)
		So(App, ShouldNotBeNil)
		So(App.Logger, ShouldNotBeNil)
		So(App.Logger, ShouldEqual, l)
	})

	Convey("Configure Database", t, func() {
		err := App.ConfigureDatabase("test", "session", "name", "user", "password", "host", "123")
		So(err, ShouldNotBeNil)
		So(App, ShouldNotBeNil)
		So(App.DB, ShouldBeNil)
	})

	Convey("Configure App", t, func() {
		err := App.ConfigureDatabase("test", "session", "name", "user", "password", "host", "123")
		So(err, ShouldNotBeNil)
		So(App, ShouldNotBeNil)
		So(App.DB, ShouldBeNil)
	})

	Convey("Configure App", t, func() {
		App.ConfigureApp("tests", "tests/test.key", "tests/test.key", 0)
		So(App.AppKey, ShouldNotBeNil)
		So(App.DB, ShouldBeNil)
	})

	Convey("Configure Chef", t, func() {
		err := App.ConfigureChefClient("tests", "testurl", "tests/test1.key")
		So(App.Client, ShouldBeNil)
		So(err, ShouldNotBeNil)

		err = App.ConfigureChefClient("tests", "testurl", "tests/test.key")
		So(App.Client, ShouldNotBeNil)
		So(err, ShouldBeNil)
		app := App.GetChefClientConfig()
		So(app, ShouldEqual, App.Client)
	})
}
