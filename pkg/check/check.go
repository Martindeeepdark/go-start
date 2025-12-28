package check

import (
	"fmt"
	"os/exec"
	"regexp"
	"strings"
)

// GoVersionInfo holds Go version information
type GoVersionInfo struct {
	Version  string
	Major    int
	Minor    int
	Patch    int
	Valid    bool
	Warnings []string
}

// CheckGoVersion checks if the current Go version is compatible
func CheckGoVersion() *GoVersionInfo {
	info := &GoVersionInfo{
		Warnings: []string{},
	}

	// Get Go version
	cmd := exec.Command("go", "version")
	output, err := cmd.CombinedOutput()
	if err != nil {
		info.Valid = false
		info.Warnings = append(info.Warnings, "❌ 无法获取 Go 版本信息")
		return info
	}

	versionStr := string(output)
	info.Version = strings.TrimSpace(versionStr)

	// Parse version (e.g., "go version go1.21.0 darwin/amd64")
	re := regexp.MustCompile(`go version go(\d+)\.(\d+)\.(\d+)`)
	matches := re.FindStringSubmatch(versionStr)

	if len(matches) < 4 {
		info.Valid = false
		info.Warnings = append(info.Warnings, "❌ 无法解析 Go 版本")
		return info
	}

	info.Major = parseInt(matches[1])
	info.Minor = parseInt(matches[2])
	info.Patch = parseInt(matches[3])

	// Check version compatibility
	// Recommended: 1.21 - 1.23
	// Warning for 1.24+ (known GORM Gen compatibility issues)
	if info.Minor < 21 {
		info.Valid = false
		info.Warnings = append(info.Warnings,
			fmt.Sprintf("❌ Go 版本过低: %s", info.Version),
			"   需要 Go 1.21 或更高版本",
			"   请升级 Go: https://go.dev/dl/",
		)
		return info
	}

	if info.Minor >= 24 {
		info.Valid = true
		info.Warnings = append(info.Warnings,
			fmt.Sprintf("⚠️  检测到 Go %d.%d", info.Minor, info.Patch),
			"   Go 1.24+ 与 GORM Gen 存在已知兼容性问题",
			"   建议使用 Go 1.21-1.23 以获得最佳体验",
			"   或在生成代码时使用: GOTOOLCHAIN=local go1.21 start gen db",
		)
		return info
	}

	// Perfect version (1.21-1.23)
	info.Valid = true
	return info
}

func parseInt(s string) int {
	result := 0
	for _, c := range s {
		result = result*10 + int(c-'0')
	}
	return result
}

// PrintVersionInfo prints the version check result
func PrintVersionInfo(info *GoVersionInfo) {
	if info.Valid {
		if len(info.Warnings) > 0 {
			// Has warnings (e.g., Go 1.24+)
			fmt.Println("📋 Go 版本检查:")
			for _, warning := range info.Warnings {
				fmt.Println(warning)
			}
			fmt.Println()
		} else {
			// Perfect version
			fmt.Printf("✅ Go 版本检查通过: %s\n\n", info.Version)
		}
	} else {
		// Invalid version
		fmt.Println("📋 Go 版本检查:")
		for _, warning := range info.Warnings {
			fmt.Println(warning)
		}
		fmt.Println()
	}
}
