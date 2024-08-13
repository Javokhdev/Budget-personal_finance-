#!/bin/bash
CURRENT_DIR=$1
rm -rf ${CURRENT_DIR}/protos
for x in $(find ${CURRENT_DIR}/submodule-personal_finance* -type d); do
  protoc -I=${x} -I=${CURRENT_DIR} -I /usr/local/go --go_out=${CURRENT_DIR} \
   --go-grpc_out=${CURRENT_DIR} ${x}/*.proto
done