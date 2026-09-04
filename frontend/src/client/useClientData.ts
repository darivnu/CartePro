import { useMutation, useQuery } from "@tanstack/react-query";
import { getBalance, getTransactions, generateQrCode } from "../api/clients";

export function useBalance(options: { refetchInterval?: number | false } = {}) {
    return useQuery({
        queryKey: ['balance'],
        queryFn: getBalance,
        refetchInterval: options.refetchInterval,
    });
}

export function useTransactions(page = 1, limit = 20) {
    return useQuery({
        queryKey: ['transactions', page, limit],
        queryFn: () => getTransactions(page, limit),
    });
}

export function useGenerateQrCode() {
    return useMutation({
        mutationFn: generateQrCode,
    });
}
