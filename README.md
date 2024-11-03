# go-annotation
> An annotation format for Go

## Format

This package defines a toolchain for an annotation format, usable in most Go comments (only type, field and named
functions are implemented) to allow code generation, meta programming and similar functions, with the source as input,
in a standardized way.

Type documentations, function documentations (including receiver functions) and struct field documentations are being
considered.

**Example**

```go
package demo
// simple_annotation.name{attribute1,attribute2=with_value,attribute3="with quoted, and formatted value"}
// simple_annotation{attribute1,attribute2,attribute3}
func AFunction() {}

/**
 inspect output on "./example/demo":
{
	"functions": [
		{
			"function": {
				"file_path": "example/demo/b.go",
				"name": "AFunction",
				"package": "demo"
			},
			"annotations": [
				{
					"identifier": "simple_annotation",
					"arguments": {
						"another_attribute": "TRUE",
						"attribute": "TRUE",
						"third_attribute": "TRUE"
					}
				},
				{
					"identifier": "simple_annotation.name",
					"arguments": {
						"another_attribute": "with_value",
						"attribute": "TRUE",
						"third_attribute": "with quoted, and formatted value"
					}
				}
			]
		}
	]
}
 */
```

### General Form

The general form of an annotation is based on the definition of a generic Lua (table) expression. 

**Annotation**

```
Regex: ^([a-zA-Z\.\_]+)\s{0,1}\{(.*)\}$
```

An annotation MUST start with line start and end with line end, no multi-line allowed.
The annotation name MUST only consist of characters from a-z (lower and uppercase allowed) or underscore (_) or
period (.).
After the name an optional whitespace may be placed.
Following the name and optional whitespace an attributes definition block, wrapped in curly brackets ({...}), MUST be
present.

**Attribute**

```
Regex: ([a-zA-Z_]+)(=("([^"]*)"|([^,]*)))?,?
```

An attribute MUST be key-only or a key-value-pair, annotated with an equality sign (=).
An attribute declaration SHOULD end with a comma (,).
An attribute key can only consist of characters from a-z (lower and uppercase allowed) and underscores (_).
If an attribute is key-only, it's attribute value will be set to TRUE (string).
If an attributes value contains a comma (,) it MUST be quoted ("..."). 
If an attribute value is quoted ("...") then, no quote character (") is allowed in the value. There is no escaping.

## CLI

### inspect

```bash
$ go run ./cmd/inspect/... -root ./example/demo
```

Renders the JSON version of the annotation output of every type and function found by traversing the file tree starting
at root.