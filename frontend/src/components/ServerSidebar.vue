<script setup lang="ts">
import type { Server } from '@types'

defineProps<{
    servers: Server[]
}>()

</script>

<template>
    <div class="overflow-y-auto w-18 shrink-0 flex flex-col items-center rounded-lg gap-4 p-0 [scrollbar-width:none] [&::-webkit-scrollbar]:hidden">

        <template v-if="servers.length">
            <div v-for="server in servers" :key="server.id" class="w-full size-16 shrink-0 rounded-lg bg-neutral-200 dark:bg-neutral-800 flex items-center justify-center cursor-pointer hover:rounded-lg transition-all">
                {{ server.name.charAt(0) }}
            </div>
        </template>

        <template v-else>
            <USkeleton v-for="_ in 20" class="w-full size-16 shrink-0 rounded-lg" />
        </template>

    </div>
</template>
