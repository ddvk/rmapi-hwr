# HWR for rmapi (wip)

## Prereq 

- register for a developer account on myScript
- get the application key and hmac

## Usage
- set the env variables:  
```
    export RMAPI_HWR_APPLICATIONKEY="some application id key"
    export RMAPI_HWR_HMAC="some hmac stuff"
```

- get the zip of a notebook or pdf with `rmapi get`
- run with `go run main.go the.zip`
- or build `go build`
- `-h` for supported options



## Status
- only a single page is being converted



