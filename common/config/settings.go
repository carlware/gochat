package config

// Configuration settings
type Configuration struct {
	Environment string `yaml:"environment" default:"local"`

	Debug struct {
		Enable bool `yaml:"enable" default:"false" comment:"allow debug mode"`
	} `yaml:"debug"`

	RabbiMQ struct {
		Host            string `yaml:"host" default:"amqp://guest:guest@localhost:5672/"`
		CommandReqQueue string `yaml:"command_req_queue" default:"commandReq"`
		CommandResQueue string `yaml:"command_res_queue" default:"commandRes"`
	} `yaml:"psql"`
}
