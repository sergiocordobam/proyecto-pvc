import { Test, TestingModule } from '@nestjs/testing';
import { OperatorRegistrationService } from '../src/operator/services/operator-registration.service';
import axios from 'axios';

jest.mock('axios');
const mockedAxios = axios as jest.Mocked<typeof axios>;

describe('OperatorRegistrationService', () => {
    let registrationService: OperatorRegistrationService;

    beforeEach(async () => {
        const module: TestingModule = await Test.createTestingModule({
            providers: [
                OperatorRegistrationService,
                {
                    provide: 'API_URL',
                    useValue: process.env.API_BASE_URL,
                },
            ],
        }).compile();

        registrationService = module.get<OperatorRegistrationService>(OperatorRegistrationService);
    });

    it('should register an operator', async () => {
        const mockDto = { name: 'Operator1', address: 'Address1', contactMail: 'test@example.com', participants: [] };
        const mockResponse = { id: '1' };
        mockedAxios.post.mockResolvedValue({ data: mockResponse });

        const result = await registrationService.registerOperator(mockDto);

        expect(mockedAxios.post).toHaveBeenCalledWith(`${process.env.API_BASE_URL}/registerOperator`, mockDto);
        expect(result).toEqual(mockResponse);
    });

    it('should register an endpoint', async () => {
        const endpointData = { idOperator: '1', endPoint: 'http://example.com', endPointConfirm: 'http://example.com/confirm' };
        mockedAxios.post.mockResolvedValue({ data: {} });

        await registrationService.registerEndPoint(endpointData);

        expect(mockedAxios.post).toHaveBeenCalledWith(`${process.env.API_BASE_URL}/registerTransferEndPoint`, endpointData);
    });
});