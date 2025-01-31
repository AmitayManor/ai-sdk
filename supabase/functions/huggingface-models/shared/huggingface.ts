
import { ProcessingResult, MODEL_CONFIG } from './types';

export class HuggingFaceClient {
    private readonly apiKey: string;
    private readonly headers: HeadersInit;

    constructor(apiKey: string) {
        this.apiKey = apiKey;
        this.headers = {
            'Authorization': `Bearer ${apiKey}`,
            'Content-Type': 'application/json'
        };
    }

    async processText(input: string): Promise<ProcessingResult> {
        try {
            const config = MODEL_CONFIG.text2text;
            const response = await fetch(config.endpoint, {
                method: 'POST',
                headers: this.headers,
                body: JSON.stringify({
                    inputs: input,
                    parameters: config.parameters
                })
            });

            if (!response.ok) {
                throw new Error(`HuggingFace API error: ${response.statusText}`);
            }

            const result = await response.json();
            return {
                status: 'success',
                data: {
                    output: result.generated_text,
                    processingTime: Date.now(),
                    tokenCount: result.usage?.total_tokens
                }
            };
        } catch (error) {
            return {
                status: 'error',
                error: error.message
            };
        }
    }

    async processImage(input: string): Promise<ProcessingResult> {
        try {
            const config = MODEL_CONFIG.text2image;
            const response = await fetch(config.endpoint, {
                method: 'POST',
                headers: this.headers,
                body: JSON.stringify({
                    inputs: input,
                    parameters: config.parameters
                })
            });

            if (!response.ok) {
                throw new Error(`HuggingFace API error: ${response.statusText}`);
            }

            // Image response comes as binary data
            const imageData = new Uint8Array(await response.arrayBuffer());

            return {
                status: 'success',
                data: {
                    output: imageData,
                    processingTime: Date.now()
                }
            };
        } catch (error) {
            return {
                status: 'error',
                error: error.message
            };
        }
    }
}