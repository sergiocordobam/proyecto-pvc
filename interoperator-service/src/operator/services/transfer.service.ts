import { Injectable,HttpException, HttpStatus,Inject, Body } from '@nestjs/common';
import { EventPattern, Payload } from '@nestjs/microservices';
import { ClientProxy, ClientProxyFactory, Transport, MessagePattern} from '@nestjs/microservices';
import { async, firstValueFrom } from 'rxjs';
import { OperatorFetchService } from './operator-fetch.service';
import axios from 'axios';
import { log } from 'console';
import { RegisterCitizenDto } from '../DTO/RegisterCitizenDTO';

@Injectable()
export class TransferService {
    private readonly apiUrl : string;
    private readonly selfId : string;
    private readonly selfName: string;
    constructor(
        @Inject('DELETE_CITIZEN_CLIENT') private readonly deleteCitizenClient: ClientProxy,
        @Inject('DELETE_DOCUMENTS_CLIENT') private readonly deleteDocumentsClient: ClientProxy,
        @Inject('REGISTER_CITIZEN_CLIENT') private readonly registerCitizenClient: ClientProxy,
        @Inject('REGISTER_DOCUMENTS_CLIENT') private readonly registerDocumentsClient: ClientProxy,
        private readonly fetchService: OperatorFetchService,

    ) {
        this.apiUrl = process.env.API_BASE_URL || 'http://localhost:3000';
        this.selfId = process.env.OPERATOR_ID || '1';
        this.selfName = process.env.OPERATOR_NAME || 'Operator1';
    }
     // Listen to the 'transfer_requests' queue
    async processTransfer(@Body() message: any): Promise<void> {
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
            
            const formattedUrls = documentUrls.reduce((acc, url, index) => {
                acc[`URL${index + 1}`] = url;
                return acc;
            }, {});
            // Format the payload to match the required structure
            console.log(`Formatted URLs:`, formattedUrls);
            const payload = {
                id: citizenId,
                citizenName: citizenInfo.name,
                citizenEmail: citizenInfo.email,
                urlDocuments: formattedUrls,
                confirmAPI: process.env.OPERATOR_TRANSFER_ENDPOINT_CONFIRM, 
            };
    
            // Fetch the receiving operator's URL
            const operatorUrl = await this.getOperatorUrl(operatorId);
            console.log(`Fetched operator URL: ${operatorUrl}`);
    
            // Send the payload to the receiving operator's queue via RabbitMQ
            const unregisterCitizen = {
                id: citizenId,
                operatorId: this.selfId,
                operatorName: this.selfName,
            }
            await axios.delete(`${this.apiUrl}/unregisterCitizen`, {
                data: unregisterCitizen,
            });

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
            const { id, req_status } = dto;

            if (req_status === 1) {

                this.deleteCitizenClient.emit('delete_citizen_queue', {
                    citizenId: id,
                });

                // Delete documents via RabbitMQ
                this.deleteDocumentsClient.emit('delete_documents_queue', {
                    citizenId: id
                });
            }
        } catch (error) {
            console.error(`Error confirming transfer: ${error.message}`);
            throw new Error(`Error confirming transfer: ${error.message}`);
        }
    }

    async registerCitizenAndDocuments(request: RegisterCitizenDto): Promise<void> {
        try {
            const { id, citizenName, citizenEmail, urlDocuments,  confirmAPI} = request;

            // Register the citizen in the Auth Service
            this.registerCitizenClient.emit('register_citizen_queue', {
                full_name: citizenName,
                document_id: id,    
                email: citizenEmail,
                password: id,       
                terms_accepted: true
            });                  

            console.log(`Citizen registered successfully`);

            const flatUrls = Object.values(urlDocuments).flat();

            // Send document URLs to the Document Service
            this.registerDocumentsClient.emit('register_documents_queue', {
                citizenId: id,
                documents: flatUrls,
            });
            console.log(`Documents registered successfully`);
            
            const confirmation = {
                id: id, 
                req_status : 1,
            };

            await axios.post(confirmAPI, confirmation);
        } catch (error) {
            const { id, confirmAPI} = request;
            const confirmation = {
                id: id, 
                req_status : 0,
            };

            await axios.post(confirmAPI, confirmation);

            console.error(`Error registering citizen and documents: ${error.message}`);
            throw new Error(`Error registering citizen and documents: ${error.message}`);
        }
    }
}