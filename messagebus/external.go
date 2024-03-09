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
	TopicOrganization          = "/messagebus/eauth/organization"
	TopicOrganizationAccount   = "/messagebus/eauth/account"
	TopicOrganizationWorkspace = "/messagebus/eauth/workspace"
)
