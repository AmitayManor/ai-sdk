-- Enable extensions
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Enable real-time
BEGIN;
  DROP PUBLICATION IF EXISTS supabase_realtime CASCADE;
  CREATE PUBLICATION supabase_realtime;
COMMIT;

-- AI Models Table
CREATE TABLE IF NOT EXISTS ai_models (
                                         id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(255) NOT NULL,
    model_type VARCHAR(50) NOT NULL,
    version VARCHAR(50) NOT NULL,
    huggingface_id VARCHAR(255) NOT NULL,
    function_url TEXT NOT NULL,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
                             CONSTRAINT valid_model_type CHECK (model_type IN ('text-to-text', 'text-to-image'))
    );

-- Users Table
CREATE TABLE IF NOT EXISTS public.users (
                                            id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    last_login TIMESTAMP WITH TIME ZONE,
                             is_admin BOOLEAN DEFAULT false,
                             is_active BOOLEAN DEFAULT true,
                             failed_login_attempts INTEGER DEFAULT 0
                             );


-- API Keys Table
CREATE TABLE IF NOT EXISTS api_keys (
                                        id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES auth.users(id),
    key_hash TEXT NOT NULL UNIQUE,
    name VARCHAR(255) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    last_used TIMESTAMP WITH TIME ZONE,
                            is_active BOOLEAN DEFAULT true,
                            rate_limit INTEGER NOT NULL DEFAULT 60,
                            CONSTRAINT rate_limit_range CHECK (rate_limit BETWEEN 1 AND 1000)
    );

-- Model Requests Table
CREATE TABLE IF NOT EXISTS model_requests (
                                              id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES auth.users(id),
    model_id UUID NOT NULL REFERENCES ai_models(id),
    api_key_id UUID NOT NULL REFERENCES api_keys(id),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    completed_at TIMESTAMP WITH TIME ZONE,
                               status VARCHAR(50) NOT NULL,
    input_data JSONB NOT NULL,
    output_data JSONB,
    error_msg TEXT,
    token_used INTEGER DEFAULT 0,
    token_count INTEGER DEFAULT 0,
    processing_time INTEGER DEFAULT 0,
    CONSTRAINT valid_status CHECK (status IN ('PENDING', 'IN_PROGRESS', 'COMPLETED', 'FAILED'))
    );



-- Indexes
CREATE INDEX idx_model_requests_user_id ON model_requests(user_id);
CREATE INDEX idx_model_requests_model_id ON model_requests(model_id);
CREATE INDEX idx_model_requests_status ON model_requests(status);
CREATE INDEX idx_model_requests_created_at ON model_requests(created_at);
CREATE INDEX idx_api_keys_user_id ON api_keys(user_id);
CREATE INDEX idx_api_keys_key_hash ON api_keys(key_hash);
CREATE INDEX idx_ai_models_type ON ai_models(model_type) WHERE is_active = true;
CREATE INDEX IF NOT EXISTS users_email_idx ON public.users(email);

-- Enable Row Level Security
ALTER TABLE public.users ENABLE ROW LEVEL SECURITY;
ALTER TABLE ai_models ENABLE ROW LEVEL SECURITY;
ALTER TABLE model_requests ENABLE ROW LEVEL SECURITY;
ALTER TABLE api_keys ENABLE ROW LEVEL SECURITY;

-- Add tables to realtime publication
ALTER PUBLICATION supabase_realtime ADD TABLE model_requests;