<template>
  <n-dropdown :options="localeOptions" @select="handleSelect">
    <n-button quaternary circle>
      <template #icon>
        <n-icon :component="LanguageOutline" />
      </template>
    </n-button>
  </n-dropdown>
</template>

<script setup lang="ts">
import { computed } from 'vue';
import { useI18n } from 'vue-i18n';
import { LanguageOutline } from '@vicons/ionicons5';
import { setLocale } from '@/locales';

const { locale } = useI18n();

const localeOptions = computed(() => [
  {
    label: '🇨🇳 简体中文',
    key: 'zh-CN',
    disabled: locale.value === 'zh-CN'
  },
  {
    label: '🇺🇸 English',
    key: 'en-US',
    disabled: locale.value === 'en-US'
  }
]);

const handleSelect = (key: string) => {
  setLocale(key);
};
</script>
