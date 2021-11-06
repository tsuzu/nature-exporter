.PHONY: build copy
build:
	GOOS=linux go build -o nature-exporter .

copy: build
	rsync -avh --rsync-path='sudo rsync' nature-exporter 10.20.40.14:/usr/local/bin/nature-exporter
	ssh pi@10.20.40.30 sudo bash -c "systemctl daemon-reload && systemctl restart nature-exporter"
