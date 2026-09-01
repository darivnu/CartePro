export type Role = 'client' | 'partner' | 'admin';

export interface User {
    id: string;
    email: string;
    role: Role;
}