<script lang="ts" setup>
import { useChannelStore } from '@stores/channelStore';
import type { Message } from '@types';
import { computed, onMounted, ref, useTemplateRef } from 'vue';

const channelStore = useChannelStore()
const messages = ref<Message[]>([])

const scrollArea = useTemplateRef('scrollArea')

type VirtualItem =
    | { id: string; type: 'header' }
    | (Message & { type: 'message' })

const virtualItems = computed<VirtualItem[]>(() => {
    return [
        { id: 'header', type: 'header' },
        ...messages.value.map(msg => ({
            ...msg,
            type: 'message' as const
        }))
    ]
})

function scrollToBottom() {
    scrollArea.value?.virtualizer?.scrollToIndex(virtualItems.value.length - 1, { align: 'end' })
}

const getIcon = (type: 0 | 1 | 2) => {
    switch (type) {
        case 0: return 'i-lucide-hash'
        case 1: return 'i-lucide-volume-2'
    }
}

onMounted(() => {
    scrollToBottom()
})
</script>

<template>
    <UScrollArea :items="virtualItems"
                 virtualize
                 orientation="vertical"
                 class="hover-scrollbar h-full w-full border rounded-lg border-neutral-200 dark:border-neutral-800 bg-white dark:bg-neutral-900 shadow-sm"
                 :ui="{ viewport: 'gap-4 p-4' }">

        <template #default="{ item }">

            <!-- Welcome Header -->
            <div v-if="item.type === 'header' && channelStore.currentChannel"
                 class="flex flex-col gap-2 p-4 border-b border-neutral-100 dark:border-neutral-800 mb-4">
                <UIcon :name="getIcon(channelStore.currentChannel?.type || 0)"
                       class="size-20 rounded-full p-3 bg-white dark:bg-neutral-800 shadow-sm" />
                <p class="select-none text-2xl font-bold">
                    Welcome to
                    <UIcon :name="getIcon(channelStore.currentChannel?.type || 0)"
                           class="inline mb-0.5 me-0" />{{ channelStore.currentChannel?.name }}!
                </p>
                <p class="text-muted">
                    This is the beginning of the
                    <UIcon :name="getIcon(channelStore.currentChannel?.type || 0)"
                           class="inline mb-0.5 me-0 max-w-4" /><strong>{{ channelStore.currentChannel?.name }}</strong>
                    channel.
                </p>
            </div>

            <!-- Messages -->
            <div v-else-if="item.type === 'message'" class="flex flex-col gap-1 py-1">
                <span class="font-bold text-sm">{{ item.author }}</span>
                <p class="text-neutral-600 dark:text-neutral-300">{{ item.content }}</p>
            </div>
        </template>
    </UScrollArea>
</template>