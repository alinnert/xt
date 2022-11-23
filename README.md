# XML Schema Tools

This is a command line tool that fetches and displays information about an XML Schema.

## Usage

```
xt ./path-to-schema.xsd element-name
```

Analyses the hierarchy of element definitions and gives all paths from root level elements to the provided `element-name` showing shorter paths first.

### Positional arguments

- **file**: The entry file of the schema
- **element-name**: The name of the element you want to look for

### Flags

| Long name   | Short | Default | Description                            |
|-------------|-------|---------|----------------------------------------|
| `--limit`   | `-l`  | 5       | Limit the number of displayed results. |
| `--verbose` | `-v`  |         | Output information about the schema.   |
