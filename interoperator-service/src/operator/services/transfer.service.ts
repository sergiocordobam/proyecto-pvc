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

import { Injectable } from '@nestjs/common';
import { EventPattern, Payload } from '@nestjs/microservices';
import { ClientProxy, ClientProxyFactory, Transport } from '@nestjs/microservices';
import { firstValueFrom } from 'rxjs';

@Injectable()
export class TransferService {
    private readonly transferConfirmationsClient: ClientProxy;
    private readonly fetchCitizenInfoClient: ClientProxy;
    private readonly fetchDocumentUrlsClient: ClientProxy;
    private readonly deleteCitizenClient: ClientProxy;
    private readonly deleteDocumentsClient: ClientProxy;
    private readonly registerCitizenClient: ClientProxy;
    private readonly registerDocumentsClient: ClientProxy;

    constructor() {
        const rabbitMQUrl = process.env.RABBITMQ_URL || 'amqp://localhost:5672';

        this.transferConfirmationsClient = ClientProxyFactory.create({
            transport: Transport.RMQ,
            options: {
                urls: [rabbitMQUrl],
                queue: 'transfer_confirmations',
                queueOptions: { durable: true },
            },
        });

        this.fetchCitizenInfoClient = ClientProxyFactory.create({
            transport: Transport.RMQ,
            options: {
                urls: [rabbitMQUrl],
                queue: 'fetch_citizen_info',
                queueOptions: { durable: true },
            },
        });

        this.fetchDocumentUrlsClient = ClientProxyFactory.create({
            transport: Transport.RMQ,
            options: {
                urls: [rabbitMQUrl],
                queue: 'fetch_document_urls',
                queueOptions: { durable: true },
            },
        });

        this.deleteCitizenClient = ClientProxyFactory.create({
            transport: Transport.RMQ,
            options: {
                urls: [rabbitMQUrl],
                queue: 'delete_citizen',
                queueOptions: { durable: true },
            },
        });

        this.deleteDocumentsClient = ClientProxyFactory.create({
            transport: Transport.RMQ,
            options: {
                urls: [rabbitMQUrl],
                queue: 'delete_documents',
                queueOptions: { durable: true },
            },
        });

        this.registerCitizenClient = ClientProxyFactory.create({
            transport: Transport.RMQ,
            options: {
                urls: [rabbitMQUrl],
                queue: 'register_citizen',
                queueOptions: { durable: true },
            },
        });

        this.registerDocumentsClient = ClientProxyFactory.create({
            transport: Transport.RMQ,
            options: {
                urls: [rabbitMQUrl],
                queue: 'register_documents',
                queueOptions: { durable: true },
            },
        });
    }
    @EventPattern('transfer_requests') // Listen to the 'transfer_requests' queue
    async processTransfer(@Payload() message: any): Promise<void> {
        try {
            const { citizenId, operatorId } = message;

            // Fetch citizen info via RabbitMQ
            const citizenInfo = await firstValueFrom(
                this.fetchCitizenInfoClient.send('fetch_citizen_info', { citizenId }),
            );
            console.log(`Fetched citizen info:`, citizenInfo);

            // Fetch document URLs via RabbitMQ
            const documentUrls = await firstValueFrom(
                this.fetchDocumentUrlsClient.send('fetch_document_urls', { citizenId }),
            );
            console.log(`Fetched document URLs:`, documentUrls);

            // Format the payload to match the required structure
            const payload = {
                id: citizenId,
                citizenName: citizenInfo.name,
                citizenEmail: citizenInfo.email,
                urlDocuments: documentUrls, // Dynamically include all document URLs
            };

            // Simulate processing logic
            console.log(`Processing transfer for citizen ${citizenId} to operator ${operatorId}`);
            console.log(`Payload to be sent:`, payload);

            // Publish confirmation to the transfer_confirmations queue
            const confirmation = {
                transferId: citizenId,
                status: 'confirmed',
            };
            await firstValueFrom(this.transferConfirmationsClient.emit('transfer_confirmations', confirmation));
        } catch (error) {
            console.error(`Error processing transfer: ${error.message}`);
        }
    }

    async confirmTransfer(dto: any): Promise<any> {
        try {
            const { operatorId, transferId, status } = dto;

            // Simulate fetching the operator's URL (this logic can be replaced with actual implementation)
            const operatorUrl = `http://example.com/operators/${operatorId}`;

            // Send confirmation to the operator
            const response = await firstValueFrom(
                this.transferConfirmationsClient.send('confirm_transfer', { operatorId, transferId, status }),
            );

            // If the confirmation is successful and the status is "confirmed"
            if (response.status === 200 && status === 'confirmed') {
                // Delete citizen via RabbitMQ
                await firstValueFrom(this.deleteCitizenClient.send('delete_citizen', { citizenId: transferId }));

                // Delete documents via RabbitMQ
                await firstValueFrom(this.deleteDocumentsClient.send('delete_documents', { citizenId: transferId }));
            }

            return response;
        } catch (error) {
            console.error(`Error confirming transfer: ${error.message}`);
            throw new Error(`Error confirming transfer: ${error.message}`);
        }
    }

    async registerCitizenAndDocuments(payload: any): Promise<void> {
        try {
            const { id, citizenName, citizenEmail, urlDocuments } = payload;

            // Register the citizen in the Auth Service
            const registerCitizenResponse = await firstValueFrom(
                this.registerCitizenClient.send('register_citizen', {
                    id: id,
                    name: citizenName,
                    email: citizenEmail,
                }),
            );
            console.log(`Citizen registered successfully:`, registerCitizenResponse);

            // Send document URLs to the Document Service
            const registerDocumentsResponse = await firstValueFrom(
                this.registerDocumentsClient.send('register_documents', {
                    citizenId: id,
                    documents: urlDocuments,
                }),
            );
            console.log(`Documents registered successfully:`, registerDocumentsResponse);
        } catch (error) {
            console.error(`Error registering citizen and documents: ${error.message}`);
            throw new Error(`Error registering citizen and documents: ${error.message}`);
        }
    }
}