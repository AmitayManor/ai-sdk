import { createClient } from '@supabase/supabase-js';
import { HuggingFaceClient } from './huggingface';
import { ModelRequest, ProcessingResult } from './types';

export class RequestProcessor {
    private supabase;
    private hfClient: HuggingFaceClient;

    constructor(supabaseUrl: string, supabaseKey: string, hfApiKey: string) {
        this.supabase = createClient(supabaseUrl, supabaseKey);
        this.hfClient = new HuggingFaceClient(hfApiKey);
    }

    private async updateRequestStatus(
        requestId: string,
        status: string,
        data?: any,
        error?: string
    ) {
        const updates = {
            status,
            ...(data && { output_data: data }),
            ...(error && { error_msg: error }),
            ...(data?.processingTime && { processing_time: data.processingTime }),
            ...(data?.tokenCount && { token_count: data.tokenCount })
        };

        await this.supabase
            .from('model_requests')
            .update(updates)
            .match({ id: requestId });
    }

    private async storeImage(userId: string, requestId: string, imageData: Uint8Array) {
        const filePath = `${userId}/${requestId}.png`;

        await this.supabase.storage
            .from('generated-images')
            .upload(filePath, imageData, {
                contentType: 'image/png',
                upsert: true
            });

        return filePath;
    }

    async processRequest(request: ModelRequest): Promise<ProcessingResult> {
        try {
            // Update status to processing
            await this.updateRequestStatus(request.id, 'processing');

            // Process based on model type
            const result = request.modelType === 'text2text'
                ? await this.hfClient.processText(request.input)
                : await this.hfClient.processImage(request.input);

            if (result.status === 'error') {
                await this.updateRequestStatus(request.id, 'failed', null, result.error);
                return result;
            }

            // For image responses, store the image and update the output with the file path
            if (request.modelType === 'text2image' && result.data) {
                const imagePath = await this.storeImage(
                    request.userId,
                    request.id,
                    result.data.output as Uint8Array
                );
                result.data.output = imagePath;
            }

            // Update request with results
            await this.updateRequestStatus(request.id, 'completed', result.data);

            return result;
        } catch (error) {
            await this.updateRequestStatus(request.id, 'failed', null, error.message);
            return {
                status: 'error',
                error: error.message
            };
        }
    }
}