all: imei-imsi-ccid simple send receive set-apn

imei-imsi-ccid:
	@cd examples/$@  && go build -o ../../bin/$@

simple:
	@cd examples/$@  && go build -o ../../bin/$@

send:
	@cd examples/$@  && go build -o ../../bin/$@

receive:
	@cd examples/$@  && go build -o ../../bin/$@

set-apn:
	@cd examples/$@  && go build -o ../../bin/$@

