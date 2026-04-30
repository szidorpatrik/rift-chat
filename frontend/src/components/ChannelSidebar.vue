<script setup lang="ts">
import { useServerStore } from '@/stores/serverStore';
import type { Channel } from '@types'

defineProps<{
    channels: Channel[]
}>()

const serverStore = useServerStore()

</script>

<template>
    <div class="overflow-y-auto w-64 shrink-0 flex flex-col rounded-lg gap-4 border border-neutral-200 dark:border-neutral-800 bg-white dark:bg-neutral-900 shadow-sm p-4 [scrollbar-width:none] [&::-webkit-scrollbar]:hidden">
        <div class="flex flex-0 w-full">
            <UButton class="p-2 max-w-fit truncate text-md font-semibold text-neutral-400 tracking-wide"
                     color="neutral"
                     variant="ghost"
                     size="md"
                     trailing-icon="i-lucide-chevron-down"
                     :label="serverStore.currentServer?.name || 'Select a server'" />
            <UButton class="p-2 max-w-fit ml-auto text-neutral-400"
                     color="neutral"
                     variant="ghost"
                     size="xl"
                     icon="i-lucide-user-plus" />
        </div>

        <template v-if="channels.length">
            <div v-for="channel in channels" :key="channel.id" class="h-8 w-full shrink-0 flex items-center px-2 rounded-lg hover:bg-neutral-100 dark:hover:bg-neutral-800 cursor-pointer">
                # {{ channel.name }}
            </div>
        </template>

        <template v-else>
            <USkeleton v-for="_ in 8" class="h-8 w-full shrink-0" />
        </template>

    </div>
</template>