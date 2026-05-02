<script lang="ts" setup>
import type { FormSubmitEvent, RadioGroupItem } from '@nuxt/ui';
import { useChannelStore } from '@stores/channelStore';
import { useServerStore } from '@stores/serverStore';
import type { ChannelRequest } from '@types';
import { reactive, ref } from 'vue';
import * as z from 'zod'

const serverStore = useServerStore()
const channelStore = useChannelStore()

const schema = z.object({
    name: z.string('Name is required').min(3, 'Must be at least 3 characters').max(100, 'Must be at most 100 characters'),
    type: z.literal(0).or(z.literal(1)).or(z.literal(2)).default(0),
})
type Schema = z.output<typeof schema>

const state = reactive<Partial<Schema>>({
    type: 0,
    name: undefined,
})

const open = defineModel<boolean>('open', { default: false })
const types = ref<RadioGroupItem[]>([
    {
        id: 0,
        label: 'Text',
        description: 'A text channel for chatting and sharing messages.',
        class: 'cursor-pointer mb-1',
    },
    {
        id: 1,
        label: 'Voice',
        description: 'A voice channel for audio communication.',
        class: 'cursor-pointer mb-1',
    },
    {
        id: 2,
        label: 'Category',
        description: 'A category for organizing channels.',
        disabled: true,
        class: 'cursor-pointer',
    },
])

const onSubmit = async (payload: FormSubmitEvent<Schema>) => {
    const req: ChannelRequest = {
        name: payload.data.name,
        type: payload.data.type,
        position: serverStore.currentServer?.channels.length || 0,
    }

    const created = await channelStore.createChannel(req)
    if (created) open.value = false
}

</script>

<template>
    <div class="flex flex-col flex-1 w-full gap-4">
        <UForm :state="state" class="w-full space-y-4" @submit="onSubmit">
            <UFormField label="Type" name="type" required>
                <URadioGroup class="w-full" variant="card" v-model="state.type" value-key="id" :items="types" />
            </UFormField>

            <UFormField label="Name" name="name" required>
                <template v-if="state.type === 0">
                    <UInput v-model="state.name" class="w-full" icon="i-lucide-hash" placeholder="Channel Name" />
                </template>
                <template v-else-if="state.type === 1">
                    <UInput v-model="state.name" class="w-full" icon="i-lucide-volume-2" placeholder="Channel Name" />
                </template>
            </UFormField>

            <div class="flex flex-1 w-full gap-4 justify-between">
                <UButton class="w-full justify-center" label="Cancel" variant="subtle" color="neutral" @click="open = false" />
                <UButton type="submit" class="w-full justify-center" label="Create" variant="solid" color="primary" />
            </div>
        </UForm>
    </div>
</template>