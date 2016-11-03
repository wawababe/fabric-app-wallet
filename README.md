# wallet application
## Set up enviroment for development
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
  using the docker-compose file in directory docker
  `$ cd docker && docker-compose up`
  
  Some extra knowledge to make use of docker: 
  
  1. If you want to stop the service, just type the command 
    `$ docker-compose down`
  2. If you want to look up the logs separately, you could watch it by
    `$ docker logs -f <container id | container name>`

### Set up mysql
1. Make sure you have installed mysql (version  Ver 14.14 Distrib 5.7.14)
2. sign in mysql as root and create a new account for development
3. modify the USERNAME and PASSWORD specified in the scripts/mysql/create_table.sh 
4. create tables by execute the script:
  `$ ./create_table.sh`


## Set up the project 
  Make sure you have installed Intellij IDEA 2016.2 and the plugin Go for golang programming language
### Set up wallet application server
1. import the source code under diretory consolesrvc as a project through IntelliJ IDEA 2016.2
2. configure the environmental variables
    PEER_ADDRESS: http://127.0.0.1:7050

### Set up chaincode for development
1. import the chaincode under directory contracts as a project through IntelliJ IDEA 2016.2
2. ...
  
