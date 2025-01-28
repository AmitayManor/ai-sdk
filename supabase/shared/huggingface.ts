import { HfInference } from "https://esm.sh/@huggingface/inference@3.1.2"
import { ModelError, ProcessingError, TimeoutError, ValidationError } from './errors.ts'
import { EdgeEnv, HFModelConfig } from './types.ts'

export class HuggingFaceClient {
    private client: HfInference
    private env: EdgeEnv
    private retryCount: number = 0

    constructor(env: EdgeEnv) {
        this.env = env
        this.client = new HfInference(this.env.HF_API_TOKEN)
    }

    private withTimeout<T>(promise: Promise<T>): Promise<T> {
        const timeout = new Promise<never>((_, reject) => {
            setTimeout(() => {
                reject(new TimeoutError())
            }, this.env.MODEL_TIMEOUT)
        })

        return Promise.race([promise, timeout])
    }

    private async withRetry<T>(operation: () => Promise<T>): Promise<T> {
        try {
            return await this.withTimeout(operation())
        } catch (error) {
            if (this.retryCount >= this.env.MAX_RETRIES || error instanceof TimeoutError) {
                throw error
            }

            this.retryCount++
            console.warn(`Retry attempt ${this.retryCount} for operation`)
            return this.withRetry(operation)
        }
    }

    async textToText(modelId: string, input: string, config?: HFModelConfig['parameters']): Promise<string> {
        try {
            const response = await this.withRetry(() =>
                this.client.textGeneration({
                    model: modelId,
                    inputs: input,
                    parameters: {
                        max_new_tokens: 512,
                        temperature: 0.7,
                        top_p: 0.95,
                        return_full_text: false,
                        ...config
                    }
                })
            )

            return response.generated_text
        } catch (error) {
            throw new ProcessingError(`Text generation failed: ${error}`)
        }
    }

    async textToImage(modelId: string, prompt: string, config?: HFModelConfig['parameters']): Promise<Uint8Array> {
        try {
            const response = await this.withRetry(() =>
                this.client.textToImage({
                    model: modelId,
                    inputs: prompt,
                    parameters: {
                        negative_prompt: "blurry, bad quality, distorted",
                        num_inference_steps: 50,
                        guidance_scale: 7.5,
                        ...config
                    }
                })
            )

            return new Uint8Array(await response.arrayBuffer())
        } catch (error) {
            throw new ProcessingError(`Image generation failed: ${error}`)
        }
    }

    async validateModel(modelId: string): Promise<void> {
        try {
            const response = await fetch(
                `https://huggingface.co/api/models/${modelId}`,
                {
                    headers: {
                        Authorization: `Bearer ${this.env.HF_API_TOKEN}`
                    }
                }
            )

            if (!response.ok) {
                throw new ModelError('Model not found or inaccessible', 404)
            }

            const data = await response.json()
            if (!data.pipeline_tag) {
                throw new ValidationError('Invalid model type')
            }
        } catch (error) {
            if (error instanceof ModelError) {
                throw error
            }
            throw new ProcessingError(`Model validation failed: ${error}`)
        }
    }
}