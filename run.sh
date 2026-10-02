#!/bin/bash
GO_LOG=$GO_LOG GO_BIN=$GOBIN docker stack deploy --compose-file=$1/run.yml $2