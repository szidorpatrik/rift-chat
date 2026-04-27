<script setup lang="ts">
import { computed } from 'vue'

const props = defineProps<{
    username?: string,
    description?: string,
    avatar_url?: string,
    status?: string
}>()

const chipColor = computed(() => {
    switch (props.status) {
        case 'online':
            return 'primary';
        case 'idle':
            return 'warning';
        case 'dnd':
            return 'error';
        default:
            return 'neutral';
    }
})

const statusMessage = computed(() => {
    switch (props.status) {
        case 'online':
            return 'Online';
        case 'idle':
            return 'Idle';
        case 'dnd':
            return 'Do Not Disturb';
        default:
            return 'Offline';
    }
})

</script>

<template>
    <div :id="username" class="shrink-0 flex justify-between min-h-16 rounded-2xl border border-neutral-200 dark:border-neutral-800 bg-white dark:bg-neutral-900 shadow-sm p-2 items-center">
        <UButton variant="ghost" color="neutral" class="w-full rounded-xl">
            <UUser class="w-full"
                   :name="username"
                   :description="description ? description : statusMessage"
                   :avatar="{
                    src: avatar_url
                }"
                   :chip="{
                    color: chipColor,
                    position: 'bottom-right',
                }" />
        </UButton>

        <div id="settings" class="h-full">
            <div class="flex gap-1 h-full p-1">
                <UButton variant="ghost" color="neutral" size="xl" class="rounded-xl">
                    <UIcon name="mdi-microphone" size="20" />
                </UButton>
                <UButton variant="ghost" color="neutral" size="xl" class="rounded-xl">
                    <UIcon name="mdi-headset" size="20" />
                </UButton>
                <UButton variant="ghost" color="neutral" size="xl" class="rounded-xl">
                    <UIcon name="mdi-cog" size="20" />
                </UButton>
            </div>
        </div>
    </div>
</template>
