<script setup lang="ts">
import { en } from '@nuxt/ui/locale'

</script>

<template>
    <UApp :locale="en">
        <router-view />
    </UApp>
</template>
