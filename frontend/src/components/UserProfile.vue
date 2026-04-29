<script setup lang="ts">
import router from '@/router';
import type { LogoutResponse } from '@types';
import { apiFetch } from '@/utils/api';
import type { DropdownMenuItem } from '@nuxt/ui';
import { computed, ref } from 'vue'

const toast = useToast()
const maxUsernameLength = 15
const maxDescriptionLength = 18

const props = defineProps<{
    username?: string,
    description?: string,
    avatar_url?: string,
    status?: string
}>()

const usernameTruncated = computed(() => {
    if (!props.username) return '';
    return props.username.length > maxUsernameLength
        ? props.username.slice(0, maxUsernameLength) + '...'
        : props.username;
})

const descriptionTruncated = computed(() => {
    if (!props.description) return '';
    return props.description.length > maxDescriptionLength
        ? props.description.slice(0, maxDescriptionLength) + '...'
        : props.description;
})

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

const items = ref<DropdownMenuItem[][]>([
    [
        {
            label: 'Profile',
            icon: 'mdi-user'
        },
    ],
    [
        {
            label: 'Logout',
            icon: 'i-lucide-log-out',
            color: 'error',
            onClick: () => {
                onLogout();
            }
        }
    ]
])

const onLogout = async () => {
    const { error } = await apiFetch<LogoutResponse>('/revoke', {
        method: 'POST',
    })
    if (error) {
        toast.add({
            title: 'Logout failed',
            description: error.error || 'An unexpected error occurred. Please try again later.',
            color: 'error'
        })
        return
    }
    toast.add({
        title: 'Logged out',
        description: 'You have been logged out successfully.',
        color: 'success'
    })
    router.push({ name: 'login' })
}

</script>

<template>
    <div id="btn-profile" class="shrink-0 flex justify-between min-h-16 rounded-lg border border-neutral-200 dark:border-neutral-800 bg-white dark:bg-neutral-900 shadow-sm p-2">
        <UDropdownMenu :items="items" :content="{ align: 'start', side: 'top', sideOffset: 16 }" size="lg">
            <UButton variant="ghost" color="neutral" class="w-full rounded-lg text-left">
                <UUser class="w-full overflow-x-hidden min-w-0"
                       :name="usernameTruncated"
                       :description="description ? descriptionTruncated : statusMessage"
                       :avatar="{ src: avatar_url, loading: 'lazy' }"
                       :chip="{ color: chipColor, position: 'bottom-right' }"
                       :ui="{ name: 'truncate', description: 'truncate' }" />
            </UButton>
        </UDropdownMenu>

        <div class="h-full">
            <div class="flex gap-1 h-full p-1">
                <UButton id="btn-mute" variant="ghost" color="neutral" size="xl" class="rounded-lg">
                    <UIcon name="i-lucide-mic" size="20" />
                </UButton>
                <UButton id="btn-headset" variant="ghost" color="neutral" size="xl" class="rounded-lg">
                    <UIcon name="i-lucide-headset" size="20" />
                </UButton>
                <UButton id="btn-settings" variant="ghost" color="neutral" size="xl" class="rounded-lg">
                    <UIcon name="i-mdi-cog-outline" size="20" />
                </UButton>
            </div>
        </div>
    </div>
</template>
