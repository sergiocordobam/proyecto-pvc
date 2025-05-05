/*import { Injectable, HttpException, HttpStatus } from '@nestjs/common';
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
}*/
import { Injectable } from '@nestjs/common';
import { EventPattern, Payload } from '@nestjs/microservices';
import axios from 'axios';

@Injectable()
export class CitizenService {
    private readonly authServiceUrl = process.env.AUTH_SERVICE_URL || 'http://localhost:4000';

    @EventPattern('fetch_citizen_info')
    async fetchCitizenInfo(@Payload() data: { citizenId: string }): Promise<any> {
        try {
            console.log(`Fetching citizen info for ID: ${data.citizenId}`);
            const response = await axios.get(`${this.authServiceUrl}/citizens/${data.citizenId}`);
            return response.data;
        } catch (error) {
            console.error(`Error fetching citizen info: ${error.message}`);
            throw new Error(`Error fetching citizen info: ${error.message}`);
        }
    }

    @EventPattern('delete_citizen')
    async deleteCitizen(@Payload() data: { citizenId: string }): Promise<void> {
        try {
            console.log(`Deleting citizen with ID: ${data.citizenId}`);
            await axios.delete(`${this.authServiceUrl}/citizens/${data.citizenId}`);
        } catch (error) {
            console.error(`Error deleting citizen: ${error.message}`);
            throw new Error(`Error deleting citizen: ${error.message}`);
        }
    }
}