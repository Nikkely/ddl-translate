# ddl-translate

translate DDL with specified range

## JSON

run
```bash
./main json -path=/path/to/json "foo" "hoge"
{"foo":["ハローワールド","ペンが大きい"],"fuga":"Life is what you make it.","hoge":"人は皆、自分の運を決めるのは自分自身です。"}
```

json
```json
{
  "foo": ["Hello world", "My pen is big"]
  "hoge": "Every man is the architect of his own fortune."
  "fuga": "Life is what you make it."
}
```

edit config.ini
