//number of nodes in testnet
variable "nodes" {
  default = 4
}

//docker user. Needs rw permission in conf/
variable "user" {
  default = "1000"
}

//evml command (solo, dag1, raft)
variable "command" {
  default = "solo"
}

//evml Docker Image version tag
variable "version" {
  default = "latest"
}

/*
  directory containing the folders to be mounted as volumes in each container. 
  These volumes will be mounted in /.evm where evml is configured to look 
  by default. For each node, there are files related to eth (accounts, genesis 
  file, keys, etc), the consensus system (ex DAG1 peers.json, key),
  and a evml.toml file containing configuration for eth and the consensus 
  system.

  ex: conf/
    node0
    │   ├── dag1
    │   │   ├── peers.json
    │   │   └── priv_key.pem
    │   ├── evml.toml
    │   └── eth
    │       ├── genesis.json
    │       ├── keystore
    │       │   └── UTC--2018-09-24T15-46-41.072334466Z--bd3ef129b4bd4336c71153b8e10b5bc1692efa3f
    │       └── pwd.txt
    ├── node1
    │   ├── dag1
    │   │   ├── peers.json
    │   │   └── priv_key.pem
    │   ├── evml.toml
    │   └── eth
    │       ├── genesis.json
    │       ├── keystore
    │       │   └── UTC--2018-09-24T15-46-43.020722903Z--81a1ca948588423582cc2649fa0362debc5a581d
    │       └── pwd.txt
*/
variable "conf" {
  default = ""
}
