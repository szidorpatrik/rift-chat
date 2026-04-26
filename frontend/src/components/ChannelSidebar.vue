<script setup lang="ts">
import type { Channel } from '@types'

defineProps<{
    channels: Channel[]
}>()
</script>

<template>
    <div class="overflow-y-auto w-64 shrink-0 flex flex-col rounded-2xl gap-4 border border-neutral-200 dark:border-neutral-800 bg-white dark:bg-neutral-900 shadow-sm p-4 [scrollbar-width:none] [&::-webkit-scrollbar]:hidden">

        <template v-if="channels.length">
            <div v-for="channel in channels" :key="channel.id" class="h-8 w-full shrink-0 flex items-center px-2 rounded hover:bg-neutral-100 dark:hover:bg-neutral-800 cursor-pointer">
                # {{ channel.name }}
            </div>
        </template>

        <template v-else>
            <USkeleton v-for="_ in 20" class="h-8 w-full shrink-0" />
        </template>

    </div>
</template>