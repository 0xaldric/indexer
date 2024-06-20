const { Kafka } = require('kafkajs')

const kafka = new Kafka({
  clientId: 'my-app',
  brokers: ['localhost:29092']
})
const consumer = kafka.consumer({ groupId: 'myGroup' })

const run = async () => {
  // Consuming
  await consumer.connect()
  await consumer.subscribe({ topic: 'jetton_transfer_notification', fromBeginning: true })

  await consumer.run({
    eachMessage: async ({ topic, partition, message }) => {
      // convert message.value into JSON
      json = JSON.parse(message.value.toString())
      try {
        // check if body is already an object
        if (typeof json.body === 'object') {
          body = json.body;
        } else {
          body = JSON.parse(json.body)
        }
        console.log(body)
      } catch(e) {
        console.log(e)
      }
      console.log(json)
      console.log({
        partition,
        offset: message.offset,
        value: message.value.toString(),
      })
    }
  })
}

run().catch(console.error)