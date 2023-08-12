objects=$(PWD)/internal/protos/*.proto
SRC_DIR=$(PWD)/internal/protos
DST_DIR=$(PWD)/pkg/protos

.PHONY:
	echo "hello"

grpc: $(objects)
	protoc -I=$(SRC_DIR) \
	    --go_out=$(DST_DIR) \
	    --go_opt=paths=source_relative \
	    --go-grpc_out=$(DST_DIR) \
	    --go-grpc_opt=paths=source_relative \
	    $(objects)

