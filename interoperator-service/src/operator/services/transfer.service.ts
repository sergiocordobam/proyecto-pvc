import { Injectable, HttpException, HttpStatus } from '@nestjs/common';
import axios from 'axios';
import { OperatorFetchService } from './operator-fetch.service';
import { CitizenService } from './Auth-service-Conection';
import { DocumentService } from './Documents-service-conection';
import { TransferCitizenDto } from '../DTO/TransferCitizenDto';
import { ConfirmTransferDto } from '../DTO/ConfirmTransferDto';
import { TransferRequestDto } from '../DTO/TransferRequestDto';

@Injectable()
export class TransferService {
    constructor(
        private readonly fetchService: OperatorFetchService,
        private readonly citizenService: CitizenService,
        private readonly documentService: DocumentService,
    ) {}

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

    async transferCitizen(dto: TransferCitizenDto): Promise<any> {
        try {
            // Fetch citizen info and document URLs
            const citizenInfo = await this.citizenService.fetchCitizenInfo(dto.citizenId);
            const documentUrls = await this.documentService.fetchDocumentUrls(dto.citizenId);

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

            // If the confirmation is successful and the status is "confirmed"
            if (response.status === 200 && dto.status === 'confirmed') {
                await this.citizenService.deleteCitizen(dto.transferId);
                await this.documentService.deleteDocuments(dto.transferId);
            }

            return response.data;
        } catch (error) {
            throw new HttpException(
                `Error confirming transfer: ${error.message}`,
                HttpStatus.INTERNAL_SERVER_ERROR,
            );
        }
    }
    async processTransfer(dto: TransferRequestDto): Promise<any> {
        try {
            const userInfo = await this.citizenService.fetchCitizenInfo(dto.userId);
            const userDocuments = await this.documentService.fetchDocumentUrls(dto.userId);

            const confirmation = {
                userId: dto.userId,
                userInfo,
                userDocuments,
            };

            return {
                success: true,
                message: 'Transfer processed successfully',
                confirmation,
            };
        } catch (error) {
            throw new HttpException(
                `Error processing transfer: ${error.message}`,
                HttpStatus.INTERNAL_SERVER_ERROR,
            );
        }
    }
}