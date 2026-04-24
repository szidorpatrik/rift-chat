<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import type { FormSubmitEvent, AuthFormField } from '@nuxt/ui'
import * as z from 'zod'

const router = useRouter()
const toast = useToast()
const fields = ref<AuthFormField[]>([
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

function onSubmit(payload: FormSubmitEvent<Schema>) {
    const data = JSON.stringify(payload.data)
    fetch('http://localhost:8080/api/v1/login', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: data
    }).then(async res => {
        const result = await res.json()
        if (res.ok) {
            toast.add({
                title: 'Login successful',
                description: `Welcome back, ${result.username}!`,
            })
            router.push({ name: 'app' })
        } else {
            toast.add({
                title: 'Login failed',
                description: result["error"] || 'An unexpected error occurred',
            })
        }
    }).catch(error => {
        toast.add({
            title: 'Network Error',
            description: error.message,
        })
    })
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
        </UPageCard>
    </UContainer>
</template>
