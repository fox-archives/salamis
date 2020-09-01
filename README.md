# sparta

Contextual vscode extension management

## Problem

-  Many extensions are installed, but only a fraction are used per project. Domain/Language specific extensions are still loaded, even if they may not be applicable
-  Enabling / disabling them per workspace is not a proper solution

## Solution

-  Extensions are categorized by domain/language/purpose
-  Extensions can be launched by any combination of these tags

Launch Command:

```sh
sparta launch --tag core --tag javascript
# is an abstraction over
code --extensions-dir extensions/core-javascript
```

Configuration File:

```toml
[[groups]]
name = "JavaScript"
description = "JavaScript / TypeScript Development Environment"
use = [
	"core", "javascript"
]

[[extensions]]
name = "editorconfig.editorconfig"
tags = [
"core"
]

[[extensions]]
name = "esbenp.prettier-vscode"
tags = [
"core"
]

[[extensions]]
name = "dbaeumer.vscode-eslint"
tags = [
"javascript"
]
```
