export type ApiKeyData = {
  id: string
  name: string
  keyPrefix: string
  enabled: boolean
  createdBy: string
  createdAt: string
  lastUsedAt?: string
}

export type CreateApiKeyResponse = {
  id: string
  name: string
  key: string
  keyPrefix: string
  createdAt: string
}
