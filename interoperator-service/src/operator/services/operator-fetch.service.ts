import { Injectable } from '@nestjs/common';
import axios from 'axios';

@Injectable()
export class OperatorFetchService {
    private readonly apiUrl: string;

    constructor() {
        this.apiUrl = process.env.API_BASE_URL || 'http://localhost:3000';
    }

    async getOperators(): Promise<any[]> {
        try {
            console.log(`Fetching operators from URL: ${this.apiUrl}/getOperators`);
            const response = await axios.get(`${this.apiUrl}/getOperators`);
            console.log('Operators fetched successfully.');
            return response.data;
        } catch (error) {
            console.error('Error fetching operators:', error.message);
            throw error;
        }
    }

    async getOperatorByName(name: string): Promise<any> {
        const operators = await this.getOperators();
        return operators.find((op) => op.operatorName === name);
    }

    async getSelfOperator(): Promise<any> {
        const operators = process.env.OPERATOR_ID;
        if (!operators) {
            const operators = await this.getOperators();
            const operator=operators.find((op) => op.operatorName ===  process.env.OPERATOR_NAME);
            return operator;
        }
        else {return this.getOperatorById(operators);}
    }

    async getOperatorById(id: string): Promise<any> {
        const operators = await this.getOperators();
        return operators.find((op) => op._id === id);
    }
}