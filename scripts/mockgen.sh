#!/usr/bin/env sh

mockgen_cmd="mockgen"
$mockgen_cmd -source=client/client.go -package mock -destination client/mock/client.go
