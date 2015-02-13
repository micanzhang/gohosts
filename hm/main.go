package main

import (
	"errors"
	"flag"
	"fmt"
	"github.com/micanzhang/gohosts"
	"io/ioutil"
	"os"
	"runtime"
	"strings"
)

const VERSION = "1.0@dev"

func usage() {
	fmt.Fprintf(os.Stderr, "usage: %s [reddit]\n", os.Args[0])
	flag.PrintDefaults()
	os.Exit(2)
}

func help() {
	usage()
}

func version() {
	fmt.Printf("%s version %s\n", os.Args[0], VERSION)
	flag.PrintDefaults()
	os.Exit(2)
}

func main() {
	parseCMD()
}

func edit() {
	cmd := make(map[string]string)
	cmd["darwin"] = "open -a Emacs /etc/hosts"
	cmd["linux"] = "emacs /etc/hosts"
	cmd["win"] = "emasc C:/Windows/System32/driver/etc/hosts"
}

func list(params map[string]string) {
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

func enable(params map[string]string) {
	groups := getHost()
	enable := true
	if name, ok := params["group"]; ok {
		if ip, ok := params["ip"]; ok {
			groups.SwitchByIp(ip, enable)
		} else if domain, ok := params["domain"]; ok {
			groups.SwitchByDomain(domain, enable)
		} else {
			groups.SwitchByName(name, enable)
		}
	} else {

	}

	if ip, ok := params["ip"]; ok {
		groups.SwitchByIp(ip, enable)
	}

	if domain, ok := params["domain"]; ok {
		groups.SwitchByDomain(domain, enable)
	}

	fmt.Println(groups)
}

func disable(params map[string]string) {

}

func parseCMD() {
	args := os.Args
	params := make(map[string]string)
	if len(args) == 1 {
		usage()
	}

	key := ""
	action := ""
	//support list, edit, disable, enable, help, version
	for _, arg := range args[1:] {
		if len(key) > 0 {
			params[key] = arg
			key = ""
		} else {
			switch arg {
			case "-h", "--help":
				action = "help"
				help()
			case "-v", "--version":
				action = "version"
				version()
			case "-l", "--list", "list":
				if len(action) == 0 {
					action = "list"
				} else if action != "list" {
					usage()
				}
				break
			case "-e", "--edit":
				edit()
			case "-r", "--remove", "remove":
				if len(action) == 0 {
					action = "disable"
				} else if action != "disable" {
					usage()
				}
				break
			case "-s", "--switch", "switch":
				if len(action) == 0 {
					action = "enable"
				} else if action != "enable" {
					usage()
				}
				break
			case "-g", "--group":
				key = "group"
				break
			case "-d", "--domain":
				key = "domain"
				break
			case "-i", "--ip":
				key = "ip"
				break
			default:
				usage()
				break
			}
		}
	}

	if action == "disable" {
		disable(params)
	} else if action == "enable" {
		enable(params)
	} else {
		list(params)
	}
}

func listHost(params map[string]string) {
	groups := getHost()
	fmt.Println(groups)
	//list all
	//list by group
	// list by domain
	// list by ip
}

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

func loadHostString() (string, error) {
	hostPath := ""
	paltform := runtime.GOOS
	if paltform == "darwin" || paltform == "linux" {
		hostPath = "/etc/hosts"
	} else {
		return "", errors.New("unsupported platform!")
	}

	bytes, err := ioutil.ReadFile(hostPath)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}
