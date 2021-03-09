# f5-api

Provides RESTful API access to the F5 BigIP/LTM

## Endpoints

```
GET /v1/f5/ping
GET /v1/f5/version
GET /v1/f5/metrics

GET /v1/f5/{host}/clientssl
GET /v1/f5/{host}/clientssl/{clientsslprofilename}
PUT /v1/f5/{host}/createclientssl/{clientclientsslprofilename}
PUT /v1/f5/{host}/updateclientssl/{updateclientsslprofilename}

```

## Usage

Configuration file with a list of LTM hosts

CRUD operations on:
  - SSL Client Certificates
  - SSL keys
  - SSL Client Profiles

Enable operations on one or more LTM hosts

### List Client SSL Profiles
GET

/v1/f5/flt-ltm-cluster.example.org/clientssl


### Show Client SSL Profile
GET

/v1/f5/flt-ltm-cluster.example.org/clientssl/{test.example.org}

### Create Client SSL Profile

PUT

curl -H -X PUT --data "@tmp/sslclientprofile" 'X-Auth-Token:{uuid}' "http://127.0.0.1:8080/v1/f5/flt-ltm-cluster.example.org/createclientssl/{test.example.org}" |jq

where tmp/sslclientprofile contains:
```{
"clientssl-profile": "test.example.org",
"chain": "intermediate-chain.crt",
"defaultsfrom": "clientssl",
"ciphergroup": "default-tlsv1.2",
"ciphers": "none",
"cert": "base64-encoded-certificate-pem",
"key": "base64-encoded-key-pem"
}```

### Update Client SSL Profile

PUT

curl -H -X PUT --data "@tmp/sslclientprofile" 'X-Auth-Token:{uuid}' "http://127.0.0.1:8080/v1/f5/flt-ltm-cluster.example.org/updateclientssl/{test.example.org}" |jq

where tmp/sslclientprofile contains:
```{
"clientssl-profile": "test.example.org",
"chain": "intermediate-chain.crt",
"defaultsfrom": "clientssl",
"ciphergroup": "default-tlsv1.2",
"ciphers": "none",
"cert": "base64-encoded-certificate-pem",
"key": "base64-encoded-key-pem"
}```

### Responses

```json
json here
```

## Authentication

Authentication is accomplished via a pre-shared key.  This is done via the `X-Auth-Token` header.

## Author

Darryl Wisneski <darryl.wisneski@yale.edu>

## License

GNU Affero General Public License v3.0 (GNU AGPLv3)  
Copyright (c) 2021 Yale University
