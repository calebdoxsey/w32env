package main

import (
	"flag"
	"log"
	"strings"
)

func add(root HKEY, path string, key string, value string) {
	v := get(root, path, key)
	parts := strings.Split(v, ";")
	if len(parts) == 1 && parts[0] == "" {
		parts = []string{}
	}
	for _, p := range parts {
		if p == value {
			return
		}
	}
	parts = append(parts, value)
	set(root, path, key, strings.Join(parts, ";"))
}
func remove(root HKEY, path string, key string, value string) {
	v := get(root, path, key)
	parts := strings.Split(v, ";")
	if len(parts) == 1 && parts[0] == "" {
		parts = []string{}
	}
	for i := 0; i < len(parts); i++ {
		if parts[i] == value {
			parts = append(parts[:i], parts[i+1:]...)
		}
	}
	set(root, path, key, strings.Join(parts, ";"))
}

func main() {
	flag.Parse()
	
	args := flag.Args()
	if len(args) < 1 {
		log.Fatalln("Expected Mode")
	}
	
	switch args[0] {
	case "add":
		if len(args) < 3 {
			log.Fatalln("Expected environment name and value")
		}
		
		add(HKEY_CURRENT_USER, "Environment", args[1], args[2])
	case "delete":
		if len(args) < 2 {
			log.Fatalln("Expected environment name")
		}
		
		delete(HKEY_CURRENT_USER, "Environment", args[1])		
	case "remove":
		if len(args) < 3 {
			log.Fatalln("Expected environment name and value")
		}
		
		remove(HKEY_CURRENT_USER, "Environment", args[1], args[2])
	case "set":
		if len(args) < 3 {
			log.Fatalln("Expected environment name and value")
		}
		
		set(HKEY_CURRENT_USER, "Environment", args[1], args[2])
	}
}
