<script setup lang="ts">
import { onMounted } from 'vue'
import { apiFetch } from '@api'
import { useAuthStore } from '@stores/authStore'
import type { LoginResponse } from '@types'
import router from '@router'

const authStore = useAuthStore()
const toast = useToast()

onMounted(async () => {
    const { data, error } = await apiFetch<LoginResponse>('/session')
    if (data) {
        authStore.setUser(data)
        return
    }

    if (error?.status_code === 500) {
        toast.add({
            title: 'Server error',
            description: error.error || 'An unexpected error occurred. Please try again later.',
            color: 'error'
        })
        return
    }

    authStore.clearUser()
    toast.add({
        title: 'Session expired',
        description: 'Please log in again.',
        color: 'warning'
    })
    router.push({ name: 'login' })
})
</script>

<template>
    <UContainer>
        <div v-if="authStore.isAuthenticated">
            Welcome, {{ authStore.currentUser?.username }}!
        </div>
    </UContainer>
</template>
