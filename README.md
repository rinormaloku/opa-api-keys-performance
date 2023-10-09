## Performance of OPA

### Setup the data

Generate `data.json` that OPA uses to evaluate API Keys.

```bash
go run generate/keys.go
```

This stores the file in `./bundle/apikeys/data.json`

Proceed to recreate the optimized OPA bundle.

```bash
bundle/bundle.sh
```

Next, we will validate OPA in two ways:
1. Consume the Bundle locally
2. Poll a WebServer that hosts the Bundle

### Consume the Bundle locally

Run the program:
```bash
go run localdata/local.go
```

Output:

```bash
Loading API Keys into memory...
API Keys loaded into memory.

Allowed count:  50
Rejected count:  50
2023/03/06 14:27:17 Avg evaluation took: 12ms
```

### Poll a WebServer that hosts the Bundle

In a second terminal switch to the `bundle` directory and run the following command:

```bash
cd bundle
python3 -m http.server 8981
```

In the first terminal, from the root of the repo run the following command:

```bash
go run polldata/poll.go
```