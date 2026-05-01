<script setup lang="ts">
import { useServerStore } from '@stores/serverStore';
import CreateServerView from '@views/CreateServerView.vue';
import JoinServerView from '@views/JoinServerView.vue';
const serverStore = useServerStore()

</script>

<template>
    <div class="overflow-y-auto w-12 shrink-0 flex flex-col items-center rounded-lg gap-4 p-0 [scrollbar-width:none] [&::-webkit-scrollbar]:hidden">

        <template v-if="serverStore.servers.length > 0">
            <div v-for="server in serverStore.servers" :key="server.id"
                 @click="serverStore.setCurrentServer(server.id)"
                 class="select-none w-full size-12 shrink-0 rounded-lg bg-neutral-200 dark:bg-neutral-800 flex items-center justify-center cursor-pointer hover:rounded-lg transition-all">
                {{ server.name.charAt(0) }}
            </div>
        </template>

        <UModal title="Create a server">
            <UButton class="p-2 justify-center" icon="i-lucide-plus" size="xl" color="neutral" variant="outline" />

            <template #body>
                <CreateServerView />
            </template>
        </UModal>

        <UModal title="Join a server" :scrollable="true" :ui="{ content: 'max-w-5xl' }">
            <UButton class="p-2 justify-center" icon="i-lucide-compass" size="xl" color="neutral" variant="outline" />

            <template #body>
                <JoinServerView />
            </template>
        </UModal>
    </div>
</template>
