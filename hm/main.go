package main

import (
	"errors"
	"flag"
	"fmt"
	"github.com/micanzhang/gohosts"
	"io/ioutil"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

const VERSION = "0.0.2"
const PLATFORM = runtime.GOOS

var pathMap map[string]string = map[string]string{
	"windows": "C:\\Windows\\System32\\drivers\\etc\\hosts",
	"darwin":  "/etc/hosts",
	"linux":   "/etc/hosts",
}

var (
	version bool
	edit    bool
	enable  bool
	disable bool
	list    bool
	domain  string
	ip      string
	help    bool
	group   string
)

func init() {
	flag.BoolVar(&version, "v", false, "Show Version Number")
	flag.BoolVar(&version, "version", false, "Show Version Number")
	flag.BoolVar(&edit, "e", false, "Edit Hosts File Directly")
	flag.BoolVar(&enable, "s", false, "Enable Hosts")
	flag.BoolVar(&disable, "r", false, "Disable Hosts")
	flag.BoolVar(&list, "l", false, "List All Hosts Config by Group Name")
	flag.StringVar(&domain, "d", "", "Filter By Domain")
	flag.StringVar(&ip, "i", "", "Filter by Ip Address")
	flag.BoolVar(&help, "h", false, "Usage")
	flag.StringVar(&group, "g", "", "Filter By Group Name")
}

func main() {
	flag.Parse()
	args := flag.Args()
	if len(args) > 0 {
		group = args[0]
	}
	if help {
		flag.Usage()
	}
	if version {
		printVersion()
	}
	if edit {
		editHosts()
	}
	params := make(map[string]string)
	if domain != "" {
		params["domain"] = domain
	}
	if ip != "" {
		params["ip"] = ip
	}
	if group != "" {
		params["group"] = group
	}
	if list {
		listHosts(params)
	}
	if enable || disable {
		switchHosts(params, enable || (disable && false))
	}
}

func usage() {
	fmt.Fprintf(os.Stderr, "usage: %s [options]\n", os.Args[0])
	os.Exit(2)
}

func printVersion() {
	fmt.Printf("hm version %s %s/%s\n", VERSION, PLATFORM, runtime.GOARCH)
	os.Exit(2)
}

func editHosts() {
	if hostFile, err := getHostFile(); err == nil {
		cmd := exec.Command("emacs", hostFile)
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			fmt.Println(err)
		}
	} else {
		fmt.Println(err)
	}
}

func listHosts(params map[string]string) {
	groups := getHost()
	if group, ok := params["group"]; ok {
		groups = groups.FindByName(group)
	}

	if ip, ok := params["ip"]; ok {
		groups = groups.FilterByIp(ip)
	}

	if domain, ok := params["domain"]; ok {
		groups = groups.FilterByDomain(domain)
	}

	fmt.Println(groups)
}

func switchHosts(params map[string]string, enable bool) {
	groups := getHost()
	hasParams := false
	if name, ok := params["group"]; ok {
		hasParams = true
		groups.SwitchByName(name, enable)
	}
	if ip, ok := params["ip"]; ok {
		hasParams = true
		groups.SwitchByIp(ip, enable)
	}
	if domain, ok := params["domain"]; ok {
		hasParams = true
		groups.SwitchByDomain(domain, enable)
	}
	if !hasParams {
		flag.Usage()
	} else if groups != nil {
		if err := updateHostString(groups.String()); err != nil {
			fmt.Printf("Error: %v\n", err)
		} else {
			fmt.Println("Hosts Switch Successfully!")
		}
	}
}

// Build Groups object
func getHost() *hosts.Groups {
	hostStr, err := loadHostString()
	if err != nil {
		panic(err)
	}
	lines := strings.Split(hostStr, "\n")
	groups := new(hosts.Groups)
	group := new(hosts.Group)
	inGroup := false
	for _, line := range lines {
		line = strings.TrimSpace(line)
		length := len(line)
		if length == 5 && inGroup {
			*groups = append(*groups, *group)
			group = new(hosts.Group)
			inGroup = false
		} else if length > 5 {
			if !inGroup && string(line[:5]) == "#====" {
				inGroup = true
				group.Name = strings.TrimSpace(string(line[5:]))
			} else if inGroup {
				host := new(hosts.Host)
				if line[0] == '#' {
					line = string(line[1:])
					host.Enable = false
				} else {
					host.Enable = true
				}
				items := strings.Fields(line)
				if len(items) > 0 {
					host.Ip = items[0]
					host.Domain = items[1:]
				}
				group.Items = append(group.Items, *host)
			}
		}
	}
	return groups
}

// Read host of supported platform
func loadHostString() (string, error) {
	if hostFile, err := getHostFile(); err == nil {
		bytes, err := ioutil.ReadFile(hostFile)
		if err != nil {
			return "", err
		}
		return string(bytes), nil
	} else {
		return "", err
	}

}

//write hosts config to system hosts config file
func updateHostString(hosts string) error {
	if hostFile, err := getHostFile(); err == nil {
		return ioutil.WriteFile(hostFile, []byte(hosts), 0777)
	} else {
		return err
	}
}

//Get host file by current os
func getHostFile() (string, error) {
	paltform := runtime.GOOS
	if hostFile, ok := pathMap[paltform]; ok {
		return hostFile, nil
	} else {
		return "", errors.New("unsupported PLATFORM!")
	}
}
