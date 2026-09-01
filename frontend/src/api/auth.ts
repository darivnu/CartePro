import { apiFetch } from './client'
import type { User } from '../types/auth'

export function login(email: string, password: string) {
    return apiFetch<{ user: User }>('/auth/login', {
      method: 'POST',
      body: { email, password },
    })
  }

export function logout() {
    return apiFetch<void>('/auth/logout', {
      method: 'POST',
    })
  }

export function getCurrentUser() {
    return apiFetch<User>('/auth/me')
}
