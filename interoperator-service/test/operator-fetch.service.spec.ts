import { Test, TestingModule } from '@nestjs/testing';
import { OperatorFetchService } from '../src/operator/services/operator-fetch.service';
import axios from 'axios';

jest.mock('axios');
const mockedAxios = axios as jest.Mocked<typeof axios>;

describe('OperatorFetchService', () => {
    let fetchService: OperatorFetchService;

    beforeEach(async () => {
        // Mock the environment variable
        process.env.API_BASE_URL = 'http://localhost:3000';

        const module: TestingModule = await Test.createTestingModule({
            providers: [OperatorFetchService],
        }).compile();

        fetchService = module.get<OperatorFetchService>(OperatorFetchService);
    });

    it('should fetch all operators', async () => {
        const mockOperators = [{ id: '1', name: 'Operator1' }];
        mockedAxios.get.mockResolvedValue({ data: mockOperators });

        const result = await fetchService.getOperators();

        expect(mockedAxios.get).toHaveBeenCalledWith(`${process.env.API_BASE_URL}/getOperators`);
        expect(result).toEqual(mockOperators);
    });

    it('should fetch operator by name', async () => {
        const mockOperators = [{ id: '1', operatorName: 'Operator1' }];
        jest.spyOn(fetchService, 'getOperators').mockResolvedValue(mockOperators);

        const result = await fetchService.getOperatorByName('Operator1');

        expect(result).toEqual(mockOperators[0]);
    });

    it('should fetch operator by ID', async () => {
        const mockOperators = [{ id: '1', operatorName: 'Operator1' }];
        jest.spyOn(fetchService, 'getOperators').mockResolvedValue(mockOperators);

        const result = await fetchService.getOperatorById('1');

        expect(result).toEqual(mockOperators[0]);
    });
});