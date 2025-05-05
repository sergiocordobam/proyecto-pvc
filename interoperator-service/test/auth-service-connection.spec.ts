import { Test, TestingModule } from '@nestjs/testing';
import { CitizenService } from '../src/operator/services/Auth-service-Conection';
import axios from 'axios';

jest.mock('axios');
const mockedAxios = axios as jest.Mocked<typeof axios>;

describe('CitizenService', () => {
    let citizenService: CitizenService;

    beforeEach(async () => {
        // Mock the environment variable
        process.env.AUTH_SERVICE_URL = 'http://localhost:4000';

        const module: TestingModule = await Test.createTestingModule({
            providers: [CitizenService],
        }).compile();

        citizenService = module.get<CitizenService>(CitizenService);
    });

    it('should fetch citizen info', async () => {
        const mockResponse = { data: { id: '123', name: 'John Doe' } };
        mockedAxios.get.mockResolvedValue(mockResponse);

        const result = await citizenService.fetchCitizenInfo({ citizenId: '123' });

        expect(mockedAxios.get).toHaveBeenCalledWith(`${process.env.AUTH_SERVICE_URL}/citizens/123`);
        expect(result).toEqual(mockResponse.data);
    });

    it('should delete a citizen', async () => {
        const mockResponse = { status: 200 };
        mockedAxios.delete.mockResolvedValue(mockResponse);

        await citizenService.deleteCitizen({ document_id: '5555555555' });

        expect(mockedAxios.delete).toHaveBeenCalledWith(`${process.env.AUTH_SERVICE_URL}/delete_user`, {
            data: { document_id: '5555555555' },
        });
    });

    it('should handle errors when deleting a citizen', async () => {
        mockedAxios.delete.mockRejectedValue(new Error('Test error'));

        await expect(citizenService.deleteCitizen({ document_id: '5555555555' })).rejects.toThrow(
            'Error deleting user: Test error',
        );

        expect(mockedAxios.delete).toHaveBeenCalledWith(`${process.env.AUTH_SERVICE_URL}/delete_user`, {
            data: { document_id: '5555555555' },
        });
    });

    it('should throw an error if document_id is missing', async () => {
        await expect(citizenService.deleteCitizen({ document_id: '' })).rejects.toThrow(
            'Document ID is required',
        );
    });
});