<script setup lang="ts">
import * as z from 'zod'
import type { FormSubmitEvent } from '@nuxt/ui'
import { reactive } from 'vue'
import { useServerStore } from '@stores/serverStore'
import type { ServerRequest } from '@types'

const open = defineModel<boolean>('open', { default: false })

const schema = z.object({
    name: z.string('Name is required').min(3, 'Must be at least 3 characters').max(100, 'Must be at most 100 characters'),
    description: z.string().max(255, 'Must be at most 255 characters').optional(),
    isPublic: z.boolean().optional().default(false),
})

type Schema = z.output<typeof schema>

const state = reactive<Partial<Schema>>({
    name: undefined,
    isPublic: false,
    description: undefined,
})

const toast = useToast()
const serverStore = useServerStore()

async function onSubmit(payload: FormSubmitEvent<Schema>) {
    const req: ServerRequest = {
        name: payload.data.name!,
        description: payload.data.description,
        is_public: payload.data.isPublic,
    }
    const { error } = await serverStore.createServer(req)

    if (error) {
        toast.add({ title: 'Error', description: error.error, color: 'error' })
        return
    }

    toast.add({ title: 'Success', description: 'Server created successfully.', color: 'success' })
    open.value = false
}
</script>

<template>
    <UForm :schema="schema" :state="state" class="w-full space-y-4" @submit="onSubmit">
        <UFormField label="Name" name="name" required>
            <UInput v-model="state.name" class="w-full" />
        </UFormField>

        <UFormField label="Description" name="description">
            <UInput v-model="state.description" class="w-full" />
        </UFormField>

        <UFormField name="isPublic">
            <UCheckbox v-model="state.isPublic" label="Public" />
        </UFormField>

        <div class="w-full flex flex-1 justify-end">
            <UButton type="submit">
                Create
            </UButton>
        </div>
    </UForm>
</template>
