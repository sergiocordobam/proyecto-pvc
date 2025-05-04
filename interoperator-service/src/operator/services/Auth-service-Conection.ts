import { Injectable, HttpException, HttpStatus } from '@nestjs/common';
import axios from 'axios';

@Injectable()
export class CitizenService {
    private readonly authServiceUrl = process.env.AUTH_SERVICE_URL || 'http://localhost:4000';

    async fetchCitizenInfo(citizenId: string): Promise<any> {
        try {
            const response = await axios.get(`${this.authServiceUrl}/citizens/${citizenId}`);
            return response.data;
        } catch (error) {
            throw new HttpException(
                `Error fetching citizen info: ${error.message}`,
                HttpStatus.INTERNAL_SERVER_ERROR,
            );
        }
    }

    async deleteCitizen(citizenId: string): Promise<void> {
        try {
            await axios.delete(`${this.authServiceUrl}/citizens/${citizenId}`);
        } catch (error) {
            throw new HttpException(
                `Error deleting citizen: ${error.message}`,
                HttpStatus.INTERNAL_SERVER_ERROR,
            );
        }
    }
}