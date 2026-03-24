# Maskfile

This is a [mask](https://github.com/jacobdeichert/mask) task runner file.

## encode

> Encode a JSON record to Avro binary. Usage: mask encode --schema item --json '{"id":"1","name":"Widget"}'

**OPTIONS**
- schema
  - flags: -s --schema
  - type: string
  - desc: Schema name (matches a file in models/)
- json
  - flags: -j --json
  - type: string
  - desc: JSON record to encode

```bash
[ -n "$schema" ] || { echo "No schema provided — use --schema <name>"; exit 1; }
[ -n "$json" ] || { echo "No JSON provided — use --json '{...}'"; exit 1; }
[ -f "models/$schema.avsc" ] || { echo "Schema not found: models/$schema.avsc"; exit 1; }
echo "$json" | avro-tools fromjson --schema-file "models/$schema.avsc" - > "$schema.avro"
echo "Written to $schema.avro"
```

## read

> Read an Avro file with a language reader. Usage: mask read --schema item --lang go

**OPTIONS**
- schema
  - flags: -s --schema
  - type: string
  - desc: Schema name (matches a <schema>.avro file)
- lang
  - flags: -l --lang
  - type: string
  - desc: Language reader to use (go, python)

```bash
[ -n "$schema" ] || { echo "No schema provided — use --schema <name>"; exit 1; }
[ -n "$lang" ] || { echo "No lang provided — use --lang <go|python>"; exit 1; }
[ -f "$schema.avro" ] || { echo "$schema.avro not found — run 'mask encode' first"; exit 1; }
[ -d "readers/$lang" ] || { echo "No reader found for lang: $lang (add one in readers/$lang/)"; exit 1; }

case "$lang" in
  go)
    cd "readers/go" && go mod tidy && go run . "$OLDPWD/$schema.avro"
    ;;
  python)
    uv run "readers/python/read.py" "$schema.avro"
    ;;
  *)
    echo "Unknown lang: $lang"
    exit 1
    ;;
esac
```
