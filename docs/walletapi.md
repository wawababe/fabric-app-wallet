# APIs of the Blockchain App-Wallet

## 1 Authentication Module
### 1.1 Sign up

  *METHOD*: POST

  *Route*: /auth/signup

  *Request*:
```
username:
password:
```
  *Response*:

```Json
{
  "status": "ok",
  "sessionid": "2ef88fc0-d333-4ea6-8b78-644d6bb8a4c5",
  "authtoken": "03b02c2c98712e79a76fe1629645c6f4a5b35b09"
}
```

```Json
{
  "status": "error",
  "message": "failed to signup duplicate user"
  ""
}
```

### 1.2 Login

*METHOD*: POST

*Route*: /auth/login

*Request*:
```
username:
password:
```
*Response*:

```Json
{
  "status": "ok",
  "sessionid": "b942c55c-4789-4a71-af46-c156ad8166b2",
  "authtoken": "44db8bb90026acfb716ad397e996c76e2d692985"
}
```

```Json
{
  "status": "error",
  "message": "bad request: user with name sss not exist"
}
```

### 1.3 Refresh

*METHOD*: POST

*Route*: /auth/refresh

*Request*:
```
username:
sessionid:
authtoken:
```
*Response*:

```Json
{
  "status": "ok",
  "message": "succeeded in refreshing session",
  "useruuid": "12753fa0-7b04-43ed-9278-979d6621e36b"
}
```

```Json
{
  "status": "error",
  "message": "unauthorized"
}
```

### 1.4 Logout

*METHOD*: POST

*Route*: /auth/refresh

*Request*:
```
username:
sessionid:
authtoken:
```
*Response*:

```Json
{
  "status": "ok",
  "message": "succeeded in logging out",
  "useruuid": "12753fa0-7b04-43ed-9278-979d6621e36b"
}
```

```Json
{
  "status": "error",
  "message": "unauthorized"
}
```


## 2 Wallet Module
### 2.1 Create Account

*METHOD*: POST

*Route*: /wallet/account/create

*Request*:
```
username:
sessionid:
authtoken:
accountname:
```
  *Response*:

```Json
{
  "status": "ok",
  "message": "succeeded in creating task for accountcreate event",
  "useruuid": "12753fa0-7b04-43ed-9278-979d6621e36b",
  "accountuuid": "9a707348-0aed-468c-8cf9-6274334fa754",
  "taskuuid": "b366b7c0-0fa0-4f32-b3af-168b3bfaf776"
}
```

```Json
{
  "status": "error",
  "message": "bad request: account name should not be empty"
}
```


### 2.2 List Account

*METHOD*: POST

*Route*: /wallet/account/list

*Request*:
```
username:
sessionid:
authtoken:
```
  *Response*:

```Json
{
  "status": "ok",
  "message": "succeeded in listing accounts",
  "useruuid": "12753fa0-7b04-43ed-9278-979d6621e36b",
  "accountlist": [
    {
      "rowid": 5,
      "accountuuid": "9a707348-0aed-468c-8cf9-6274334fa754",
      "useruuid": "12753fa0-7b04-43ed-9278-979d6621e36b",
      "accountname": "testaccount",
      "accountid": "94567239fbf815cc7d7db87fb94ea6c6",
      "amount": 1000,
      "bc_txuuid": "d5c2895d-c646-4b8d-94d2-11bd58b44eab",
      "status": "created"
    }
  ]
}
```

```Json
{
  "status": "error",
  "message": "unauthorized"
}
```

### 2.3 Account Transfer

*METHOD*: POST

*Route*: /wallet/transaction/transfer

*Request*:
```
username:
sessionid:
authtoken:
payeraccountid:
payeeaccountid:
amount:
```
  *Response*:

```Json
{
  "status": "ok",
  "message": "succeeded in creating task for transferring",
  "useruuid": "e0d3ad74-4b8e-4f73-9d15-8f393b3d2dd9",
  "taskuuid": "1619ea02-2890-4e90-b684-ae42f6fd8960"
}
```

```Json
{
  "status": "error",
  "message": "bad request: account residual amount is not enough to pay",
  "useruuid": "e0d3ad74-4b8e-4f73-9d15-8f393b3d2dd9"
}
```

### 2.4 Account Transfer History

*METHOD*: POST

*Route*: /wallet/transaction/list

*Request*:
```
username:
sessionid:
authtoken:
```
  *Response*:

