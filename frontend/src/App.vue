<template>
  <n-config-provider>
    <n-dialog-provider>
      <n-message-provider>
        <router-view></router-view>
      </n-message-provider>
    </n-dialog-provider>
  </n-config-provider>
</template>

<script setup lang="ts">
import { onMounted } from 'vue';
import { NConfigProvider, NMessageProvider, NDialogProvider } from 'naive-ui';
import { useSettingStore } from '@/store/setting';

const settingStore = useSettingStore();

onMounted(() => {
  settingStore.fetchSettings();
});
</script>

<style>
body {
  margin: 0;
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, 'Helvetica Neue', Arial, 'Noto Sans', sans-serif;
}
</style>
