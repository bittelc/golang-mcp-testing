# MCP Server in Go

This is a custom MCP (Model Context Protocol) server implementation written entirely in Go. I built this project for the fun of it, and also because I wanted full transparency and control over how my interactions and data are handled.

## Why Build This?

Rather than relying on closed-source packages from companies that aren't clear about how they use user interactions and data, I wanted to create something where I know exactly how the internals work. When a client connects to this server, I have complete visibility into what's happening under the hood.

## Implementation

This server is written entirely in Go and is based on the excellent work from:
- [gomcp library](https://github.com/localrivet/gomcp/)
- [Tutorial: How to Build Your Own MCP Vibe Coding Server in Go](https://medium.com/@alma.tuck/how-to-build-your-own-mcp-vibe-coding-server-in-go-using-gomcp-c80ad2e2377c)
