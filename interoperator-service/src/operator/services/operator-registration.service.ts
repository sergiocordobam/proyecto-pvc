import { Injectable, Inject } from '@nestjs/common';
import axios from 'axios';
import { RegisterOperatorDto } from '../DTO/RegisterOperatorDto';
import { RegisterEndpointDto } from '../DTO/RegisterEndpointDto';

@Injectable()
export class OperatorRegistrationService {
    constructor(@Inject('API_URL') private readonly apiUrl: string = process.env.API_BASE_URL!,) {}

    async registerOperator(dto: RegisterOperatorDto): Promise<any> {
        try {
            console.log(`${this.apiUrl}/registerOperator`);
            const response = await axios.post(`${this.apiUrl}/registerOperator`, dto);
            return response.data;
        } catch (error) {
            throw new Error(`Error registering operator: ${error.message}`);
        }
    }

    async registerEndPoint(dto: RegisterEndpointDto): Promise<void> {
        try {
            console.log(`Registering endpointData`);
            const response = await axios.post(`${this.apiUrl}/registerTransferEndPoint`, dto);
            console.log("Operator registered successfully");
            return response.data;
        } catch (error) {
            console.error('Error registering operator:', error.message);
            throw error;
        }
    }
}