import { Injectable, Inject } from '@nestjs/common';
import axios from 'axios';

@Injectable()
export class OperatorRegistrationService {
    constructor(@Inject('API_URL')private readonly apiUrl: string = process.env.API_BASE_URL!) {
        this.apiUrl = process.env.API_BASE_URL || 'http://localhost:3000';
    }

    async registerOperator(operatorData: any): Promise<any> {
        try {
            console.log(`Registering operator: ${operatorData.name}`);
            const response = await axios.post(`${this.apiUrl}/registerOperator`, operatorData);
            console.log(`Operator registered successfully: ${operatorData.name}`);
            return response.data;
        } catch (error) {
            console.error('Error registering operator:', error.message);
            throw error;
        }
    }

    async registerEndPoint(endpointData: any): Promise<void> {
        try {
            console.log(`Registering endpointData`);
            const response = await axios.post(`${this.apiUrl}/registerTransferEndPoint`, endpointData);
            console.log("Operator registered successfully");
            return response.data;
        } catch (error) {
            console.error('Error registering operator:', error.message);
            throw error;
        }
    }
}