```Json
{
  "status": "ok",
  "message": "succeeded in listing transactions",
  "useruuid": "e0d3ad74-4b8e-4f73-9d15-8f393b3d2dd9",
  "payertransactionlist": [
    {
      "rowid": 2,
      "txuuid": "7b0ff296-c895-4bf0-887d-56f1d92c4805",
      "payeruuid": "e0d3ad74-4b8e-4f73-9d15-8f393b3d2dd9",
      "payeeuuid": "12753fa0-7b04-43ed-9278-979d6621e36b",
      "payeraccountid": "2de1b1453679a6af7b7c9b52b5106922",
      "payeeaccountid": "94567239fbf815cc7d7db87fb94ea6c6",
      "amount": 21,
      "bc_txuuid": "30074477-d506-4ac4-aa72-ae389b4c3a62",
      "bc_blocknum": 0,
      "status": "transferred"
    }
  ],
  "payeetransactionlist": [
    {
      "rowid": 1,
      "txuuid": "6e3c1f98-ffb2-483e-9d16-ca33336fa723",
      "payeruuid": "358184ec-28b8-463b-8c7f-a930c3d31e9d",
      "payeeuuid": "e0d3ad74-4b8e-4f73-9d15-8f393b3d2dd9",
      "payeraccountid": "2e45cf5ccb19e8f6df6869d183a10170",
      "payeeaccountid": "2de1b1453679a6af7b7c9b52b5106922",
      "amount": 25,
      "bc_txuuid": "",
      "bc_blocknum": 0,
      "status": "pending"
    }
  ]
}
```

```Json
{
  "status": "error",
  "message": "failed to validate request, user not exist"
}
```



### 2.5 Query Transaction Record in Blockchain

*METHOD*: POST

*Route*: /blockchain/transaction

*Request*:
```
username:
sessionid:
authtoken:
bc_txuuid: 
```
  *Response*:

```Json
{
  "status": "ok",
  "useruuid": "12753fa0-7b04-43ed-9278-979d6621e36b",
  "txdetail": {
    "type": 2,
    "chaincodeID": "EgZ3YWxsZXQ=",
    "payload": "Cm8IARIIEgZ3YWxsZXQaYQoNY3JlYXRlYWNjb3VudAokMTI3NTNmYTAtN2IwNC00M2VkLTkyNzgtOTc5ZDY2MjFlMzZiCiQ5YTcwNzM0OC0wYWVkLTQ2OGMtOGNmOS02Mjc0MzM0ZmE3NTQKBDEwMDA=",
    "txid": "d5c2895d-c646-4b8d-94d2-11bd58b44eab",
    "timestamp": {
      "seconds": 1477963792,
      "nanos": 155681597
    },
    "nonce": "7IcQ/Dz2Cfcl9T4/rkCL4XrCVChiNiUy",
    "cert": "MIICUTCCAfegAwIBAgIQPD7rSagkSlypufvZrbuEGTAKBggqhkjOPQQDAzAxMQswCQYDVQQGEwJVUzEUMBIGA1UEChMLSHlwZXJsZWRnZXIxDDAKBgNVBAMTA3RjYTAeFw0xNjEwMjcwNzM3NDBaFw0xNzAxMjUwNzM3NDBaMEUxCzAJBgNVBAYTAlVTMRQwEgYDVQQKEwtIeXBlcmxlZGdlcjEgMB4GA1UEAxMXVHJhbnNhY3Rpb24gQ2VydGlmaWNhdGUwWTATBgcqhkjOPQIBBggqhkjOPQMBBwNCAAQkEDviH82lmDOi6vWgSeu/rFS0U0kKiviDP8BCbd/glHjaRwU5PtvVJPRbrehbTF/OinJyJmjztzEZenatu7FFo4HcMIHZMA4GA1UdDwEB/wQEAwIHgDAMBgNVHRMBAf8EAjAAMA0GA1UdDgQGBAQBAgMEMA8GA1UdIwQIMAaABAECAwQwTQYGKgMEBQYHAQH/BEAGgdzMT4c1BLov5W0NGDM4lJgeVG1KVjFtJYcx2sF9p6xii7iw/DWtIsjbOm04qz/VWqMWeJRP5gGjzX6/flZxMEoGBioDBAUGCARAhcZ47O2P+DqBmgNdfOQDTlLmwp/GDFdHKWFYuwLYbQ2FjslAT12HvXdV3JkM6PhVK+OZ+a49Cs/oNwp1qyXWtzAKBggqhkjOPQQDAwNIADBFAiAqHiIeX653n+DHibE4xbthiuj1YkYcpG7asFSUaTHzQQIhAK/CYz7qqLjmYS3eloEtnjxvxRF/5uyNxRVKw7+K1DSZ",
    "signature": "MEYCIQDBbqhwyQ4pq9XE8FOgB2bMlu4PTIjpEaeQXgbaCCxPMQIhAOhvddrC21DMzxrPbGtR8xIwAIktRmlaVIBE+rIuwtRM"
  }
}
```

```Json
{
  "status": "error",
  "message": "not found: transaction d5c2895d-c646-4b8d-94d2-11bd58b44ea not exist in the blockchain",
  "useruuid": "12753fa0-7b04-43ed-9278-979d6621e36b"
}
```
