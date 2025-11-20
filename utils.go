package appium_cli

import (
	"context"
	"fmt"
	"os/exec"
	"runtime"
	"strings"
	"time"
)

func (platform PlatformType) ToString() string {
	var ret string
	switch platform {
	case Android:
		ret = "Android"
	case IOS:
		ret = "IOS"
	case Mac:
		ret = "Mac"
	case Windows:
		ret = "Windows"
	}
	return ret
}

func GetOutPutString(commandShell string, commandList []string) (info string, error *AppiumError) {
	out, err := exec.Command(commandShell, commandList...).Output()
	if err != nil {
		error = &AppiumError{
			Message:   "Get shell output error",
			ErrorCode: OsShellError,
		}
		return
	}
	info = string(out)
	return
}

func NoOutPutString(commandShell string, commandList []string) (error *AppiumError) {
	_, err := exec.Command(commandShell, commandList...).Output()
	if err != nil {
		error = &AppiumError{
			Message:   "Get shell output error",
			ErrorCode: OsShellError,
		}
		return
	}
	return
}

func GetAdbOutputString(commandShell string, commandList []string) (error *AppiumError) {
	cmd := exec.Command(commandShell, commandList...)

	// 执行命令并忽略输出
	err := cmd.Start()
	if err != nil {
		error = &AppiumError{
			Message:   "Get shell output error",
			ErrorCode: OsShellError,
		}
		return
	}
	return
}

// KillLoopCmd
// @Note: this function can not kill subprocess, e,g
// "python3 main.py"
// 修复后的 KillLoopCmd
// 建议：不要返回自定义的 *AppiumError，内部工具函数返回标准 error 更加通用
func KillLoopCmd(commandShell string, commandList []string) (bool, error) {
	// 缩短超时时间，Ping 不需要 3 秒那么久
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	process := exec.CommandContext(ctx, commandShell, commandList...)

	// CombinedOutput 可以同时获取 stdout 和 stderr，便于调试
	outputBytes, err := process.CombinedOutput()
	outputStr := string(outputBytes)

	// 1. 如果 context 超时了
	if ctx.Err() == context.DeadlineExceeded {
		return false, fmt.Errorf("command timed out: %s", outputStr)
	}

	// 2. 如果命令执行出错 (比如 adb 没找到，或者 ping 失败)
	if err != nil {
		return false, fmt.Errorf("command failed: %v, output: %s", err, outputStr)
	}

	// 3. 简单的判断逻辑：如果有输出通常意味着执行了
	// 对于 Ping，我们通常判断是否包含 "ttl=" 或 "bytes from"
	if len(outputStr) > 0 && !strings.Contains(outputStr, "unknown host") {
		return true, nil
	}

	return false, nil
}

func GetAdbPath() string {
	if runtime.GOOS == "windows" {
		adbPath, err := exec.LookPath("adb")
		if err != nil {
			fmt.Println("找不到 adb 命令：", err)
			return "D:\\AndroidSDK\\android-sdk_r24.4.1-windows\\android-sdk_r24.4.1-windows\\android-sdk-windows\\platform-tools\\adb.exe"
		}
		return adbPath
	}
	return "adb"
}
