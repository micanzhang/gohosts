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
	flag.BoolVar(&edit, "edit", false, "Edit Hosts File Directly")
	flag.StringVar(&domain, "d", "", "Filter By Domain")
	flag.StringVar(&domain, "domain", "", "Filter By Domain")
	flag.StringVar(&group, "g", "", "Filter By Group Name")
	flag.StringVar(&group, "group", "", "Filter By Group Name")
	flag.BoolVar(&help, "h", false, "Usage Infomtion")
	flag.BoolVar(&help, "help", false, "Usage Infomation")
	flag.StringVar(&ip, "i", "", "Filter by Ip Address")
	flag.StringVar(&ip, "ip", "", "Filter by Ip Address")
	flag.BoolVar(&list, "l", false, "List All Hosts Config by Group Name")
	flag.BoolVar(&list, "list", false, "List All Hosts Config by Group Name")
	flag.BoolVar(&enable, "s", false, "Enable Hosts")
	flag.BoolVar(&enable, "enable", false, "Enable Hosts")
	flag.BoolVar(&disable, "r", false, "Disable Hosts")
	flag.BoolVar(&disable, "disable", false, "Disable Hosts")
}

func main() {
	flag.Usage = usage
	flag.Parse()
	args := flag.Args()
	var arg string
	if len(args) > 0 {
		arg = args[0]
		if group == "" {
			group = arg
		}
	}
	if help {
		flag.Usage()
	}
	if version {
		printVersion()
	}
	if edit {
		editHosts(arg)
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
	flag.Usage()
}

func printVersion() {
	fmt.Printf("hm version %s %s/%s\n", VERSION, PLATFORM, runtime.GOARCH)
	os.Exit(2)
}

func editHosts(editor string) {
	if editor == "" {
		editor = "emacs"
	}
	if hostFile, err := getHostFile(); err == nil {
		cmd := exec.Command(editor, hostFile)
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			fmt.Println(err)
		}
	} else {
		fmt.Println(err)
	}
	os.Exit(2)
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
	os.Exit(2)
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
	os.Exit(2)
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

var usageinfo string = `hm is a command line tool for hosts manager.

Usage:

	hm [flags] 
	
flags:
  -e, -edit                 Edit Hosts File Directly
  -d, -domain               Filter By Domain
  -g, -group                Filter By Group Name
  -h, -help                 Usage Infomation
  -i, -ip                   Filter by Ip Address
  -l, -list                 List All Hosts Config by Group Name
  -r, -disable              Disable Hosts
  -s, -enable               Enable Hosts
  -v, -version              Show Version Number 

Example:
    
	hm -l default
	
more help information please refer to https://github.com/micanzhang/gohosts	
`

func usage() {
	fmt.Println(usageinfo)
	os.Exit(2)
}
