import { Test, TestingModule } from '@nestjs/testing';
import { DocumentService } from '../src/operator/services/Documents-service-conection';
import axios from 'axios';

jest.mock('axios');
const mockedAxios = axios as jest.Mocked<typeof axios>;

describe('DocumentService', () => {
    let documentService: DocumentService;

    beforeEach(async () => {
        // Mock the environment variable
        process.env.DOCUMENT_SERVICE_URL = 'http://localhost:5000';

        const module: TestingModule = await Test.createTestingModule({
            providers: [DocumentService],
        }).compile();

        documentService = module.get<DocumentService>(DocumentService);
    });

    it('should fetch document URLs', async () => {
        const mockResponse = { data: ['http://example.com/doc1', 'http://example.com/doc2'] };
        mockedAxios.get.mockResolvedValue(mockResponse);

        const result = await documentService.fetchDocumentUrls({ citizenId: '123' });

        expect(mockedAxios.get).toHaveBeenCalledWith(`${process.env.DOCUMENT_SERVICE_URL}/documents/123`);
        expect(result).toEqual(mockResponse.data);
    });

    it('should delete documents', async () => {
        mockedAxios.delete.mockResolvedValue({});

        await documentService.deleteDocuments({ citizenId: '123' });

        expect(mockedAxios.delete).toHaveBeenCalledWith(`${process.env.DOCUMENT_SERVICE_URL}/documents/123`);
    });
});