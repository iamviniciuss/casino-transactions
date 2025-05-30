const { Kafka } = require('kafkajs');

const kafka = new Kafka({
  clientId: 'demo-producer',
  brokers: ['kafka:29092'],
});

const producer = kafka.producer();

const userIds = [
  "573a37e7-832a-4ecd-9691-41ff29afb955",
  "912457bc-7eed-4170-aba5-8a13c35a9d8a",
  "62d54d96-88f4-4111-8564-c043d710bdcd",
];

const getRandomUserId = () => {
  const index = Math.floor(Math.random() * userIds.length);
  return userIds[index];
};

const run = async () => {
  await producer.connect();

  setInterval(async () => {
    const message = {
      user_id: getRandomUserId(),
      transaction_type: Math.random() > 0.5 ? 'bet' : 'win',
      amount: Number((Math.random() * 1000).toFixed(2)),
      timestamp: new Date().toISOString(),
    };

    await producer.send({
      topic: 'casino-transactions',
      messages: [
        {
          key: 'transaction',
          value: JSON.stringify(message),
        },
      ],
    });

    console.log('Message sent:', message);
  }, 500);
};

run().catch(console.error);