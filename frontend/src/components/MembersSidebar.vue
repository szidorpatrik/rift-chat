<script setup lang="ts">
import { useServerStore } from '@stores/serverStore';
import { statusColorMap } from '@types';

const serverStore = useServerStore()
</script>

<template>
    <UScrollArea :items="serverStore.currentServer?.members || []"
                 virtualize
                 orientation="vertical"
                 class="select-none h-full w-full [scrollbar-width:none] [&::-webkit-scrollbar]:hidden">

        <template #default="{ item }">
            <UButton variant="ghost" color="neutral" class="w-full rounded-lg text-left">
                <UUser class="w-full overflow-x-hidden min-w-0"
                       :description="item.description"
                       :avatar="{ src: item.avatar_url, loading: 'lazy' }"
                       :chip="{ color: statusColorMap[item.status], position: 'bottom-right' }">
                    <template #name>
                        <div class="flex flex-1 min-w-0 max-w-38 items-center gap-1 overflow-x-hidden">
                            <template v-if="item.id === serverStore.currentServer?.owner_id">
                                <span class="truncate min-w-0 max-w-32">
                                    {{ item.username }}
                                </span>

                                <UIcon name="i-mdi-crown"
                                       class="shrink-0 text-warning" />
                            </template>
                            <template v-else>
                                <span class="truncate min-w-0 w-full">
                                    {{ item.username }}
                                </span>
                            </template>
                        </div>
                    </template>
                </UUser>
            </UButton>
        </template>

    </UScrollArea>
</template>