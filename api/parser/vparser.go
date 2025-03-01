package parser

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

type Version struct {
	Major int
	Minor int
	Patch int
	Tag   string
	Build string
	Date  string
	Error bool
}

type OSInfo struct {
	Os           string
	Architecture string
}

type LanguageInfo struct {
	Name string
	Version string
}

type ParsedInfo struct {
	Name     string
	Label    string
	Version  Version
	Os       OSInfo
	Language LanguageInfo
}

var reLanguage = regexp.MustCompile(`(?P<name>[a-zA-Z]+)?-?(?P<version>[\d+.?]+)`)

func (p *ParsedInfo) String() string {
	return fmt.Sprintf("%v (%v) %v %v", p.Name, p.Version, p.Os, p.Language)
}
func ParseVersionString(input string) *ParsedInfo {
	var output ParsedInfo

	// Quick fix to filter out enodes.
	if strings.Contains(input, "enode://") {
		return nil
	}

	// If there exists more than one version, don't count it
	if strings.Count(input, "/v") > 1 {
		return nil
	}

	// version string consists of four components, divided by /
	s := strings.Split(strings.ToLower(input), "/")
	output.Name = strings.ToLower(s[0])
	if output.Name == "" {
		return nil
	}

	// Check if the length of the 's' array is at least 4 before accessing its elements
	if len(s) >= 4 {
		if strings.HasPrefix(s[1], "v") {
			output.Version = parseVersion(s[1])
			output.Os = parseOS(s[2])
			output.Language = parseLanguage(s[3])
		} else {
			output.Label = s[1]
			output.Version = parseVersion(s[2])
			output.Os = parseOS(s[3])
			output.Language = parseLanguage(s[4])
		}
	} else {
		return nil
	}

	if output.Version.Error {
		fmt.Printf(" -> Error Parsing: '%s', %v\n", input, output)
		return nil
	}
	return &output
}



func parseLanguage(input string) LanguageInfo {
	var languageInfo LanguageInfo
	match := reLanguage.FindStringSubmatch(input)

	if len(match) > 0 {
		languageInfo.Name = match[reLanguage.SubexpIndex("name")]
		languageInfo.Version = match[reLanguage.SubexpIndex("version")]
	}

	return languageInfo
}

func parseVersion(input string) Version {
	var vers Version
	split := strings.Split(input, "-")
	split_length := len(split)
	switch len(split) {
	case 8:
		fallthrough
	case 7:
		fallthrough
	case 6:
		fallthrough
	case 5:
		vers.Date = split[split_length - 1]
		vers.Build = split[split_length - 2]
		vers.Tag = strings.Join(split[1:split_length - 3], "")
		vers.Major, vers.Minor, vers.Patch = parseVersionNumber(split[0])
	case 4:
		// Date
		vers.Date = split[3]
		fallthrough
	case 3:
		// Build
		vers.Build = split[2]
		fallthrough
	case 2:
		// Tag
		vers.Tag = split[1]
		fallthrough
	case 1:
		// Version
		vers.Major, vers.Minor, vers.Patch = parseVersionNumber(split[0])
	}

	if vers.Major == 0 && vers.Minor == 0 && vers.Patch == 0 {
		fmt.Println("Version string is invalid:", input)
		vers.Error = true
	}
	
	return vers
}

func parseVersionNumber(input string) (int, int, int) {
	// Version
	trimmed := strings.TrimLeft(input, "v")
	vSplit := strings.Split(trimmed, ".")
	var major, minor, patch int

	switch len(vSplit) {
	case 4:
		fallthrough
	case 3:
		patch, _ = strconv.Atoi(vSplit[2])
		fallthrough
	case 2:
		minor, _ = strconv.Atoi(vSplit[1])
		fallthrough
	case 1:
		major, _ = strconv.Atoi(vSplit[0])
	}

	return major, minor, patch
}

func parseOS(input string) OSInfo {
	reOSArch := regexp.MustCompile(`(?P<os>[a-zA-Z]+)(?:-(?P<arch>[a-zA-Z0-9_]+))?`)
	matches := reOSArch.FindStringSubmatch(input)

	var osInfo OSInfo
	if len(matches) > 0 {
		osInfo.Os = matches[reOSArch.SubexpIndex("os")]
		osInfo.Architecture = matches[reOSArch.SubexpIndex("arch")]
	}
	return osInfo
}
