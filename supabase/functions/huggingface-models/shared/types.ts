export interface ModelRequest {
    id: string
    input: string | Record<string, unknown>
    modelId: string
}

export interface ModelResponse {
    output: string | Uint8Array
    processingTime: number
    tokenCount?: number
    status: 'success' | 'error'
    error?: string
}

export interface EdgeEnv {
    SUPABASE_URL: string
    SUPABASE_ANON_KEY: string
    HF_API_TOKEN: string
    MODEL_TIMEOUT: number
    MAX_RETRIES: number
}

export type ModelType = 'text-to-text' | 'text-to-image'

export interface HFModelConfig {
    type: ModelType
    parameters?: Record<string, unknown>
}