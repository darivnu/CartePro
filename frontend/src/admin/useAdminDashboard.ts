import { useQuery } from '@tanstack/react-query'                                                                                                                                       
import { getAdminDashboard } from '../api/admin'
import type { GetAdminDashboardParams } from '../api/admin'

export function useAdminDashboard(params: GetAdminDashboardParams = {}) {
const { from, to } = params

return useQuery({
    queryKey: ['admin-dashboard', from, to],
    queryFn: () => getAdminDashboard({ from, to }),
})
}