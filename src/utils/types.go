package utils

import "fmt"

type HostResponse struct {
	ArpCache       []*ArpCache       `json:"arp_cache,omitempty"`
	CurrentRam     *CurrentRam       `json:"current_ram,omitempty"`
	DiskPartitions []*DiskPartitions `json:"disk_partitions,omitempty"`
	GeneralInfo    *GeneralInfo      `json:"general_info,omitempty"`
	NeedsUpgrades  []*NeedsUpgrade   `json:"needs_upgrades,omitempty"`
	Host           *ConfigHost       `json:"host"`
}

func (r *HostResponse) Print(short bool) {
	if r == nil {
		return
	}

	if short {
		r.printShort()
		return
	}
	r.printLong()
}

func (r *HostResponse) printShort() {
	if len(r.NeedsUpgrades) > 0 {
		printShortNeedsUpgrade(r.NeedsUpgrades)
		fmt.Println("")
	}
}

func (r *HostResponse) printLong() {
	if len(r.ArpCache) > 0 {
		printLongArpCache(r.ArpCache)
		fmt.Println("")
	}
	if r.CurrentRam != nil {
		printLongCurrentRam(r.CurrentRam)
		fmt.Println("")
	}
	if len(r.DiskPartitions) > 0 {
		printLongDiskPartitions(r.DiskPartitions)
		fmt.Println("")
	}
	if r.GeneralInfo != nil {
		printLongGeneralInfo(r.GeneralInfo)
		fmt.Println("")
	}
	if len(r.NeedsUpgrades) > 0 {
		printLongNeedsUpgrade(r.NeedsUpgrades)
		fmt.Println("")
	}
}

type ArpCache struct {
	Addr   string `json:"addr"`
	HwType string `json:"hw_type"`
	HwAddr string `json:"hw_addr"`
	Mask   string `json:"mask"`
}

func printLongArpCache(a []*ArpCache) {
	fmt.Println("Arp Cache:")
	for _, arp := range a {
		fmt.Printf("Address: %s\n", arp.Addr)
		fmt.Printf("HwType: %s\n", arp.HwType)
		fmt.Printf("HwAddr: %s\n", arp.HwAddr)
		fmt.Printf("Mask: %s\n\n", arp.Mask)
	}
}

type CurrentRam struct {
	Total     float64 `json:"total"`
	Used      float64 `json:"used"`
	Available float64 `json:"available"`
}

func printLongCurrentRam(c *CurrentRam) {
	fmt.Println("Current Ram:")
	fmt.Printf("Total: %.2fMB\n", c.Total)
	fmt.Printf("Used: %.2fMB\n", c.Used)
	fmt.Printf("Available: %.2fMB\n", c.Available)
}

type DiskPartitions struct {
	FileSystem  string `json:"file_system"`
	Size        string `json:"size"`
	Used        string `json:"used"`
	Avail       string `json:"avail"`
	UsedPercent string `json:"used%"`
	Mounted     string `json:"mounted"`
}

func printLongDiskPartitions(d []*DiskPartitions) {
	fmt.Println("Disk Partitions:")
	for _, partition := range d {
		fmt.Printf("Filesystem: %s\n", partition.FileSystem)
		fmt.Printf("Size: %s\n", partition.Size)
		fmt.Printf("Used: %s\n", partition.Used)
		fmt.Printf("Available: %s\n", partition.Avail)
		fmt.Printf("Used %%: %s\n", partition.UsedPercent)
		fmt.Printf("Mount Point: %s\n\n", partition.Mounted)
	}
}

type GeneralInfo struct {
	OS         string `json:"os"`
	Kernel     string `json:"kernel"`
	Hostname   string `json:"hostname"`
	Uptime     string `json:"uptime"`
	ServerTime string `json:"server_time"`
}

func printLongGeneralInfo(g *GeneralInfo) {
	fmt.Println("General Info:")
	fmt.Printf("OS: %s\n", g.OS)
	fmt.Printf("Kernel: %s\n", g.Kernel)
	fmt.Printf("Hostname: %s\n", g.Hostname)
	fmt.Printf("Uptime: %s\n", g.Uptime)
	fmt.Printf("Server Time: %s\n", g.ServerTime)
}

type NeedsUpgrade struct {
	Package   string `json:"package"`
	Installed string `json:"installed"`
	Available string `json:"available"`
}

func printShortNeedsUpgrade(u []*NeedsUpgrade) {
	fmt.Println("Upgrades Available:")
	fmt.Printf("Total: %d\n", len(u))
}

func printLongNeedsUpgrade(u []*NeedsUpgrade) {
	fmt.Println("Upgrades Available:")
	fmt.Printf("Total: %d\n\n", len(u))
	for _, p := range u {
		fmt.Printf("Package: %s\n", p.Package)
		fmt.Printf("Installed: %s\n", p.Installed)
		fmt.Printf("Available: %s\n\n", p.Available)
	}
}
