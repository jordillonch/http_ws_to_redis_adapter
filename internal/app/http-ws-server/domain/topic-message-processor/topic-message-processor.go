package topic_message_processor

type TopicMessageProcessor interface {
	Process(topicMessage TopicMessage) (TopicMessage, error)
}
