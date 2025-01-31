import { serve } from 'https://deno.land/std@0.168.0/http/server.ts';
import { RequestProcessor } from './processor';
import { ModelRequest } from './types';

const processor = new RequestProcessor(
    Deno.env.get('_SUPABASE_URL')!,
    Deno.env.get('_SUPABASE_ANON_KEY')!,
    Deno.env.get('HF_API_TOKEN')!
);

serve(async (req: Request) => {
    try {
        if (req.method !== 'POST') {
            return new Response(
                JSON.stringify({ error: 'Method not allowed' }),
                { status: 405, headers: { 'Content-Type': 'application/json' } }
            );
        }

        const request = await req.json() as ModelRequest;
        if (!request.id || !request.modelType || !request.input || !request.userId) {
            return new Response(
                JSON.stringify({ error: 'Invalid request format' }),
                { status: 400, headers: { 'Content-Type': 'application/json' } }
            );
        }

        const result = await processor.processRequest(request);

        return new Response(
            JSON.stringify(result),
            {
                status: result.status === 'success' ? 200 : 500,
                headers: { 'Content-Type': 'application/json' }
            }
        );

    } catch (error) {
        return new Response(
            JSON.stringify({ error: error.message }),
            { status: 500, headers: { 'Content-Type': 'application/json' } }
        );
    }
});