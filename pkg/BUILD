package(default_visibility = ["PUBLIC"])

go_library(
  name = "pkg",
  srcs = ["main.go", "plumbing.go", "porcelain.go", "util.go"],
  visibility = ["//..."],
  deps = [":go-toml"],
)

go_module(
  name = "go-toml",
  module = "github.com/pelletier/go-toml",
  version = "v1.8.0",
)
