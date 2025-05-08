// __tests__/citizen.service.spec.ts
import { Test, TestingModule } from '@nestjs/testing';
import { CitizenService } from '../src/operator/services/Auth-service-Conection';
import axios from 'axios';

jest.mock('axios');
const mockedAxios = axios as jest.Mocked<typeof axios>;

describe('CitizenService', () => {
  let service: CitizenService;

  beforeEach(async () => {
    const module: TestingModule = await Test.createTestingModule({
      providers: [CitizenService],
    }).compile();

    service = module.get<CitizenService>(CitizenService);
  });

  describe('fetchCitizenInfo', () => {
    it('should return citizen info from axios', async () => {
      const mockData = { name: 'Jane Doe', email: 'jane@example.com' };
      mockedAxios.get.mockResolvedValueOnce({ data: mockData });

      const result = await service.fetchCitizenInfo({ citizenId: '123' });
      expect(mockedAxios.get).toHaveBeenCalledWith('http://localhost:4000/citizens/123');
      expect(result).toEqual(mockData);
    });
  });

  describe('deleteCitizen', () => {
    it('should send delete request via axios', async () => {
      mockedAxios.delete.mockResolvedValueOnce({ status: 200 });

      await service.deleteCitizen({ document_id: '123' });
      expect(mockedAxios.delete).toHaveBeenCalledWith(
        'http://localhost:4000/delete_user',
        { data: { document_id: '123' } },
      );
    });

    it('should throw if document_id is missing', async () => {
      await expect(service.deleteCitizen({ document_id: '' })).rejects.toThrow('Document ID is required');
    });
  });
});