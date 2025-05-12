import { Injectable,HttpException, HttpStatus,Inject, Body } from '@nestjs/common';
import { ClientProxy} from '@nestjs/microservices';
import { OperatorFetchService } from './operator-fetch.service';
import axios from 'axios';
import { RegisterCitizenDto } from '../DTO/RegisterCitizenDTO';
import { PubSubService } from './PubSubService';
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
        private readonly PubSubService: PubSubService,

    ) {
        this.apiUrl = process.env.API_BASE_URL || 'http://localhost:3000';
        this.selfId = process.env.OPERATOR_ID || '1';
        this.selfName = process.env.OPERATOR_NAME || 'Operator1';
    }
     // Listen to the 'transfer_requests' queue
    async processTransfer(@Body() message: any): Promise<void> {
        try {
            const { citizenId, operatorId, citizenName, citizenEmail } = message;
            console.log("payload citizen recibido");
            console.log("id", citizenId);
            console.log("name", citizenName);
            console.log("email", citizenEmail);
    
            // Fetch document URLs via Kong Gateway
            const documentUrlsResponse = await axios.get(`${process.env.DOCUMENT_SERVICE_URL}/files/download/${citizenId}/all`);
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
                citizenName: citizenName,
                citizenEmail: citizenEmail,
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
            console.log("unregistered citizen from govcarpeta")

            await axios.post(operatorUrl, payload);
            console.log("operador transferido")
    
        } catch (error) {
            console.error(`Error processing transfer: ${error.message}`);
        }
    }

    private async getOperatorUrl(operatorId: string): Promise<string> {
        const operator = await this.fetchService.getOperatorById(operatorId);
        console.log("url de trasnferencia obtenida")

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

                console.log("delete citizen from operador pvc")

                // Delete documents via RabbitMQ
                this.deleteDocumentsClient.emit('delete_documents_queue', {
                    citizenId: id
                });
                console.log("delete documents from citizen")
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

                    // Enviar evento a través de Pub/Sub
            const topicName = process.env.GCP_PUBSUB_TOPIC_EMAIL || 'email-topic';
            const externalUrl = 'https://example.com/document';
            const messageData = {
            event: 'transfer',
                user: id,
                name: citizenName,
                user_email: citizenEmail,
                extra_data: {
                    URL: externalUrl,
                    Asunto:'Confirmacion de datos para transferencia a Carpeta PVC',
                    Body:'Tus datos han sido transferidos a la Carpeta PVC, necesitamos confirmar ciertos datos y que nos des tu nueva contraseña, muchas graciaspor seleccionarnos como tu operador de confianza',
                },
            };
            await this.PubSubService.publishMessage(topicName, messageData);
            console.log(`Evento enviado exitosamente a través de Pub/Sub`);
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