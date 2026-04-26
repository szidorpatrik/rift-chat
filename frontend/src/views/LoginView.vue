<script setup lang="ts">
import { reactive } from 'vue'
import { useRouter } from 'vue-router'
import { apiFetch } from '@api'
import { useAuthStore } from '@stores/authStore'
import * as z from 'zod'
import type { FormSubmitEvent, AuthFormField } from '@nuxt/ui'
import type { LoginRequest, LoginResponse } from '@types'

const router = useRouter()
const toast = useToast()
const authStore = useAuthStore()

const fields = reactive<AuthFormField[]>([
    {
        name: 'username',
        type: 'text',
        label: 'Username',
        required: true,
    },
    {
        name: 'password',
        type: 'password',
        label: 'Password',
        required: true,
    }
])
const schema = z.object({
    username: z.string('Username is required').min(3, 'Must be at least 3 characters'),
    password: z.string('Password is required').min(8, 'Must be at least 8 characters')
})
type Schema = z.output<typeof schema>

const onSubmit = async (payload: FormSubmitEvent<Schema>) => {
    const req: LoginRequest = {
        username: payload.data.username,
        password: payload.data.password,
    }

    const { data, error } = await apiFetch<LoginResponse>('/login', {
        method: 'POST',
        body: JSON.stringify(req)
    })

    if (error) {
        toast.add({
            title: 'Login failed',
            description: error.error,
            color: 'error'
        })
        return
    }

    authStore.setUser(data!)
    toast.add({ title: 'Login successful', description: `Welcome back, ${data?.username}!` })
    router.push({ name: 'app' })
}
</script>

<template>
    <UContainer class="w-screen h-screen flex flex-col items-center justify-center gap-4 p-4">
        <UPageCard class="w-full max-w-md">
            <UAuthForm :schema="schema"
                       :fields="fields"
                       title="Login"
                       class="max-w-md"
                       icon="i-lucide-user"
                       @submit="onSubmit" />

            <div class="text-sm mt-1 flex flex-row gap-1">
                <span class="text-muted">Don't have an account?</span>
                <ULink :to="{ name: 'register' }" class="text-secondary">Register</ULink>
            </div>
        </UPageCard>
    </UContainer>
</template>
