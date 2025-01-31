export interface ModelRequest {
    id: string
    input: string
    modelType: 'text2text' | 'text2image'
    userId: string
}

export interface ModelResponse {
    output: string | Uint8Array
    processingTime: number
    tokenCount?: number
}

export interface ProcessingResult {
    status: 'success' | 'error';
    data?: ModelResponse;
    error?: string;
}

export const MODEL_CONFIG = {
    text2text: {
        endpoint: 'https://api-inference.huggingface.co/models/google/gemma-2-2b-it',
        parameters: {
            max_new_tokens: 512,
            temperature: 0.7,
            top_p: 0.95,
            return_full_text: false,
        }
    },
    text2image: {
        endpoint: 'https://api-inference.huggingface.co/models/black-forest-labs/FLUX.1-dev',
        parameters: {
            negative_prompt: "blurry, bad quality, distorted",
            num_inference_steps: 50,
            guidance_scale: 7.5,
        }
    }
} as const;

export interface EdgeEnv {
    SUPABASE_URL: string
    SUPABASE_ANON_KEY: string
    HF_API_TOKEN: string
    MODEL_TIMEOUT: number
    MAX_RETRIES: number
}


//TODO: Delete those:
export type ModelType = 'text-to-text' | 'text-to-image'

export interface HFModelConfig {
    type: ModelType
    parameters?: Record<string, unknown>
}