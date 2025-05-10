import { Module } from '@nestjs/common';
import { TransferController } from './controllers/transfer.controler';
import { TransferService } from './services/transfer.service';
import { OperatorFetchService } from './services/operator-fetch.service';
import { PubSubService } from './services/PubSubService';
import { ClientsModule, Transport } from '@nestjs/microservices';

@Module({
  imports: [
    ClientsModule.register([
      {
        name: 'DELETE_CITIZEN_CLIENT',
        transport: Transport.RMQ,
        options: {
          urls: [process.env.RABBITMQ_HOST || 'amqp://localhost:5672'],
          queue: 'delete_citizen_queue',
          queueOptions: { durable: true },
        },
      },
      {
        name: 'DELETE_DOCUMENTS_CLIENT',
        transport: Transport.RMQ,
        options: {
          urls: [process.env.RABBITMQ_HOST || 'amqp://localhost:5672'],
          queue: 'delete_documents_queue',
          queueOptions: { durable: true },
        },
      },
      {
        name: 'REGISTER_CITIZEN_CLIENT',
        transport: Transport.RMQ,
        options: {
          urls: [process.env.RABBITMQ_HOST || 'amqp://localhost:5672'],
          queue: 'register_citizen_queue',
          queueOptions: { durable: true },
        },
      },
      {
        name: 'REGISTER_DOCUMENTS_CLIENT',
        transport: Transport.RMQ,
        options: {
          urls: [process.env.RABBITMQ_HOST || 'amqp://localhost:5672'],
          queue: 'register_documents_queue',
          queueOptions: { durable: true },
        },
      },
    ]),
  ],
  controllers: [TransferController],
  providers: [
    TransferService,
    OperatorFetchService, // <-- This is what was missing
    PubSubService
  ],
})
export class TransferModule {}