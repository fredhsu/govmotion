package main

import (
	"fmt"
	"net/url"

	"github.com/vmware/govmomi"
	// "github.com/vmware/govmomi/vim25/methods"
	"github.com/vmware/govmomi/vim25/mo"
	"github.com/vmware/govmomi/vim25/types"
)

func getHost(id int, crs []govmomi.Reference, c *govmomi.Client) (types.ManagedObjectReference, *types.ManagedObjectReference) {
	fmt.Println("getHost")
	var cr mo.ComputeResource
	fmt.Println("No crs: ", len(crs))
	// 4 = .146
	// 2 = .218
	// 7 = .150
	ref := crs[7]
	fmt.Println(ref)

	// err := c.Properties(ref.Reference(), []string{"host", "vm"}, &cr)
	err := c.Properties(ref.Reference(), nil, &cr)
	if err != nil {
		fmt.Println("Error:")
		fmt.Println(err.Error())
	}
	fmt.Println(cr.Host[0])
	fmt.Println(cr.ResourcePool)
	fmt.Println(cr.Summary)
	// h := cr.Host[0]
	// var hs mo.HostSystem
	// err = c.Properties(h, []string{"host"}, &hs)
	// fmt.Println(err)
	return cr.Host[0], cr.ResourcePool
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

	vms, err := folders.VmFolder.Children()

	if err != nil {
		fmt.Println("Error")
	}
	fmt.Println(vms)
	var vm mo.VirtualMachine
	for _, ref := range vms {
		err = c.Properties(ref.Reference(), []string{"config", "guest"}, &vm)
		if err != nil {
			fmt.Println("Not a vm")
		} else {
			if vm.Config.Name == "vxlan-vm1" {
				fmt.Println(vm.Config.Name)
				break
			}
		}
	}

	// host1, err := s.FindByIp(dc, "172.22.29.146", false)
	// var mh mo.HostSystem
	// h := host1.(*govmomi.HostSystem)
	// err = c.Properties(h.Reference(), []string{"parent"}, &mh)
	// fmt.Println(mh)

	// host2, err := si.FindByIp(dc, "172.22.29.146", false)
	// var hr mo.HostSystem
	// c.Properties(host1.Reference(), nil, &hr)
	// fmt.Println(hr)
	// fmt.Println(host2.Reference())
	// dstHost := mh.Reference()
	// mvmtask := types.MigrateVM_Task{
	// 	This:     vm.Reference(),
	// 	Host:     &dstHost, // *ManagedObjectReference
	// 	Priority: types.VirtualMachineMovePriorityDefaultPriority}
	// mvmresp, err := methods.MigrateVM_Task(c, &mvmtask)
	// fmt.Println(mvmresp)
	// fmt.Println(err)

	crs, _ := folders.HostFolder.Children()
	// fmt.Println(crs[0])
	dHost, pool := getHost(4, crs, c)
	mvmtask := types.MigrateVM_Task{
		This:     vm.Reference(),
		Pool:     pool,
		Host:     &dHost, // *ManagedObjectReference
		Priority: types.VirtualMachineMovePriorityDefaultPriority}
	// methods.MigrateVM_Task(c, &mvmtask)
	fmt.Println(mvmtask)
	// vMotion to .146 verify
	// vMotion to .150 verify
}
