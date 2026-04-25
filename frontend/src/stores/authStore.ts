import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import type { LoginResponse } from '@types'

export const useAuthStore = defineStore('auth', () => {
    const user = ref<LoginResponse | null>(null)

    const isAuthenticated = computed(() => user.value !== null)
    const currentUser = computed(() => user.value)

    function setUser(userData: LoginResponse) {
        user.value = userData
    }

    function clearUser() {
        user.value = null
    }

    return {
        user,
        isAuthenticated,
        currentUser,
        setUser,
        clearUser
    }
})
