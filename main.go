package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strings"

	sh "github.com/cartersusi/script_helper"
)

const (
	MAX_ATTEMPTS = 10
)

var version string

func confirm_cmd(cmd string) bool {
	reader := bufio.NewReader(os.Stdin)
	attempts := 0

	for {
		if attempts >= MAX_ATTEMPTS {
			sh.Error("Too many attempts", true)
			return false
		} else if attempts > MAX_ATTEMPTS/2 {
			sh.Warning(fmt.Sprintf("You have %d attempts remaining", MAX_ATTEMPTS-attempts))
		}

		fmt.Printf("\n%s\n\tAre you sure you want to run this command? (y/n): ", cmd)
		input, err := reader.ReadString('\n')
		if err != nil {
			sh.Error(fmt.Sprintf("Error reading input: %s", err.Error()))
			return false
		}

		input = strings.TrimSpace(strings.ToLower(input))

		switch input {
		case "y", "yes":
			return true
		case "n", "no":
			return false
		default:
			attempts++
			fmt.Println("Invalid input. Please enter 'y' or 'n'.")
		}
	}

}

func main() {
	dir := flag.String("dir", "", "Directory to search")
	d := flag.String("d", "", "Directory to search")

	replace := flag.String("replace", "", "String to replace")
	r := flag.String("r", "", "String to replace")

	with := flag.String("with", "", "String to replace with")
	w := flag.String("w", "", "String to replace with")

	v := flag.Bool("v", false, "Print version")
	_version := flag.Bool("version", false, "Print version")

	flag.Parse()
	//version_flag := sh.GetFlag(_version, v, false, "version")
	if *v || *_version {
		fmt.Println("replace | Version:", version)
		return
	}

	directory := sh.GetFlag(dir, d, "", "directory")
	replaceString := sh.GetFlag(replace, r, "", "replace")
	withString := sh.GetFlag(with, w, "", "with")
	sh.Success(fmt.Sprintf("Replacing `%s` with `%s` in `%s`", replaceString, withString, directory))

	//grep -rl 'this-text' ./ | xargs -I {} sed -i '' 's/this-text/this-text/g' {}
	cmd := fmt.Sprintf("grep -rl '%s' %s | xargs -I {} sed -i '' 's/%s/%s/g' {}", replaceString, directory, replaceString, withString)
	if !confirm_cmd(cmd) {
		sh.Error("Bad Command. Exiting...", true)
	}

	out, err := exec.Command("sh", "-c", cmd).Output()
	if err != nil {
		sh.Error(err.Error(), true)
	}

	out_msg := fmt.Sprintf("DONE: Replaced `%s` with `%s` in `%s`", replaceString, withString, directory)
	if len(out) == 0 || string(out) == "" {
		sh.Success(out_msg)
	} else {
		sh.Success(string(out))
	}
}
