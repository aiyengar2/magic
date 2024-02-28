package cmd

import (
	"fmt"
	"sort"
	"strings"
)

const (
	aliasSource   = "internal:go:alias"
	warningSource = "internal:warning"
	finishSource  = "mage:finish"
)

func PrintFinish(fmtStr string, args ...any) {
	Print(finishSource, fmt.Sprintf(fmtStr, args...))
}

func PrintWarning(fmtStr string, args ...any) {
	Print(warningSource, fmt.Sprintf(fmtStr, args...))
}

func PrintGo(fmtStr string, args ...any) {
	Print(aliasSource, fmt.Sprintf(fmtStr, args...))
}

func Print(source string, cmd string, args ...string) {
	PrintWith(nil, source, cmd, args...)
}

func PrintWith(env map[string]string, source string, cmd string, args ...string) {
	envSlice := MapToEnvSlice(env)
	fullCmd := strings.Join(append(append(envSlice, cmd), args...), " ")
	if Verbose {
		fmt.Printf("+ (%s) %s\n", source, fullCmd)
	}
}

func MapToEnvSlice(m map[string]string) []string {
	var mSlice []string
	for k, v := range m {
		mSlice = append(mSlice, fmt.Sprintf("%s=%s", k, v))
	}
	sort.Strings(mSlice)
	return mSlice
}
