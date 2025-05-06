/*import { Injectable, HttpException, HttpStatus } from '@nestjs/common';
import axios from 'axios';

@Injectable()
export class DocumentService {
    private readonly documentServiceUrl = process.env.DOCUMENT_SERVICE_URL || 'http://localhost:5000';

    async fetchDocumentUrls(citizenId: string): Promise<any> {
        try {
            const response = await axios.get(`${this.documentServiceUrl}/documents/${citizenId}`);
            return response.data;
        } catch (error) {
            throw new HttpException(
                `Error fetching document URLs: ${error.message}`,
                HttpStatus.INTERNAL_SERVER_ERROR,
            );
        }
    }

    async deleteDocuments(citizenId: string): Promise<void> {
        try {
            await axios.delete(`${this.documentServiceUrl}/documents/${citizenId}`);
        } catch (error) {
            throw new HttpException(
                `Error deleting documents: ${error.message}`,
                HttpStatus.INTERNAL_SERVER_ERROR,
            );
        }
    }
}*/
import { Injectable } from '@nestjs/common';
import { EventPattern, Payload } from '@nestjs/microservices';
import axios from 'axios';

@Injectable()
export class DocumentService {
    private readonly documentServiceUrl = process.env.DOCUMENT_SERVICE_URL || 'http://localhost:5000';

    @EventPattern('fetch_document_urls')
    async fetchDocumentUrls(@Payload() data: { citizenId: string }): Promise<any> {
        try {
            console.log(`Fetching document URLs for citizen ID: ${data.citizenId}`);
            const response = await axios.get(`${this.documentServiceUrl}/documents/${data.citizenId}`);
            return response.data;
        } catch (error) {
            console.error(`Error fetching document URLs: ${error.message}`);
            throw new Error(`Error fetching document URLs: ${error.message}`);
        }
    }

    @EventPattern('delete_documents')
    async deleteDocuments(@Payload() data: { citizenId: string }): Promise<void> {
        try {
            console.log(`Deleting documents for citizen ID: ${data.citizenId}`);
            await axios.delete(`${this.documentServiceUrl}/documents/${data.citizenId}`);
        } catch (error) {
            console.error(`Error deleting documents: ${error.message}`);
            throw new Error(`Error deleting documents: ${error.message}`);
        }
    }
}