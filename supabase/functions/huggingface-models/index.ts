// Follow this setup guide to integrate the Deno language server with your editor:
// https://deno.land/manual/getting_started/setup_your_environment
// This enables autocomplete, go to definition, etc.

import { serve } from "https://deno.land/std@0.214.0/http/server.ts"
import { createClient } from "https://esm.sh/@supabase/supabase-js@2.48.1"
import { HuggingFaceClient } from './shared/huggingface.ts'
import { handleError } from './shared/errors.ts'
import { EdgeEnv, ModelRequest, ModelResponse } from './shared/types.ts'
import "jsr:@supabase/functions-js/edge-runtime.d.ts"

// Get environment variables
const env: EdgeEnv = {
  SUPABASE_URL: Deno.env.get('SUPABASE_URL') ?? '',
  SUPABASE_ANON_KEY: Deno.env.get('SUPABASE_ANON_KEY') ?? '',
  HF_API_TOKEN: Deno.env.get('HF_API_TOKEN') ?? '',
  MODEL_TIMEOUT: Number(Deno.env.get('MODEL_TIMEOUT') ?? 30000),
  MAX_RETRIES: Number(Deno.env.get('MAX_RETRIES') ?? 3)
}

serve(async (req: Request) => {
  try {
    // Parse request
    const { id, input, modelId } = await req.json() as ModelRequest

    // Initialize clients
    const hfClient = new HuggingFaceClient(env)
    const supabase = createClient(env.SUPABASE_URL, env.SUPABASE_ANON_KEY)

    // Validate input
    if (!input || typeof input !== 'string') {
      throw new Error('Invalid input: must be a non-empty string')
    }

    // Validate model
    await hfClient.validateModel(modelId)

    // Update request status
    await supabase
        .from('model_requests')
        .update({ status: 'IN_PROGRESS' })
        .eq('id', id)

    // Process request
    const startTime = Date.now()
    const result = await hfClient.textToText(modelId, input)
    const processingTime = Date.now() - startTime

    // Prepare response
    const response: ModelResponse = {
      output: result,
      processingTime,
      status: 'success',
      tokenCount: result.split(/\s+/).length // Basic token count
    }

    // Update request with results
    await supabase
        .from('model_requests')
        .update({
          status: 'COMPLETED',
          output_data: response,
          completed_at: new Date().toISOString(),
          processing_time: processingTime,
          token_count: response.tokenCount
        })
        .eq('id', id)

    return new Response(
        JSON.stringify(response),
        {
          headers: { 'Content-Type': 'application/json' }
        }
    )
  } catch (error) {
    return handleError(error)
  }
})

/* To invoke locally:

  1. Run `supabase start` (see: https://supabase.com/docs/reference/cli/supabase-start)
  2. Make an HTTP request:

  curl -i --location --request POST 'http://127.0.0.1:54321/functions/v1/huggingface-models' \
    --header 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJzdXBhYmFzZS1kZW1vIiwicm9sZSI6ImFub24iLCJleHAiOjE5ODM4MTI5OTZ9.CRXP1A7WOeoJeXxjNni43kdQwgnWNReilDMblYTn_I0' \
    --header 'Content-Type: application/json' \
    --data '{"name":"Functions"}'

*/
