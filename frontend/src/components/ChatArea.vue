<script setup lang="ts">
import { ref } from 'vue'
import type { Message } from '@types'

defineProps<{
    messages: Message[]
}>()

const input = ref('')
</script>

<template>
    <div class="flex flex-1 flex-col gap-4 justify-between overflow-hidden rounded-lg">
        <UScrollArea class="hover-scrollbar h-full w-full border rounded-lg border-neutral-200 dark:border-neutral-800 bg-white dark:bg-neutral-900 shadow-sm" :ui="{ viewport: 'gap-4 p-4' }">
                <template v-if="messages.length">
                    <div v-for="msg in messages" :key="msg.id" class="flex flex-col gap-1">
                        <span class="font-bold text-sm">{{ msg.author }}</span>
                        <p class="text-neutral-600 dark:text-neutral-300">{{ msg.content }}</p>
                    </div>
                </template>
                <template v-else>
                    <USkeleton v-for="_ in 60" class="h-8 w-full shrink-0" />
                </template>
        </UScrollArea>
            
        <UTextarea v-model="input" 
                autoresize 
                :autofocus="true" 
                :rows="2" 
                :ui="{ base: 'min-h-[4rem] rounded-lg max-h-128 focus:ring-0 border border-neutral-200 dark:border-neutral-800 bg-neutral-50 dark:bg-neutral-950' }"
                placeholder="Message @placeholder..."
                variant="ghost" 
                size="xl"/>
    </div>
</template>
