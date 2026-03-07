-- Add TMDB integration provider
INSERT INTO integration_providers (provider_key, provider_type, name, enabled, config_json, meta_json, created_at, updated_at)
VALUES (
    'tmdb',
    'llm',
    'TMDB',
    false,
    '{"api_key": ""}',
    '{"description": "The Movie Database API for fetching movie and TV show metadata"}',
    NOW(),
    NOW()
)
ON CONFLICT (provider_key) DO NOTHING;
