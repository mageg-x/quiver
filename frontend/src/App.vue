<template>
  <div id="app-body">
    <Toast />
    <ConfirmDialog />
    <router-view />
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'

const data = ref('')

onMounted(() => {
  console.log('Component mounted.')
})
</script>

<style scoped>

</style>