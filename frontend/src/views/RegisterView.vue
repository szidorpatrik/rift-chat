<script setup lang="ts">
import { reactive } from 'vue'
import { useRouter } from 'vue-router'
import type { FormSubmitEvent, AuthFormField } from '@nuxt/ui'
import * as z from 'zod'
import { useAuthStore } from '@stores/authStore'

const router = useRouter()
const toast = useToast()
const authStore = useAuthStore()

const fields = reactive<AuthFormField[]>([
    {
        name: 'username',
        type: 'text',
        label: 'Username',
        required: true,
        autofocus: true,
    },
    {
        name: 'email',
        type: 'email',
        label: 'Email',
        required: true,
    },
    {
        name: 'password',
        type: 'password',
        label: 'Password',
        required: true,
    },
    {
        name: 'confirmPassword',
        type: 'password',
        label: 'Confirm Password',
        required: true,

    }
])

const schema = z.object({
    username: z.string('Username is required').min(3, 'Must be at least 3 characters'),
    email: z.email('Invalid email address'),
    password: z.string('Password is required').min(8, 'Must be at least 8 characters'),
    confirmPassword: z.string()
}).refine((data) => data.password === data.confirmPassword, {
    message: "Passwords don't match",
    path: ["confirmPassword"],
});

type Schema = z.output<typeof schema>

const onSubmit = async (payload: FormSubmitEvent<Schema>) => {
    const { data, error } = await authStore.register({
        username: payload.data.username,
        email: payload.data.email,
        password: payload.data.password,
    })

    if (error) {
        toast.add({
            title: 'Registration failed',
            description: error.error,
            color: 'error'
        })
        return
    }

    toast.add({ title: 'Registration successful', description: `You can now log in ${data?.username}!` })
    router.push({ name: 'login' })
}

</script>

<template>
    <UContainer class="w-screen h-screen flex flex-col items-center justify-center gap-4 p-4">
        <UPageCard class="w-full max-w-md">
            <UAuthForm :schema="schema"
                       :fields="fields"
                       title="Registration"
                       class="max-w-md"
                       icon="i-lucide-user-plus"
                       @submit="onSubmit" />

            <div class="text-sm flex flex-row gap-1 items-center">
                <span class="text-muted">Already have an account?</span>
                <ULink :to="{ name: 'login' }" class="text-secondary py-2 pr-3">Log in</ULink>
            </div>
        </UPageCard>
    </UContainer>
</template>