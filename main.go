package main

import (
	"fmt"
	"net/url"

	"github.com/vmware/govmomi"
	"github.com/vmware/govmomi/vim25/mo"
	"github.com/vmware/govmomi/vim25/types"
)

func findCrs(name string, crs []types.ManagedObjectReference, c Client) mo.ComputeResource {
	var cr mo.ComputeResource
	for _, r := range crs {
		err := c.Properties(crs.Reference(), []string{"host"}, &cr)
		if err == nil {
			if cr.Name == name {
				return cr
			}
		}
	}
	return cr
}

func main() {
	u, err := url.Parse("https://root:vmware@172.22.28.190/sdk")
	if err != nil {
		fmt.Println("URl parse error")
	}
	c, err := govmomi.NewClient(*u, true)
	if err != nil {
		fmt.Println("Connect error", err)
	}
	s := c.SearchIndex()
	ref, err := s.FindChild(c.RootFolder(), "SEDEMO")
	dc, ok := ref.(*govmomi.Datacenter)
	if !ok {
		fmt.Println("DC error")
	}
	folders, err := dc.Folders()
	if err != nil {
		fmt.Println("Error")
	}
	crs, err := folders.HostFolder.Children()
	if err != nil {
		fmt.Println("Error")
	}

	vms, err := folders.VmFolder.Children()

	if err != nil {
		fmt.Println("Error")
	}
	fmt.Println(vms)
	var vm mo.VirtualMachine
	for _, ref := range vms {
		err = c.Properties(ref.Reference(), []string{"config", "guest"}, &vm)
		//fmt.Printf("%+v\n", vm)
		if err != nil {
			fmt.Println("Not a vm")
		} else {
			fmt.Println(vm.Config.Name)
		}
	}
}
