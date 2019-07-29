

## Integration of LCD and RPC URLs
- `default_settings.json` : add lcd ,rpc urls and chain-id

## How to run The Commit Explorer

1. Copy `default_settings.json` to `settings.json`.
2. Update the RPC and LCD URLs.
3. Update Bech32 address prefixes.
4. Update genesis file location.

### Run in local

```
meteor npm install
meteor update
meteor --settings settings.json
```

