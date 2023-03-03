## Performance of OPA

Run the program:
```bash
go run main.go
```

Output:

```bash
Loading API Keys into memory...
API Keys loaded into memory.

Allowed count:  50
Rejected count:  50
2023/03/03 15:54:36 Avg evaluation took: 9.803µs
```

## Reproduce the test

Generate `data.json` that OPA uses to evaluate API Keys.

```bash
go run generate/keys.go
```

This stores the data.json in `./bundle/apikeys/data.json`

Proceed to recreate the optimized OPA bundle.

```bash
bundle/bundle.sh
```

And then you can run the program as shown at the top of the article.