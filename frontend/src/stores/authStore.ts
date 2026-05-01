import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import type { LoginResponse, LoginRequest, RegisterRequest, RegisterResponse, LogoutResponse } from '@types'
import { apiFetch } from '@utils/api'
import { useRouter } from 'vue-router'

export const useAuthStore = defineStore('auth', () => {
    const toast = useToast()
    const router = useRouter()
    const user = ref<LoginResponse | null>(null)

    const isAuthenticated = computed(() => user.value !== null)
    const currentUser = computed(() => user.value)

    function setUser(userData: LoginResponse) {
        user.value = userData
    }

    function clearUser() {
        user.value = null
    }

    const handleSessionCheck = async (): Promise<boolean> => {
        const { data, error } = await apiFetch<LoginResponse>('/session')

        if (data) {
            setUser(data)
            return true
        }

        if (error?.status_code === 500) {
            toast.add({
                title: 'Server error',
                description: error.error || 'An unexpected error occurred. Please try again later.',
                color: 'error'
            })
            return false
        }

        clearUser()
        toast.add({
            title: 'Session expired',
            description: 'Please log in again.',
            color: 'warning'
        })
        router.push({ name: 'login' })
        return false
    }

    const login = async (req: LoginRequest) => {
        const { data, error } = await apiFetch<LoginResponse>('/login', {
            method: 'POST',
            body: JSON.stringify(req)
        })

        if (data) {
            setUser(data)
        }
        
        return { data, error }
    }

    const register = async (req: RegisterRequest) => {
        const { data, error } = await apiFetch<RegisterResponse>('/users', {
            method: 'POST',
            body: JSON.stringify(req)
        })

        return { data, error }
    }

    const logout = async () => {
        const { error } = await apiFetch<LogoutResponse>('/revoke', {
            method: 'POST',
        })

        // If the server successfully revokes the token, clear the local state
        if (!error) {
            clearUser()
        }

        return { error }
    }

    return {
        user,
        isAuthenticated,
        currentUser,
        setUser,
        clearUser,
        handleSessionCheck,
        login,
        register,
        logout
    }
})