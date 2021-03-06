package main

import (
	"strings"

	// "github.com/freeconf/yang/testdata"
	"github.com/freeconf/examples/basics"
	"github.com/freeconf/restconf"
	"github.com/freeconf/restconf/device"
	"github.com/freeconf/yang/source"
	//"www.xdn.kr/qkd"
)

func main() {

	// This is an example for a main() body to configure a typical application, in
	// this case a car.

	// Any typical application. Here we create a car.
	car := basics.New()

	// Where to looks for yang files, this tells library to use these
	// two relative paths. This is just one of many ways to control loading yang files.
	yangPath := source.Path("./yang:./testdata")

	// A device is just a container of modules.  Modules are independent services inside
	// your microservice.ß
	d := device.New(yangPath)
	// Here we are adding the "car" module to our microservice. "car" - the name of the
	// module causing car.yang to load from yang path.  You can add as many modules as you
	// want, here we're just adding one.
	if err := d.Add("abc", basics.Manage(car)); err != nil {
		panic(err)
	}

	// This is an example for a main() body to configure a typical application, in
	// this case a car.

	/*
		// Any typical application. Here we create a car.
		//car := basics.New()
		qkdNetworkTopology := qkd.Mew()

		// Where to looks for yang files, this tells library to use these
		// two relative paths. This is just one of many ways to control loading yang files.
		yangPath := source.Path("./yang:./testdata")

		// A device is just a container of modules.  Modules are independent services inside
		// your microservice.
		d := device.New(yangPath)
		// Here we are adding the "car" module to our microservice. "car" - the name of the
		// module causing car.yang to load from yang path.  You can add as many modules as you
		// want, here we're just adding one.
		if err := d.Add("qkdNetworkTopology", qkd.Manage(qkdNetworkTopology)); err != nil {
			panic(err)
		}
		//d.Add("qkdNetworkTopology", qkd.Manage(qkdNetworkTopology))
	*/

	// Adding RESTCONF protocol support.  Should you want an alternate protocol, IETF defines
	// other YANF-based protocols like NETCONF or COMI. While FreeCONF doesn't have support for
	// these protocols it could at some point.
	restconf.NewServer(d)

	// Even though the main configuration comes from the application management
	// system after call-home has registered this system it's often useful/neccessary
	// to bootstrap config for some of the local modules
	// Normally stored in a file, this is example "startup" configuration. Any additional
	// configuration can be applied *after" car starts, i.e. live configuration changes.
	staticStartupConfig := `
		{
			"fc-restconf" : {
				"debug": true,
				"web" : {
					"port" : ":8080"
				}
			},
			"abc" : {}
		}`
	if err := d.ApplyStartupConfig(strings.NewReader(staticStartupConfig)); err != nil {
		panic(err)
	}

	// trick to wait for cntrl-c... but you might have other services to start
	select {}
}
