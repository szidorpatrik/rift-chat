<script setup lang="ts">
import * as z from 'zod'
import type { FormSubmitEvent } from '@nuxt/ui'
import { reactive } from 'vue'
import { apiFetch } from '@utils/api'
import type { ServerResponse } from '@types'
import { useServerStore } from '@stores/serverStore'

const schema = z.object({
    invite_code: z.string().max(32, 'Must be at most 32 characters').optional(),
})

type Schema = z.output<typeof schema>

const state = reactive<Partial<Schema>>({
    invite_code: undefined,
})

const toast = useToast()
const serverStore = useServerStore()

async function onSubmit(payload: FormSubmitEvent<Schema>) {
    toast.add({ title: 'TODO: Joining server...', description: 'This feature is not implemented yet.', color: 'warning' })
    return

    const { data, error } = await apiFetch<ServerResponse>('/servers', {
        method: 'POST',
        body: JSON.stringify({
            invite_code: payload.data.invite_code,
        }),
    })

    if (error) {
        toast.add({ title: 'Error', description: error!.error, color: 'error' })
        return
    }

    serverStore.addServer(data!)
    toast.add({ title: 'Success', description: 'Server created successfully.', color: 'success' })
}
</script>

<template>
    <div class="flex flex-col flex-1 gap-6 overflow-auto overflow-y-hidden">

        <UForm :schema="schema" :state="state" class="flex flex-1 gap-6 items-end w-full" @submit="onSubmit">
            <UFormField name="invite_code" class="w-full">
                <UInput v-model="state.invite_code" placeholder="Enter invite code" class="w-full" />
            </UFormField>
            <UButton class="max-h-fit" type="submit" trailing-icon="i-lucide-log-in">
                Join
            </UButton>
        </UForm>

        <USeparator />

        <div class="flex flex-1 flex-wrap justify-center gap-6 overflow-y-auto">
            <template v-for="server in serverStore.publicServers" :key="server.id">
                <div class="w-60 rounded-lg bg-neutral-200 dark:bg-neutral-800 flex flex-col items-center">
                    <div class="w-full h-30 mb-3 select-none pointer-events-none">
                        <img
                             class="w-full h-full object-cover rounded-t-lg"
                             :src="server.icon_url ? server.icon_url : 'https://placehold.co/200x120?text=Server+Background'"
                             alt="server background">
                    </div>

                    <div class="flex flex-1 flex-col w-full gap-3 justify-between items-start mb-3 px-3 overflow-x-auto break-all">
                        <div>
                            <p class="truncate max-w-50 mb-1">{{ server.name }}</p>
                            <p class="text-sm text-muted">{{ server.description }}</p>
                        </div>

                        <template v-if="serverStore.servers.some(s => s.id === server.id)">
                            <UButton disabled icon="i-lucide-check" class="w-full justify-center"></UButton>
                        </template>
                        <template v-else>
                            <UButton class="w-full justify-center"
                                     trailing-icon="i-lucide-log-in"
                                     @click="serverStore.joinPublicServer(server.id)">
                                Join
                            </UButton>
                        </template>
                    </div>
                </div>
            </template>
        </div>
    </div>
</template>
