all:
	@cd examples/imei-imsi-ccid && go build -o ../../bin/imei-imsi-ccid
	@cd examples/simple && go build -o ../../bin/simple
	@cd examples/send && go build -o ../../bin/send
	@cd examples/receive && go build -o ../../bin/receive
