import { Injectable } from '@nestjs/common';
import { PubSub } from '@google-cloud/pubsub';

@Injectable()
export class PubSubService {
    private readonly pubSubClient: PubSub;

    constructor() {
        this.pubSubClient = new PubSub({
            projectId: process.env.GCP_PROJECT_ID || 'zeta-matrix', // Explicitly set the project ID
        });
    }

    async publishMessage(topicName: string, message: object): Promise<void> {
        try {
            // Use the publishMessage method with the json property
            await this.pubSubClient.topic(topicName).publishMessage({
                json: message, // Pass the message data as JSON
            });
            console.log(`Message published to topic ${topicName}`);
        } catch (error) {
            console.error(`Error publishing message to Pub/Sub: ${error.message}`);
            throw new Error(`Error publishing message to Pub/Sub: ${error.message}`);
        }
    }
}