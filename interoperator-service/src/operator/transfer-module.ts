import { Module } from '@nestjs/common';
import { TransferController } from './controllers/transfer.controler';
import { TransferService } from './services/transfer.service';
import { ClientsModule, Transport } from '@nestjs/microservices';

@Module({
    imports: [
        ClientsModule.register([
            {
                name: 'TRANSFER_CONFIRMATIONS',
                transport: Transport.RMQ,
                options: {
                    urls: [process.env.RABBITMQ_URL || 'amqp://localhost:5672'],
                    queue: 'transfer_confirmations',
                    queueOptions: { durable: true },
                },
            },
        ]),
    ],
    controllers: [TransferController],
    providers: [TransferService],
})
export class TransferModule {}