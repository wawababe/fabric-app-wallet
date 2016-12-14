# wallet application



## Prepare for development
### Deploy hyperledger server: peer and membersrvc 
1. Make sure you have installed docker(version newer than 1.12.0).
2. Pulled the latest image from docker hub:
  `$ docker pull hyperledger/fabric-peer`

  `$ docker pull hyperledger/fabric-membersrvc`

  PS: the latest tag for fabric-baseimage has not been issued officially,
  therefore, it would be ok for you to download other substitute:
  `$ docker pull hyperledger/fabric-baseimage:x86_64-0.2.0`
  , then re-tag it as latest:
  `$ docker tag hyperledger/fabric-baseimage:x86_64-0.2.0 hyperledger/fabric-baseimage:latest`

3. Start the hyperledger/fabric-peer and hyperledger/fabric-membersrvc 
  using the [docker-compose file](docker/docker-compose.yml) in directory docker
  `$ cd docker && docker-compose up`
  
  PS: 
  
  1. If you want to stop the service, just type the command 
    `$ docker-compose down`
  2. If you want to look up the logs separately, you could watch it by
    `$ docker logs -f <container id | container name>`

### Create mysql Tables
1. Make sure you have installed mysql (version  Ver 14.14 Distrib 5.7.14) or have access to a mysql server
2. sign in mysql as root and create a new account for development
3. modify the USERNAME and PASSWORD specified in the [create_table script](scripts/mysql/create_table.sh)
4. create tables by execute the script:
  `$ ./create_table.sh`


## Set up the project 
  Make sure you have installed Intellij IDEA 2016.2 and the plugin Go for golang programming language

  ```bash
  mkdir $GOPATH/src/baas && cd $GOPATH/src/baas
  git clone URL/to/app-wallet
  ```

### Set up wallet application server
1. import the source code under diretory [consolesrvc](consolesrvc) and [vendor](vendor) as a project through IDEA
2. modify the [configuration file](consolesrvc/wallet.yaml)
    - `database.mysql.dsn`
    - `fabric.peer.address`
3. run the wallet project

### Set up chaincode for development
1. import the chaincode under directory [contracts](contracts) and [vendor](vendor) as a project through IDEA
2. set the env variables:
    - `CORE_CHAINCODE_ID_NAME=wallet;`
    - `CORE_LOGGING_CHAINCODE=DEBUG;`
    - `CORE_PEER_ADDRESS=127.0.0.1:7051;`
    - `SHIM_LOGGING_LEVEL=INFO`
3. run the chaincode
4. register the chaincode to peer by sending the following request
```json
{
  "jsonrpc": "2.0",
  "method": "deploy",
  "params": {
      "type": 1,
      "chaincodeID":{
          "name":"wallet"
        },
      "ctorMsg": {
          "function":"init",
          "args":[]
      }
  },
  "id": 1
}
```


## The RESTful APIs of the app-wallet
please refer to [restful_api](docs/walletapi.md) which is in the directory `docs`

