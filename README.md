# numerology

HTTP API that finds arithmetic expressions equaling a target number from input digits.

## Usage

```bash
go build && ./numerology
```

```bash
curl localhost:8080/23/123456
# 1*2*3*4+5-6

curl localhost:8080/23/2023
# 0*2+23

curl localhost:8080/23
# Uses today's date (DDMMYYYY) as input

curl localhost:8080/42/123456
# Find expression equaling 42

curl localhost:8080/23/123456?format=json
# {"input":"123456","target":23,"expression":"1*2*3*4+5-6","result":23}

curl localhost:8080/23/123456?format=text
# Using 123456 to reach 23: 1*2*3*4+5-6 = 23

curl localhost:8080/23?format=text
# Today is 18012026 and 1+8+0+1+2+0+2+6+3 = 23
```

## Configuration

| Variable | Default | Description |
|----------|---------|-------------|
| `HOST` | `0.0.0.0` | Bind address |
| `PORT` | `8080` | Listen port |

## API

### `GET /<target>/<input>?format=<format>`

Returns an expression that evaluates to the target using the provided digits.

- `target`: The number the expression should equal (required)
- `input`: Digits to use (optional, defaults to today's date as DDMMYYYY)
- `format`: Output format (optional, defaults to `raw`)
  - `raw`: Expression only (e.g., `1+2*3+4+5+6`)
  - `json`: JSON object with input, target, expression, result
  - `text`: Human-readable sentence
- Non-digit characters in input are filtered out
- Uses standard operator precedence (`*` `/` before `+` `-`)
- Integer division only (non-divisible operations skipped)
- Returns first valid expression found
- Returns 404 if no solution exists
- Returns 400 if target is not a valid integer

## Algorithm

1. Generate all permutations of input digits
2. For each permutation, generate all groupings (partitions into multi-digit numbers)
3. For each grouping, try all operator combinations (`+` `-` `*` `/`)
4. Evaluate with standard precedence, return first expression equaling the target
