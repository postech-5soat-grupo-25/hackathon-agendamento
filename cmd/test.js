const amqp = require('amqplib');

async function sendMessage() {
    try {
        // Connect to RabbitMQ server
        const connection = await amqp.connect('amqp://guest:guest@localhost:5672');
        const channel = await connection.createChannel();

        // Declare a request queue
        const requestQueue = 'appointments';
        await channel.assertQueue(requestQueue, {
            durable: true
        });

        // Declare a response queue
        const responseQueue = 'responses'; // You can also create a unique queue per request
        await channel.assertQueue(responseQueue, {
            durable: true
        });

        // Message to be sent
        const message = {
            type: "getdoctorworkinghours",
            body: {
                DoctorID: 123
            }
        };

        // Send message to the request queue
        const correlationId = generateUuid(); // Generate a unique correlation ID
        channel.sendToQueue(requestQueue, Buffer.from(JSON.stringify(message)), {
            correlationId: correlationId,
            replyTo: responseQueue // Set the replyTo field to indicate where to send the response
        });

        console.log(" [x] Sent '%s'", JSON.stringify(message));

        // Set up a consumer for the response queue
        channel.consume(responseQueue, (msg) => {
            if (msg.properties.correlationId === correlationId) {
                console.log(" [.] Received response: '%s'", msg.content.toString());
                channel.ack(msg); // Acknowledge the message
                // Close the connection after receiving the response
                setTimeout(() => {
                    channel.close();
                    connection.close();
                }, 500);
            }
        }, {
            noAck: false // Ensure messages are acknowledged after processing
        });
    } catch (error) {
        console.error("Error sending message to RabbitMQ:", error);
    }
}

// Helper function to generate a unique ID
function generateUuid() {
    return Math.random().toString() + Math.random().toString() + Math.random().toString();
}

// Call the function to send a message
sendMessage();
