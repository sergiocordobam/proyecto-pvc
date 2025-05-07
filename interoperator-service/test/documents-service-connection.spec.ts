import { DocumentService } from '../src/operator/services/Documents-service-conection';
import axios from 'axios';

jest.mock('axios');
const mockedAxios = axios as jest.Mocked<typeof axios>;

describe('DocumentService', () => {
  let service: DocumentService;

  beforeEach(() => {
    service = new DocumentService();
    jest.clearAllMocks();
  });

  describe('fetchDocumentUrls', () => {
    it('should return document URLs when axios call is successful', async () => {
      const mockData = { urls: ['doc1.pdf', 'doc2.pdf'] };
      mockedAxios.get.mockResolvedValueOnce({ data: mockData });

      const result = await service.fetchDocumentUrls({ citizenId: '123' });

      expect(mockedAxios.get).toHaveBeenCalledWith('http://localhost:5000/documents/123');
      expect(result).toEqual(mockData);
    });

    it('should throw an error if axios call fails', async () => {
      mockedAxios.get.mockRejectedValueOnce(new Error('Network error'));

      await expect(service.fetchDocumentUrls({ citizenId: '123' }))
        .rejects
        .toThrow('Error fetching document URLs: Network error');
    });
  });

  describe('deleteDocuments', () => {
    it('should call axios.delete with correct URL', async () => {
      mockedAxios.delete.mockResolvedValueOnce({});

      await service.deleteDocuments({ citizenId: '456' });

      expect(mockedAxios.delete).toHaveBeenCalledWith('http://localhost:5000/documents/456');
    });

    it('should throw an error if axios.delete fails', async () => {
      mockedAxios.delete.mockRejectedValueOnce(new Error('Delete failed'));

      await expect(service.deleteDocuments({ citizenId: '456' }))
        .rejects
        .toThrow('Error deleting documents: Delete failed');
    });
  });
});