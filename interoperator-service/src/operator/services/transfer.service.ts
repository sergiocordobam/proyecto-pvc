/*import { Injectable, HttpException, HttpStatus } from '@nestjs/common';
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
}*/

import { Injectable,HttpException, HttpStatus,Inject } from '@nestjs/common';
import { EventPattern, Payload } from '@nestjs/microservices';
import { ClientProxy, ClientProxyFactory, Transport } from '@nestjs/microservices';
import { async, firstValueFrom } from 'rxjs';
import { OperatorFetchService } from './operator-fetch.service';
import axios from 'axios';
import { log } from 'console';
import { RegisterCitizenDto } from '../DTO/RegisterCitizenDTO';

@Injectable()
export class TransferService {

    constructor(
        @Inject('DELETE_CITIZEN_CLIENT') private readonly deleteCitizenClient: ClientProxy,
        @Inject('DELETE_DOCUMENTS_CLIENT') private readonly deleteDocumentsClient: ClientProxy,
        @Inject('REGISTER_CITIZEN_CLIENT') private readonly registerCitizenClient: ClientProxy,
        @Inject('REGISTER_DOCUMENTS_CLIENT') private readonly registerDocumentsClient: ClientProxy,
        private readonly fetchService: OperatorFetchService,
    ) {}
     // Listen to the 'transfer_requests' queue
    async processTransfer(@Payload() message: any): Promise<void> {
        try {
            const { citizenId, operatorId } = message;
    
            // Fetch citizen info via Kong Gateway
            const citizenInfoResponse = await axios.get(`${process.env.AUTH_SERVICE_URL}/citizen-info/${citizenId}`);
            const citizenInfo = citizenInfoResponse.data;
            console.log(`Fetched citizen info:`, citizenInfo);
    
            // Fetch document URLs via Kong Gateway
            const documentUrlsResponse = await axios.get(`${process.env.DOCUMENT_SERVICE_URL}/files/${citizenId}/all`);
            const documentUrls = documentUrlsResponse.data.files;
            console.log(`Fetched document URLs:`, documentUrls);
    
            // Format the payload to match the required structure
            const payload = {
                id: citizenId,
                citizenName: citizenInfo.name,
                citizenEmail: citizenInfo.email,
                urlDocuments: documentUrls,
                confirmAPI: process.env.OPERATOR_TRANSFER_ENDPOINT_CONFIRM, 
            };
    
            // Fetch the receiving operator's URL
            const operatorUrl = await this.getOperatorUrl(operatorId);
            console.log(`Fetched operator URL: ${operatorUrl}`);
    
            // Send the payload to the receiving operator's queue via RabbitMQ

            await axios.post(operatorUrl, payload);
    
        } catch (error) {
            console.error(`Error processing transfer: ${error.message}`);
        }
    }

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

    async confirmTransfer(dto: any): Promise<any> {
        try {
            const { operatorId, transferId, status } = dto;

/*            // Send confirmation to the operator
            const response = await firstValueFrom(
                this.transferConfirmationsClient.send('confirm_transfer', { operatorId, transferId, status }),
            );

            // If the confirmation is successful and the status is "confirmed"
            if (response.status === 200 && status === 'confirmed') {*/
                // Delete citizen via RabbitMQ
            // Delete citizen via RabbitMQ
            this.deleteCitizenClient.emit('delete_citizen', {
                citizenId: transferId,
            });

            // Delete documents via RabbitMQ
            this.deleteDocumentsClient.emit('delete_documents', {
                citizenId: transferId,
            });
        } catch (error) {
            console.error(`Error confirming transfer: ${error.message}`);
            throw new Error(`Error confirming transfer: ${error.message}`);
        }
    }

    async registerCitizenAndDocuments(payload: RegisterCitizenDto): Promise<void> {
        try {
            const { id, citizenName, citizenEmail, urlDocuments,  } = payload;

            // Register the citizen in the Auth Service
            this.registerCitizenClient.emit('register_citizen', {
                full_name: citizenName,
                document_id: id,    
                email: citizenEmail,
                password: id,       
                terms_accepted: true
            });                  

            console.log(`Citizen registered successfully`);

            const flatUrls = Object.values(urlDocuments).flat();

            // Send document URLs to the Document Service
            this.registerDocumentsClient.emit('register_documents', {
                citizenId: id,
                documents: flatUrls,
            });
            console.log(`Documents registered successfully`);
            
        } catch (error) {
            console.error(`Error registering citizen and documents: ${error.message}`);
            throw new Error(`Error registering citizen and documents: ${error.message}`);
        }
    }
}