package config

type Settings struct {
	Port      int       `json:"port"`
	Databases Databases `json:"databases"`
	Cluster   Cluster   `json:"cluster"`
}

type Databases struct {
	Postgres string `json:"postgres"`
	Kafka    Kafka  `json:"kafka"`
}

type Kafka struct {
	Brokers []string `json:"brokers"`
	Topics  Topics   `json:"topics"`
}

type Topics struct {
	Inspections string `json:"inspections"`
}

type Cluster struct {
	TaskService string `json:"taskService"`
}
