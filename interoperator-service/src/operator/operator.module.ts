import { Module } from '@nestjs/common';
import { OperatorController } from './controllers/operator.controller';
import { OperatorFetchService } from './services/operator-fetch.service';
import { OperatorRegistrationService } from './services/operator-registration.service';
import { TokenService } from './services/token.service';
import { TransferService } from './services/transfer.service';

@Module({
    controllers: [OperatorController],
    providers: [
        OperatorFetchService,
        OperatorRegistrationService,
        TokenService,
        TransferService,
        {
            provide: 'API_URL',
            useValue: process.env.API_BASE_URL || 'http://localhost:3000',
        },
    ],
})
export class OperatorModule {}