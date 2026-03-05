// File: api/modules/integrations.ts
// Purpose: Provide integration-center API calls for provider listing, update, and health test.
// Module: frontend-admin/api/integrations, domain gateway layer.
// Related: SystemView and backend admin integration handlers.
import { httpClient } from '@/api/http';
import type { ApiEnvelope, IntegrationProvider } from '@/api/types';

interface ProviderTestResult {
  ok: boolean;
  message: string;
  latencyMs: number;
}

export async function listIntegrations(): Promise<IntegrationProvider[]> {
  const response = await httpClient.get<ApiEnvelope<IntegrationProvider[]>>('/admin/integrations');
  return response.data.data;
}

export async function updateIntegration(
  providerKey: string,
  payload: {
    enabled: boolean;
    configJson: Record<string, unknown>;
    metaJson: Record<string, unknown>;
  },
): Promise<void> {
  await httpClient.put<ApiEnvelope<boolean>>(`/admin/integrations/${providerKey}`, payload);
}

export async function testIntegration(providerKey: string): Promise<ProviderTestResult> {
  const response = await httpClient.post<ApiEnvelope<ProviderTestResult>>(`/admin/integrations/${providerKey}/test`);
  return response.data.data;
}
