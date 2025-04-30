import { Module } from '@nestjs/common';
import { OperatorController } from './controllers/operator.controller';
import { OperatorFetchService } from './services/operator-fetch.service';
import { OperatorRegistrationService } from './services/operator-registration.service';
import { TokenService } from './services/token.service';

@Module({
    controllers: [OperatorController],
    providers: [
        OperatorFetchService,
        OperatorRegistrationService,
        TokenService,
        {
            provide: 'API_URL', // Custom provider token
            useValue: process.env.API_BASE_URL || 'http://localhost:3000', // Provide the value
        },
    ],
})
export class OperatorModule {}