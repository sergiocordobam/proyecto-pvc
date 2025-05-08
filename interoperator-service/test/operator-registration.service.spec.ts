import { OperatorRegistrationService } from '../src/operator/services/operator-registration.service';
import axios from 'axios';
import { RegisterOperatorDto } from '../src/operator/DTO/RegisterOperatorDto';
import { RegisterEndpointDto } from '../src/operator/DTO/RegisterEndpointDto';
jest.mock('axios');
const mockedAxios = axios as jest.Mocked<typeof axios>;

describe('OperatorRegistrationService', () => {
  let service: OperatorRegistrationService;

  beforeEach(() => {
    service = new OperatorRegistrationService();
    jest.clearAllMocks();
  });

  describe('registerOperator', () => {
    it('should post data to /registerOperator and return response', async () => {
      const dto: RegisterOperatorDto = { name: 'Alpha', address:'test', contactMail:'a@a.com', participants: ['a', 'b'] };
      const mockResponse = { success: true, id: 'op-123' };
      mockedAxios.post.mockResolvedValueOnce({ data: mockResponse });

      const result = await service.registerOperator(dto);

      expect(mockedAxios.post).toHaveBeenCalledWith('http://localhost:3000/registerOperator', dto);
      expect(result).toEqual(mockResponse);
    });

    it('should throw error when registration fails', async () => {
      mockedAxios.post.mockRejectedValueOnce(new Error('Post failed'));

      await expect(service.registerOperator({ name: 'Beat', address:'test1', contactMail:'b@b.com', participants: ['c', 'd'] } as RegisterOperatorDto))
        .rejects
        .toThrow('Error registering operator: Post failed');
    });
  });

  describe('registerEndPoint', () => {
    it('should post data to /registerTransferEndPoint and return response', async () => {
      const endpointData = { idOperator: '123', endPoint: 'http://endpoint.com',endPointConfirm:'http://endpoint.com/confirm' } as RegisterEndpointDto; ;
      const mockResponse = { status: 'ok' };
      mockedAxios.post.mockResolvedValueOnce({ data: mockResponse });

      const result = await service.registerEndPoint(endpointData);

      expect(mockedAxios.post).toHaveBeenCalledWith('http://localhost:3000/registerTransferEndPoint', endpointData);
      expect(result).toEqual(mockResponse);
    });
  });
});