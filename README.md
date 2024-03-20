# go-schematized-config

check your os.Environ / .env with a json(net) schema

# usage

```
$ cat .env
MY_INTEGER_VALUE=2
$ export MY_INTEGER_VALUE=23.45
$ schematized-config ./testdata/test-schema-1.schema.jsonnet
2024/03/21 01:12:00 jsonschema: '' does not validate with .../schema#/required: missing properties: 'string_value_with_enum'
$ export string_value_with_enum=only
$ schematized-config testdata/test-schema-1.schema.jsonnet
{
  "A_NUMERIC_VALUE": 23.45,
  "MY_INTEGER_VALUE": 2,
  "_____A_STRING_VALUE____with_default__": "underscores_and spaces",
  "string_value_with_enum": "only"
}
```
