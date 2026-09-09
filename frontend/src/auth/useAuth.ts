import { useMutation, useQuery } from '@tanstack/react-query';
import { getCurrentUser, login, logout, registerClient } from '../api/auth';
import { queryClient } from '../api/queryClient';

export function useAuth() {
const query = useQuery({
    queryKey: ['me'],
    queryFn: getCurrentUser,
    retry: false,
})

if (query.isError || query.data === null) {
    return { status: 'unauthenticated' as const, user: null }
}

if (query.data === undefined) {
    return { status: 'loading' as const, user: null }
}

return { status: 'authenticated' as const, user: query.data }
}


export function useLogin() {
    return useMutation({
        mutationFn: ({ email, password }: {email: string, password: string}) => login(email, password),
        onSuccess: (data) => {
            queryClient.setQueryData(['me'], data.user)
        },
    })
}

export function useLogout() {
    return useMutation({
        mutationFn: logout,
        onSuccess: () => {
            queryClient.setQueryData(['me'], null)
        },
    })
}

export function useRegisterClient() {
    return useMutation({
        mutationFn: ({ email, password, name }: { email: string; password: string; name: string }) => registerClient({ email, password, name }),
        onSuccess: async () => {
            const user = await getCurrentUser()
            queryClient.setQueryData(['me'], user)
        },
    })
}