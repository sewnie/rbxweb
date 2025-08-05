# rbxweb
[pkg.go.dev]:     https://pkg.go.dev/github.com/sewnie/rbxweb
[pkg.go.dev_img]: https://img.shields.io/badge/%E2%80%8B-reference-007d9c?logo=go&logoColor=white&style=flat-square

[![Godoc Reference][pkg.go.dev_img]][pkg.go.dev]

Go package that provides access to hand-picked Roblox web-based APIs.

The currently implemented APIs are considered legacy and are based on cookies;
they are suffixed with their version numbers to standout from rbxweb source code.

No tests are performed, and stability is not guranteed; this API is susceptible to breaking changes from both Roblox and code changes.

#### Example

```
client := rbxweb.NewClient()
cv, err := client.ClientSettingsV2.GetClientVersion(rbxweb.BinaryTypeWindowsPlayer, "")
```