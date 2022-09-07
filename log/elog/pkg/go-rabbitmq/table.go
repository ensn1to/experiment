package gorabbitmq

import amqp "github.com/rabbitmq/amqp091-go"

// Wrapper of the amqp table
type Table map[string]interface{}

func tableToAMQPTable(table Table) amqp.Table {
	new := amqp.Table{}
	for k, v := range table {
		new[k] = v
	}
	return new
}
