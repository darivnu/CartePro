import { useQuery, useMutation } from "@tanstack/react-query";
import { getPartners, getPartner, collectPayment, getOwnTransactions, getOwnDashboard } from "../api/partners";
import type { GetPartnersParams, GetOwnTransactionsParams, GetOwnDashboardParams } from "../api/partners";

export function usePartners(params: GetPartnersParams = {}) {
    const { search, category, page = 1, limit = 20 } = params;

    return useQuery({
        queryKey: ['partners', search, category, page, limit],
        queryFn: () => getPartners({ search, category, page, limit }),
    });
}

export function usePartner(id: string) {
    return useQuery({
        queryKey: ['partner', id],
        queryFn: () => getPartner(id),
    });
}

export function useCollectPayment() {
    return useMutation({
        mutationFn: collectPayment,
    });
}

export function useOwnTransactions(params: GetOwnTransactionsParams = {}) {
    const { page = 1, limit = 20, from, to } = params;

    return useQuery({
        queryKey: ['own-transactions', page, limit, from, to],
        queryFn: () => getOwnTransactions({ page, limit, from, to }),
    });
}

export function useOwnDashboard(params: GetOwnDashboardParams = {}) {
    const { from, to } = params;

    return useQuery({
        queryKey: ['own-dashboard', from, to],
        queryFn: () => getOwnDashboard({ from, to }),
    });
}