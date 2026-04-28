<script setup lang="ts">
import { useServerStore } from '@/stores/serverStore';
import CreateServerView from '@/views/CreateServerView.vue';

const serverStore = useServerStore()

</script>

<template>
    <div class="overflow-y-auto w-12 shrink-0 flex flex-col items-center rounded-lg gap-4 p-0 [scrollbar-width:none] [&::-webkit-scrollbar]:hidden">

        <template v-if="serverStore.servers.length > 0">
            <div v-for="server in serverStore.servers" :key="server.id" class="w-full size-12 shrink-0 rounded-lg bg-neutral-200 dark:bg-neutral-800 flex items-center justify-center cursor-pointer hover:rounded-lg transition-all">
                {{ server.name.charAt(0) }}
            </div>
        </template>

        <UModal title="Create a server">
            <UButton class="p-2 justify-center" icon="i-lucide-plus" size="xl" color="primary" variant="outline" />

            <template #body>
                <CreateServerView />
            </template>
        </UModal>

    </div>
</template>
