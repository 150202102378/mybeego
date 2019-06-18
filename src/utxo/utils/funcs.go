package utils

import (
	"fmt"
	"os/exec"
	"regexp"
	"time"
)

//GetProjectPath is get path of jobName
func GetProjectPath(jobName string) string {
	cmd := exec.Command("sh", "-c", "pwd")
	out, _ := cmd.Output()
	goal := fmt.Sprintf("(.*%s/)", jobName)
	r, _ := regexp.Compile(goal)
	return r.FindString(string(out))
}

//FormatTimeToDateTime : format time to datetime str
func FormatTimeToDateTime(t time.Time) string {
	return t.Format(TimeFormatStr)
}
