import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import type { LoginResponse } from '@types'
import { apiFetch } from '@/utils/api'
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

    return {
        user,
        isAuthenticated,
        currentUser,
        setUser,
        clearUser,
        handleSessionCheck,
    }
})
