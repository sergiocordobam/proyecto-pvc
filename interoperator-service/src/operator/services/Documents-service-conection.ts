import { Injectable, HttpException, HttpStatus } from '@nestjs/common';
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
}