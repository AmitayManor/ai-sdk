-- AI Models policies
CREATE POLICY "Public models are viewable by everyone"
    ON ai_models FOR SELECT
                                USING (is_active = true);

CREATE POLICY "Only admins can insert"
    ON ai_models FOR INSERT
    WITH CHECK ((SELECT is_admin FROM public.users WHERE id = auth.uid()) = true);

CREATE POLICY "Only admins can update"
    ON ai_models FOR UPDATE
                                       USING ((SELECT is_admin FROM public.users WHERE id = auth.uid()) = true);

-- Model Requests policies
CREATE POLICY "Users can view own requests"
    ON model_requests FOR SELECT
                                     USING (auth.uid() = user_id);

CREATE POLICY "Admins can view all requests"
    ON model_requests FOR SELECT
                                     USING ((SELECT is_admin FROM public.users WHERE id = auth.uid()) = true);

CREATE POLICY "Users can create requests"
    ON model_requests FOR INSERT
    WITH CHECK (auth.uid() = user_id);

-- API Keys policies
CREATE POLICY "Users can view own keys"
    ON api_keys FOR SELECT
                                      USING (auth.uid() = user_id);

CREATE POLICY "Users can create own keys"
    ON api_keys FOR INSERT
    WITH CHECK (auth.uid() = user_id);

CREATE POLICY "Users can update own keys"
    ON api_keys FOR UPDATE
                                      USING (auth.uid() = user_id);

CREATE POLICY "Users can delete own keys"
    ON api_keys FOR DELETE
USING (auth.uid() = user_id);

-- Add helper functions
CREATE OR REPLACE FUNCTION check_api_key_rate_limit(key_id UUID)
RETURNS boolean AS $$
DECLARE
rate_limit INTEGER;
    current_count INTEGER;
BEGIN
SELECT api_keys.rate_limit INTO rate_limit
FROM api_keys
WHERE id = key_id AND is_active = true;

IF rate_limit IS NULL THEN
        RETURN false;
END IF;

SELECT COUNT(*) INTO current_count
FROM model_requests
WHERE api_key_id = key_id
  AND created_at > NOW() - INTERVAL '1 minute';

RETURN current_count < rate_limit;
END;
$$ LANGUAGE plpgsql SECURITY DEFINER;