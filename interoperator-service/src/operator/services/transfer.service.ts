import { Injectable, HttpException, HttpStatus } from '@nestjs/common';
import axios from 'axios';
import { OperatorFetchService } from './operator-fetch.service';
import { TransferCitizenDto } from '../DTO/TransferCitizenDto';
import { ConfirmTransferDto } from '../DTO/ConfirmTransferDto';

@Injectable()
export class TransferService {
    constructor(private readonly fetchService: OperatorFetchService) {}

    private async getOperatorUrl(operatorId: string): Promise<string> {
        const operator = await this.fetchService.getOperatorById(operatorId);

        if (!operator || !operator.transferAPIURL) {
        throw new HttpException(
            `Interop URL not found for operator with ID: ${operatorId}`,
            HttpStatus.NOT_FOUND,
        );
        }

        return operator.transferAPIURL;
    }

    private async fetchCitizenInfo(citizenId: string): Promise<any> {
        try {
            const authServiceUrl = process.env.AUTH_SERVICE_URL || 'http://localhost:4000';
            const response = await axios.get(`${authServiceUrl}/citizens/${citizenId}`);
            return response.data;
        } catch (error) {
            throw new HttpException(
                `Error fetching citizen info: ${error.message}`,
                HttpStatus.INTERNAL_SERVER_ERROR,
            );
        }
    }

    private async fetchDocumentUrls(citizenId: string): Promise<any> {
        try {
            const documentServiceUrl = process.env.DOCUMENT_SERVICE_URL || 'http://localhost:5000';
            const response = await axios.get(`${documentServiceUrl}/documents/${citizenId}`);
            return response.data;
        } catch (error) {
            throw new HttpException(
                `Error fetching document URLs: ${error.message}`,
                HttpStatus.INTERNAL_SERVER_ERROR,
            );
        }
    }

    async transferCitizen(dto: TransferCitizenDto): Promise<any> {
        try {
            // Fetch citizen info from Auth Service
            const citizenInfo = await this.fetchCitizenInfo(dto.citizenId);

            // Fetch document URLs from Document Service
            const documentUrls = await this.fetchDocumentUrls(dto.citizenId);

            // Combine data
            const payload = {
                id: dto.citizenId,
                citizenName: citizenInfo.name,
                citizenEmail: citizenInfo.email,
                urlDocuments: documentUrls,
            };

            // Get the target operator's URL
            const interopUrl = await this.getOperatorUrl(dto.operatorId);

            // Send the combined data to the target operator
            const response = await axios.post(`${interopUrl}/transfer`, payload);
            return response.data;
        } catch (error) {
            throw new HttpException(
                `Error transferring citizen: ${error.message}`,
                HttpStatus.INTERNAL_SERVER_ERROR,
            );
        }
    }


    async confirmTransfer(dto: ConfirmTransferDto): Promise<any> {
        try {
        const interopUrl = await this.getOperatorUrl(dto.operatorId);
        const response = await axios.post(`${interopUrl}/confirm`, {
            transferId: dto.transferId,
            status: dto.status,
        });
        return response.data;
        } catch (error) {
        throw new HttpException(
            `Error confirming transfer: ${error.message}`,
            HttpStatus.INTERNAL_SERVER_ERROR,
        );
        }
    }
}
