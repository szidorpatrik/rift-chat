<script setup lang="ts">
import type { RadioGroupItem } from '@nuxt/ui';
import { useServerStore } from '@stores/serverStore';
import CreateServerView from '@views/CreateServerView.vue';
import JoinServerView from '@views/JoinServerView.vue';
import { computed, ref } from 'vue';

const openCreateServer = ref(false)
const serverStore = useServerStore()

const selectedServerId = computed({
    get: () => serverStore.currentServerId,
    set: (newId) => {
        if (newId) {
            serverStore.setCurrentServer(newId)
        }
    }
})

const servers = computed<RadioGroupItem[]>(() =>
    serverStore.servers.map(s => ({
        label: s.name[0].toUpperCase(),
        value: s.id,
        class: 'size-12 mb-2 cursor-pointer flex items-center justify-center font-semibold text-white',
    }))
)
</script>

<template>
    <div class="select-none overflow-y-auto w-12 shrink-0 flex flex-col items-center rounded-lg p-0 [scrollbar-width:none] [&::-webkit-scrollbar]:hidden">
        <template v-if="serverStore.servers.length > 0">
            <URadioGroup v-model="selectedServerId"
                         :items="servers"
                         size="md"
                         indicator="hidden"
                         variant="card" />
        </template>

        <UModal title="Create a server" v-model:open="openCreateServer">
            <UButton class="p-2 mt-2 mb-4 justify-center" icon="i-lucide-plus" size="xl" color="neutral" variant="outline" />

            <template #body>
                <CreateServerView v-model:open="openCreateServer" />
            </template>
        </UModal>

        <UModal title="Join a server" :scrollable="true" :ui="{ content: 'max-w-5xl' }">
            <UButton class="p-2 mb-4 justify-center" icon="i-lucide-compass" size="xl" color="neutral" variant="outline" />

            <template #body>
                <JoinServerView />
            </template>
        </UModal>
    </div>
</template>
