package main

import (
	"fmt"
	"github.com/guelfey/go.dbus"
	"github.com/guelfey/go.dbus/prop"
	"os"
)

type foo string

func (f foo) Foo() (string, *dbus.Error) {
	fmt.Println(f)
	return string(f), nil
}

func main() {
	conn, err := dbus.ConnectSessionBus()
	if err != nil {
		panic(err)
	}
	reply, err := conn.RequestName("com.github.guelfey.Demo",
		dbus.NameFlagDoNotQueue)
	if err != nil {
		panic(err)
	}
	if reply != dbus.RequestNameReplyPrimaryOwner {
		fmt.Fprintln(os.Stderr, "name already taken")
		os.Exit(1)
	}
	c := make(chan interface{})
	propsSpec := map[string]map[string]prop.Prop{
		"com.github.guelfey.Demo": map[string]prop.Prop{
			"SomeInt": prop.Prop{int32(0), true, c},
		},
	}
	props := prop.New(propsSpec)
	f := foo("Bar")
	conn.Export(f, "/com/github/guelfey/Demo", "com.github.guelfey.Demo")
	conn.Export(props, "/com/github/guelfey/Demo",
		"org.freedesktop.DBus.Properties")
	fmt.Println("Listening on com.github.guelfey.Demo / /com/github/guelfey/Demo ...")
	for v := range c {
		fmt.Println("SomeInt changed to", v)
	}
}
