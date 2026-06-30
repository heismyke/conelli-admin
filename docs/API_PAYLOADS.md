# API Payloads

## Health

`GET /health`

```json
{
  "status": "healthy",
  "version": "1.0.0",
  "service": "Conelli Admin API"
}
```

`GET /health/ping`

```json
{
  "message": "pong"
}
```

## Admin

`GET /admin`

```json
{
  "service": "Conelli Admin API",
  "status": "ready"
}
```
