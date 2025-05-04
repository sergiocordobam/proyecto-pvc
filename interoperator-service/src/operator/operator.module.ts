import { Module } from '@nestjs/common';
import { OperatorController } from './controllers/operator.controller';
import { OperatorFetchService } from './services/operator-fetch.service';
import { OperatorRegistrationService } from './services/operator-registration.service';
import { TokenService } from './services/token.service';
import { TransferService } from './services/transfer.service';
import { CitizenService } from './services/Auth-service-Conection';
import { DocumentService } from './services/Documents-service-conection';

@Module({
    controllers: [OperatorController],
    providers: [
        OperatorFetchService,
        OperatorRegistrationService,
        TokenService,
        TransferService,
        CitizenService,
        DocumentService,
        {
            provide: 'API_URL',
            useValue: process.env.API_BASE_URL || 'http://localhost:3000',
        },
    ],
})
export class OperatorModule {}