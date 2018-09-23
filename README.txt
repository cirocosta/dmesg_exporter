dmesg_exporter - metrics from the kernel messages ring buffer

Usage:
  dmesg_exporter [OPTIONS] start [start-OPTIONS]

Help Options:
  -h, --help         Show this help message

[start command options]
          --path=    path to serve metrics (default: /metrics)
          --address= address to listen for prometheus scraping (default: :9000)
      -t, --tail     whether the reader should seek to the end

Sample metrics:
	dmesg_logs{facility="daemon",priority="info"} 10
	dmesg_logs{facility="kern",priority="debug"} 69
	dmesg_logs{facility="kern",priority="err"} 4
	dmesg_logs{facility="kern",priority="info"} 380
	dmesg_logs{facility="kern",priority="notice"} 47
	dmesg_logs{facility="kern",priority="warning"} 38
	dmesg_logs{facility="syslog",priority="info"} 1
	dmesg_logs{facility="user",priority="warning"} 4

Install:
	go get github.com/cirocosta/dmesg_exporter

