# NSOne Maсros

## Commands

`nom newenv $BUILD` - clean env, pull $BUILD, run containers, wait for healthy containers and make bootstrap-api

#### Flags

`--up` list of containers for docker-compose up. "core dhcp data dns" by default (TODO)
 
`--noup` skip docker-compose up

`--nobs` skip bootstrap-api


