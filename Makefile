CURRENT_DIR=$(shell pwd)
DBURL=postgres://postgres:root@localhost:5432/?sslmode=disable

proto-gen:
	./scripts/gen-proto.sh ${CURRENT_DIR}

exp:
	export DBURL='postgres://postgres:root@localhost:5432/learning_db?sslmode=disable'

mig-up:
	migrate -path migrations -database '$(DBURL)' -verbose up

mig-down:
	migrate -path migrations -database '$(DBURL)' -verbose down

mig-create:
	migrate create -ext sql -dir migrations -seq create_table

mig-insert:
	migrate create -ext sql -dir db/migrations -seq insert_table
gen-protoAll:
	protoc --go_out=./ \
	--go-grpc_out=./ \
	submodule/Competition-Service/*.proto
