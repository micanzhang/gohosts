package hosts

import (
	"fmt"
	"strings"
)

type Host struct {
	Ip     string
	Domain []string
	Enable bool
}

func (host Host) String() string {
	ip := ""
	if host.Enable {
		ip = host.Ip
	} else {
		ip = fmt.Sprintf("%s%s", "#", host.Ip)
	}
	return fmt.Sprintf("%-17s%s", ip, strings.Join(host.Domain, " "))
}

// Check host is none or not
func (host Host) IsNone() bool {
	return len(host.Ip) == 0 && host.IsEmpty()
}

// Check host has domain or not
func (host Host) IsEmpty() bool {
	return len(host.Domain) == 0
}

// Filter by domain, if not exists return nil
func (host Host) FilterByDomain(domain string) *Host {
	h := new(Host)
	for _, d := range host.Domain {
		if d == domain {
			if h.IsNone() {
				h.Ip = host.Ip
			}
			h.Domain = append(h.Domain, domain)
			return h
		}
	}
	return nil
}

// Filter by ip
func (host Host) FilterByIp(ip string) *Host {
	if host.Ip == ip {
		return &host
	} else {
		return nil
	}
}

// Group hosts in the same group for control easily
type Group struct {
	Name  string
	Items []Host
}

func (group Group) String() string {
	hosts := []string{}
	for _, host := range group.Items {
		hosts = append(hosts, host.String())
	}
	return fmt.Sprintf("Group: %s\n\t%s\n", group.Name, strings.Join(hosts, "\n\t"))
}

// Check group is none or not
func (group Group) IsNone() bool {
	return len(group.Name) == 0 && group.IsEmpty()
}

// Check group has host or not
func (group Group) IsEmpty() bool {
	return len(group.Items) == 0
}

// Filter group by group name
func (group Group) FilterByDomain(domain string) *Group {
	nGroup := new(Group)
	for _, h := range group.Items {
		host := h.FilterByDomain(domain)
		if host == nil {
			continue
		}
		if nGroup.IsNone() {
			nGroup.Name = group.Name
		}
		nGroup.Items = append(nGroup.Items, *host)
	}

	if nGroup.IsNone() {
		return nil
	}
	return nGroup
}

// Filter group by ip
func (group Group) FilterByIp(ip string) *Group {
	nGroup := new(Group)
	for _, h := range group.Items {
		host := h.FilterByIp(ip)
		if host == nil {
			continue
		}
		if nGroup.IsNone() {
			nGroup.Name = group.Name
		}
		nGroup.Items = append(nGroup.Items, *host)
	}

	if nGroup.IsNone() {
		return nil
	}
	return nGroup
}

func (group Group) Switch(enable bool) {
	for i := 0; i < len(group.Items); i++ {
		group.Items[i].Enable = enable
	}
}

func (group Group) SwitchByIp(ip string, enable bool) {
	for i := 0; i < len(group.Items); i++ {
		if group.Items[i].Ip == ip {
			group.Items[i].Enable = enable
		}
	}
}

func (group Group) SwitchByDomain(domain string, enable bool) {
	for i := 0; i < len(group.Items); i++ {
		if group.Items[i].IsEmpty() {
			break
		}
		for _, d := range group.Items[i].Domain {
			if d == domain {
				group.Items[i].Enable = enable
				break
			}
		}
	}
}

type Groups []Group

func (groups Groups) String() string {
	groupStrs := []string{}
	for _, group := range groups {
		groupStrs = append(groupStrs, group.String())
	}
	return strings.Join(groupStrs, "\n")
}

// Find groups by group name , if not exists return nil
func (groups Groups) FindByName(name string) *Groups {
	for _, group := range groups {
		if group.Name == name {
			gs := new(Groups)
			*gs = append(*gs, group)
			return gs
		}
	}
	return nil
}

// Filter by domain
func (groups Groups) FilterByDomain(domain string) *Groups {
	nGroups := new(Groups)
	for _, group := range groups {
		nGroup := group.FilterByDomain(domain)
		if nGroup != nil {
			*nGroups = append(*nGroups, *nGroup)
		}
	}
	return nGroups
}

// Filter by ip
func (groups Groups) FilterByIp(ip string) *Groups {
	nGroups := new(Groups)
	for _, group := range groups {
		nGroup := group.FilterByIp(ip)
		if nGroup != nil {
			*nGroups = append(*nGroups, *nGroup)
		}
	}
	return nGroups
}

func (groups *Groups) Disable() {
	for _, group := range *groups {
		group.Switch(false)
	}
}

func (groups *Groups) SwitchByName(name string, enable bool) {
	fmt.Println(name)
	for _, group := range *groups {
		if group.Name == name {
			group.Switch(enable)
			break
		}
	}
}

func (groups *Groups) SwitchByIp(ip string, enable bool) {
	for _, group := range *groups {
		group.SwitchByIp(ip, enable)
	}
}

func (groups *Groups) SwitchByDomain(domain string, enable bool) {
	for _, group := range *groups {
		group.SwitchByDomain(domain, enable)
	}
}
