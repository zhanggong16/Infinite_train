package metric

// 存储名称和底层的映射
var linuxMetricsMap = map[string]string {
	"docker.network.incoming": "docker.network.incoming",
	"docker.network.outgoing": "docker.network.outgoing",
	"docker.cpu.load": "docker.cpu.load",
	"docker.memory.pused": "docker.memory.pused",
	"docker.network.total": "docker.network.total",
	"docker.disk0.used": "docker.disk0.used",
	"docker.disk1.used": "docker.disk1.used",
	"docker.disk1.total.kb": "docker.disk1.total.kb",
	"docker.disk1.used.kb": "docker.disk1.used.kb",
	"docker.cpu.util": "docker.cpu.util",
	"docker.disk.iops.read": "docker.disk.iops.read",
	"docker.disk.iops.write": "docker.disk.iops.write",
}
