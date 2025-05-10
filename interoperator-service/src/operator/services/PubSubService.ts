import { Injectable } from '@nestjs/common';
import { PubSub } from '@google-cloud/pubsub';

@Injectable()
export class PubSubService {
    private readonly pubSubClient: PubSub;

    constructor() {
        this.pubSubClient = new PubSub();
    }

    async publishMessage(topicName: string, message: object): Promise<void> {
        try {
            const messageBuffer = Buffer.from(JSON.stringify(message));
            await this.pubSubClient.topic(topicName).publish(messageBuffer);
            console.log(`Message published to topic ${topicName}`);
        } catch (error) {
            console.error(`Error publishing message to Pub/Sub: ${error.message}`);
            throw new Error(`Error publishing message to Pub/Sub: ${error.message}`);
        }
    }
}