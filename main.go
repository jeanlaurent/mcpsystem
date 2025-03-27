package main

import (
	"context"
	"fmt"
	"log"
	"os/exec"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func main() {
	s := server.NewMCPServer("OSXSystem", "0.0.1")

	tool := mcp.NewTool("battery", mcp.WithDescription("Get the battery status"))

	s.AddTool(tool, buildMCPHandler("pmset", "-g", "batt"))

	tool = mcp.NewTool("disk_space", mcp.WithDescription("Get the disk space"))
	s.AddTool(tool, buildMCPHandler("df", "-h"))

	tool = mcp.NewTool("cpu_usage", mcp.WithDescription("Get Cpu Usage"))
	s.AddTool(tool, buildMCPHandler("top", "-l", "1"))

	tool = mcp.NewTool("memory_usage", mcp.WithDescription("Get the memory usage"))
	s.AddTool(tool, buildMCPHandler("vm_stat"))

	if err := server.ServeStdio(s); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}

func buildMCPHandler(command string, args ...string) func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		output, err := runCommand(ctx, command, args...)

		if err != nil {
			result := mcp.NewToolResultText(fmt.Sprintf("Failed to execute %s command: %s %s with output %s", command, args, err, string(output)))
			result.IsError = true
			return result, err
		}

		return mcp.NewToolResultText(output), nil
	}
}

func runCommand(ctx context.Context, command string, args ...string) (string, error) {
	cmd := exec.CommandContext(ctx, command, args...)

	output, err := cmd.CombinedOutput()

	if err != nil {
		return "", err
	}

	return string(output), nil
}
