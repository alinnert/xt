# XML Schema Tools

This is a command line tool that analyzes an XML Schema and finds all element definitions and their relations with each other. This tool's primary function is to find the shortest path from all possible root elements to a given target element. 

## Example

### Input

```shell
xt ./parsx.xsd bild --limit 0
```

### Output

```
Possible paths for element "bild" (showing all 11)
- <bild>
- <alternativinhalt> => <bild>
- <bildgruppe> => <bild>
- <szene> => <bild>
- <tabellenfussnote> => <bild>
- <seite> => <bild>
- <titelei> => <seite> => <bild>
- <ausgabe> => <titelei> => <seite> => <bild>
- <ausgaben> => <ausgabe> => <titelei> => <seite> => <bild>
- <meta> => <ausgaben> => <ausgabe> => <titelei> => <seite> => <bild>
- <projekt> => <meta> => <ausgaben> => <ausgabe> => <titelei> => <seite> => <bild>
```

### Positional arguments

- **file**: The entry file of the schema
- **element-name**: The name of the element you want to look for

### Flags

| Long name   | Short | Default | Description                                                                                            |
|-------------|-------|---------|--------------------------------------------------------------------------------------------------------|
| `--limit`   | `-l`  | 5       | Limit the number of results. Only the shortest results are shown. Use "--limit 0" to show all results. |
| `--exact`   | `-e`  |         | If flag is set and search term is "elem" only "elem" is found. Otherwise "parent/elem" is also found.  |
| `--verbose` | `-v`  |         | Output additional information about the parsed XML Schema.                                             |
