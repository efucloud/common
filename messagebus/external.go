package messagebus

const (
	MessagePubDapr     = "dapr"
	MessagePubKafka    = "kafka"
	MessagePubPulsar   = "pulsar"
	MessagePubRocketMQ = "rocketmq"
	MessagePubRedis    = "redis"
	MessagePubNats     = "nats"
	MessagePubRabbitMQ = "rabbitmq"
)
const (
	TopicOrganization          = "/eauth/organization"
	TopicOrganizationAccount   = "/eauth/account"
	TopicOrganizationWorkspace = "/eauth/workspace"
)
