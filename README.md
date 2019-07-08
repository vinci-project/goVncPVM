# goVncPVM

This module is a base network implementation of Vncsphere nodes in PVM-shard.
It provides REST endpoints for work with smart-contracts and it commands.
List of endpoints will be expanding in future.
Also goVncPVM provides simple TCP-sockets for data exchange between nodes of shard.

For more information read Docs.

## REST endpoints

### Smart-conracts endpoints
* **/wallet/transaction** Endpoint for making simple transactions   
* **/wallet/getBalance** Endpoint for requesting user balance
* **/wallet/tranStatus** Endpoint for checking transaction status  

### Blockhain endpoints
* **/blockchain/getBHeight** Endpoint for requesting current blockchain height
* **/blockchain/getTran** Endpoint for getting transaction by key
* **/blockchain/getBlock** Endpoint for getting block by height
* **/blockchain/getVersion** Endpoint for getting version of node software
* **/blockchain/getNodes** Endpoint for getting current nodes (twigs and stem) of shard

## Default ports

Default port for REST is 5000  
Default port for TCP is 3333
