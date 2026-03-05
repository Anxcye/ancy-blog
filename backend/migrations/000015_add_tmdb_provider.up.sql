-- Add TMDB integration provider
INSERT INTO integration_providers (provider_key, provider_type, name, description, enabled, config_json, created_at, updated_at)
VALUES (
    'tmdb',
    'llm',
    'TMDB',
    'The Movie Database API for fetching movie and TV show metadata',
    false,
    '{"api_key": ""}',
    NOW(),
    NOW()
)
ON CONFLICT (provider_key) DO NOTHING;
