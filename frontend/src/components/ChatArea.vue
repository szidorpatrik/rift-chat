<script setup lang="ts">
import { computed, onMounted, ref, useTemplateRef } from 'vue'
import type { Message } from '@types'
import { useChannelStore } from '@stores/channelStore';

const props = defineProps<{
    messages: Message[]
}>()

const channelStore = useChannelStore()
const input = ref('')

const scrollArea = useTemplateRef('scrollArea')

type VirtualItem =
    | { id: string; type: 'header' }
    | (Message & { type: 'message' })

const virtualItems = computed<VirtualItem[]>(() => {
    return [
        { id: 'header', type: 'header' },
        ...props.messages.map(msg => ({
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
    <div class="flex flex-1 w-full gap-4">
        <div class="flex flex-1 flex-col justify-between overflow-hidden rounded-lg">
            <!-- Header -->
            <div v-if="channelStore.currentChannel"
                 class="flex gap-1 pt-1.5 pb-2.5 items-center select-none font-bold text-muted">

                <UIcon :name="getIcon(channelStore.currentChannel?.type || 0)" class="size-5" />
                <span>{{ channelStore.currentChannel?.name }}</span>
            </div>

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

            <UTextarea v-model="input" v-if="channelStore.currentChannel" autoresize :autofocus="true" :rows="2" :ui="{ base: 'min-h-[4rem] mt-4 rounded-lg max-h-128 focus:ring-0 border border-neutral-200 dark:border-neutral-800 bg-neutral-50 dark:bg-neutral-950' }" :placeholder="`Message #${channelStore.currentChannel?.name || '...'}`" variant="ghost" size="xl" />
        </div>
        <!-- Server Members -->
        <div v-if="channelStore.currentChannel" class="flex flex-col w-56 select-none">
            <div class="flex flex-1 max-w-56 max-h-9 mb-1 items-center gap-1">
                <UIcon name="i-lucide-users" class="size-5 text-muted" />
                <span class="font-bold text-muted">Members</span>
            </div>
            <div class="flex flex-1 items-center gap-2 p-4 border-b rounded-lg border border-neutral-200 dark:border-neutral-800 bg-white dark:bg-neutral-900 shadow-sm">
            </div>
        </div>
    </div>
</template>
