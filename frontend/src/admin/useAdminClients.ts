import { useMutation, useQuery } from '@tanstack/react-query'
import {
  getAdminClients,
  getAdminClientDetail,
  createTopup,
  cancelTransaction,
} from '../api/admin'
import type { GetAdminClientsParams } from '../api/admin'
import { queryClient } from '../api/queryClient'

export function useAdminClients(params: GetAdminClientsParams = {}) {
  const { employerId, page = 1, limit = 20 } = params

  return useQuery({
    queryKey: ['admin-clients', employerId, page, limit],
    queryFn: () => getAdminClients({ employerId, page, limit }),
  })
}

export function useAdminClientDetail(id: number) {
  return useQuery({
    queryKey: ['admin-client', id],
    queryFn: () => getAdminClientDetail(id),
  })
}

export function useCreateTopup() {
  return useMutation({
    mutationFn: createTopup,
    onSuccess: (_data, variables) => {
      queryClient.invalidateQueries({ queryKey: ['admin-client', variables.client_id] })
      queryClient.invalidateQueries({ queryKey: ['admin-clients'] })
    },
  })
}

export function useCancelTransaction() {
  return useMutation({
    mutationFn: (variables: { id: number; reason: string; clientId: number }) =>
      cancelTransaction(variables.id, variables.reason),
    onSuccess: (_data, variables) => {
      queryClient.invalidateQueries({ queryKey: ['admin-client', variables.clientId] })
    },
  })
}
