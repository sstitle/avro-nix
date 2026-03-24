# Maskfile

This is a [mask](https://github.com/jacobdeichert/mask) task runner file.

## hello

> This is an example command you can run with `mask hello`

```bash
echo "Hello World!"
```

## item-encode

> Encode a JSON record into an Avro binary file. Usage: mask item-encode '{"id":"1","name":"Widget"}'

**OPTIONS**
- json
  - flags: -j --json
  - type: string
  - desc: JSON record to encode

```bash
[ -n "$json" ] || { echo "No JSON provided — usage: mask item-encode --json '{\"id\":\"1\",\"name\":\"Widget\"}'"; exit 1; }
echo "$json" | avro-tools fromjson --schema-file models/item.avsc - > item.avro
echo "Written to item.avro"
```

## item-decode

> Decode item.avro back to JSON

```bash
[ -f item.avro ] || { echo "item.avro not found — run 'mask item-encode' first"; exit 1; }
avro-tools tojson item.avro
```

## item-roundtrip-go

> Encode a JSON record to Avro and read it back via Go

**OPTIONS**
- json
  - flags: -j --json
  - type: string
  - desc: JSON record to encode

```bash
[ -n "$json" ] || { echo "No JSON provided — usage: mask item-roundtrip-go --json '{\"id\":\"1\",\"name\":\"Widget\"}'"; exit 1; }
echo "$json" | avro-tools fromjson --schema-file models/item.avsc - > item.avro
cd scripts/go && go mod tidy && go run . "$OLDPWD/item.avro"
```

## item-roundtrip

> Encode a JSON record to Avro and read it back via Python

**OPTIONS**
- json
  - flags: -j --json
  - type: string
  - desc: JSON record to encode

```bash
[ -n "$json" ] || { echo "No JSON provided — usage: mask item-roundtrip --json '{\"id\":\"1\",\"name\":\"Widget\"}'"; exit 1; }
echo "$json" | avro-tools fromjson --schema-file models/item.avsc - > item.avro
uv run scripts/read_item.py item.avro
```

## item-schema

> Print the embedded schema from item.avro

```bash
[ -f item.avro ] || { echo "item.avro not found — run 'mask item-encode' first"; exit 1; }
avro-tools getschema item.avro
```
