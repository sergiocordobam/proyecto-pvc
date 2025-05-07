import { TransferService } from '../src/operator/services/transfer.service';
import { OperatorFetchService } from '../src/operator/services/operator-fetch.service';
import axios from 'axios';

jest.mock('axios');
const mockedAxios = axios as jest.Mocked<typeof axios>;

describe('TransferService', () => {
  let service: TransferService;

  const mockDeleteCitizenClient = { emit: jest.fn() };
  const mockDeleteDocumentsClient = { emit: jest.fn() };
  const mockRegisterCitizenClient = { emit: jest.fn() };
  const mockRegisterDocumentsClient = { emit: jest.fn() };

  const mockFetchService = {
    getOperatorById: jest.fn().mockResolvedValue({ transferAPIURL: 'http://operator-url.com' }),
  } as unknown as OperatorFetchService;

  beforeEach(() => {
    jest.clearAllMocks();

    service = new TransferService(
      mockDeleteCitizenClient as any,
      mockDeleteDocumentsClient as any,
      mockRegisterCitizenClient as any,
      mockRegisterDocumentsClient as any,
      mockFetchService
    );
  });

  it('should fetch info and post to operator successfully', async () => {
    mockedAxios.get.mockImplementation((url) => {
      if (url.includes('citizen-info')) {
        return Promise.resolve({ data: { name: 'John', email: 'john@email.com' } });
      }
      if (url.includes('documents')) {
        return Promise.resolve({ data: ['doc1.pdf', 'doc2.pdf'] });
      }
      return Promise.reject(new Error('Unknown URL'));
    });

    mockedAxios.post.mockResolvedValue({ data: { success: true } });

    const dto = { citizenId: '123', operatorId: '456' };

    await service.processTransfer(dto);

    expect(mockedAxios.get).toHaveBeenCalledWith(expect.stringContaining('citizen-info/123'));
    expect(mockedAxios.get).toHaveBeenCalledWith(expect.stringContaining('documents/123'));
    expect(mockedAxios.post).toHaveBeenCalledWith(
      'http://operator-url.com',
      expect.objectContaining({
        id: '123',
        citizenName: 'John',
        citizenEmail: 'john@email.com',
        urlDocuments: ['doc1.pdf', 'doc2.pdf'],
      })
    );
  });
});