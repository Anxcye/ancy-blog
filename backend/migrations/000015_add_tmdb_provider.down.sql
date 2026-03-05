-- Remove TMDB integration provider
DELETE FROM integration_providers WHERE provider_key = 'tmdb';
