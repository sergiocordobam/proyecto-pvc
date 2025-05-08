import { OperatorFetchService } from '../src/operator/services/operator-fetch.service';
import axios from 'axios';

jest.mock('axios');
const mockedAxios = axios as jest.Mocked<typeof axios>;

describe('OperatorFetchService', () => {
  let service: OperatorFetchService;

  beforeEach(() => {
    service = new OperatorFetchService();
    jest.clearAllMocks();
  });

  describe('getOperators', () => {
    it('should return list of operators on success', async () => {
      const mockOperators = [{ id: '1', operatorName: 'Alpha' }];
      mockedAxios.get.mockResolvedValueOnce({ data: mockOperators });

      const result = await service.getOperators();

      expect(mockedAxios.get).toHaveBeenCalledWith('http://localhost:3000/getOperators');
      expect(result).toEqual(mockOperators);
    });

    it('should throw error if axios call fails', async () => {
      mockedAxios.get.mockRejectedValueOnce(new Error('Fetch failed'));

      await expect(service.getOperators()).rejects.toThrow('Fetch failed');
    });
  });

  describe('getOperatorByName', () => {
    it('should return operator matching the name', async () => {
      const mockOperators = [
        { id: '1', operatorName: 'Alpha' },
        { id: '2', operatorName: 'Beta' },
      ];
      mockedAxios.get.mockResolvedValueOnce({ data: mockOperators });

      const result = await service.getOperatorByName('Beta');

      expect(result).toEqual({ id: '2', operatorName: 'Beta' });
    });

    it('should return undefined if no operator matches the name', async () => {
      mockedAxios.get.mockResolvedValueOnce({ data: [] });

      const result = await service.getOperatorByName('Nonexistent');

      expect(result).toBeUndefined();
    });
  });

  describe('getOperatorById', () => {
    it('should return operator matching the id', async () => {
      const mockOperators = [
        { id: '1', operatorName: 'Alpha' },
        { id: '2', operatorName: 'Beta' },
      ];
      mockedAxios.get.mockResolvedValueOnce({ data: mockOperators });

      const result = await service.getOperatorById('1');

      expect(result).toEqual({ id: '1', operatorName: 'Alpha' });
    });

    it('should return undefined if no operator matches the id', async () => {
      mockedAxios.get.mockResolvedValueOnce({ data: [] });

      const result = await service.getOperatorById('999');

      expect(result).toBeUndefined();
    });
  });
});