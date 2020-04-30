package config

type Config struct {
	RedisCfg *RedisCfg
	ZkCfg    *ZkCfg
	LogCfg   *LogCfg
	KafkaCfg *KafkaCfg
	MysqlCfg *MysqlCfg
}

type RedisCfg struct {
	IP   string
	Port int
	//单机
	Addr     string
	Password string
	//集群
	Address []string
}

type ZkCfg struct {
	Servers []string
}

type LogCfg struct {
	ProgramName string
	Debug       bool
}

type KafkaCfg struct {
	IP         string
	Port       int
	Address    []string
	GroupName  string
	ClientName string
}

type MysqlCfg struct {
	Address  string
	IP       string
	Port     int
	User     string
	Password string
	DBName   string
	DSN      string
}
