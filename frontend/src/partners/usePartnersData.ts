import { useQuery } from "@tanstack/react-query";
import { getPartners, getPartner } from "../api/partners";
import type { GetPartnersParams } from "../api/partners";

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
