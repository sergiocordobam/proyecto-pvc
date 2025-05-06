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
    async deleteCitizen(@Payload() data: { document_id: string }): Promise<void> {
        const { document_id } = data;
    
        if (!document_id) {
            console.error('Document ID is missing in the payload');
            throw new Error('Document ID is required');
        }
    
        try {
            console.log(`Deleting user with document ID: ${document_id}`);
            const response = await axios.delete(`${this.authServiceUrl}/delete_user`, {
                data: { document_id },
            });
    
            if (response.status === 200) {
                console.log(`User with document ID ${document_id} deleted successfully`);
            } else {
                console.error(`Failed to delete user with document ID ${document_id}. Status: ${response.status}`);
                throw new Error(`Failed to delete user with document ID ${document_id}`);
            }
        } catch (error) {
            console.error(`Error deleting user with document ID ${document_id}: ${error.message}`);
            throw new Error(`Error deleting user: ${error.message}`);
        }
    }
}