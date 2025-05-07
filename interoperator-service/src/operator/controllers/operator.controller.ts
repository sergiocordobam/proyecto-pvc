import { Controller, Get, Req, Res } from '@nestjs/common';
import { Request, Response } from 'express';
import { OperatorFetchService } from '../services/operator-fetch.service';
import { OperatorRegistrationService } from '../services/operator-registration.service';
import { TokenService } from '../services/token.service';
import { RegisterEndpointDto } from '../DTO/RegisterEndpointDto';

@Controller('operators')
export class OperatorController {
    constructor(
        private readonly fetchService: OperatorFetchService,
        private readonly registrationService: OperatorRegistrationService,
        private readonly tokenService: TokenService,
    ) {
        this.initialize();
    }

    @Get()
    async fetchOperators(@Req() req: Request, @Res() res: Response): Promise<void> {
        try {
            const operators = await this.fetchService.getOperators();
            res.status(200).json(operators);
        } catch (error) {
            console.error('Error fetching operators:', error.message);
            res.status(500).json({ error: 'Failed to fetch operators' });
        }
    }
    @Get('self')
    async fetchSelfOperator(@Req() req: Request, @Res() res: Response): Promise<void> {
        try {
            const operator = await this.fetchService.getSelfOperator();
            res.status(200).json(operator);
        } catch (error) {
            console.error('Error fetching self operator:', error.message);
            res.status(500).json({ error: 'Failed to fetch self operator' });
        }
    }

    async initialize(): Promise<void> {
        try {
            const operatorName = process.env.OPERATOR_NAME!;
            const operator = await this.fetchService.getOperatorByName(operatorName);
            if (operator) {
                console.log(`Operator ${operatorName} already exists.`);
                const registeredOperator = await this.fetchService.getSelfOperator();
                console.log(`token: ${registeredOperator._id}`);
                this.tokenService.saveToken(registeredOperator._id);
                console.log(`token saved: ${process.env.OPERATOR_ID}`)
            } else {
                console.log(`Operator ${operatorName} does not exist. Registering...`);
                const operatorData = {
                    name: process.env.OPERATOR_NAME!,
                    address: process.env.OPERATOR_ADDRESS!,
                    contactMail: process.env.OPERATOR_CONTACT_EMAIL!,
                    participants: [
                        process.env.OPERATOR_PARTICIPANTS1!,
                        process.env.OPERATOR_PARTICIPANTS2!,
                        process.env.OPERATOR_PARTICIPANTS3!,
                    ],
                };
                const registeredOperator = await this.registrationService.registerOperator(operatorData);
                this.tokenService.saveToken(registeredOperator);
                const endpointData = {
                    idOperator: registeredOperator,
                    endPoint: process.env.OPERATOR_TRANSFER_ENDPOINT,
                    endPointConfirm: process.env.OPERATOR_TRANSFER_ENDPOINT_CONFIRM,
                } as RegisterEndpointDto;
                this.registrationService.registerEndPoint(endpointData);
                console.log(`token: ${registeredOperator}`);
            }
            console.log(`Operator ${operatorName} registered successfully.`);
        } catch (error) {
            console.error('Error checking operator in system:', error.message);
        }
    }
}
