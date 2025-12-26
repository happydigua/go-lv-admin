<template>
  <n-card title="系统设置">
    <n-spin :show="loading">
      <n-form
        ref="formRef"
        :model="formData"
        label-placement="left"
        label-width="100"
        style="max-width: 600px;"
      >
        <n-form-item label="系统名称" path="site_name">
          <n-input v-model:value="formData.site_name" placeholder="显示在标题栏和登录页" />
        </n-form-item>

        <n-form-item label="系统Logo" path="site_logo">
          <n-space vertical style="width: 100%;">
            <n-input v-model:value="formData.site_logo" placeholder="Logo图片URL" />
            <n-image v-if="formData.site_logo" :src="formData.site_logo" width="120" />
          </n-space>
        </n-form-item>

        <n-form-item label="底部版权" path="site_footer">
          <n-input v-model:value="formData.site_footer" placeholder="页面底部显示的版权信息" />
        </n-form-item>

        <n-form-item>
          <n-button type="primary" :loading="saving" @click="handleSave">
            保存设置
          </n-button>
        </n-form-item>
      </n-form>
    </n-spin>

    <n-divider />

    <n-alert type="info" title="说明">
      <ul style="margin: 0; padding-left: 20px;">
        <li>系统名称将显示在浏览器标题栏和登录页面</li>
        <li>存储配置（OSS/COS/R2）请在后端 <n-text code>config/config.yaml</n-text> 中配置</li>
        <li>修改后刷新页面生效</li>
      </ul>
    </n-alert>
  </n-card>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue';
import { useMessage } from 'naive-ui';
import { getSettings, updateSettings } from '@/api/system/setting';

const message = useMessage();
const loading = ref(false);
const saving = ref(false);
const formRef = ref();

const formData = ref({
  site_name: '',
  site_logo: '',
  site_footer: ''
});

const fetchSettings = async () => {
  loading.value = true;
  try {
    const data: any = await getSettings();
    formData.value = {
      site_name: data.site_name || '',
      site_logo: data.site_logo || '',
      site_footer: data.site_footer || ''
    };
  } catch (error) {
    console.error('Failed to fetch settings:', error);
  } finally {
    loading.value = false;
  }
};

const handleSave = async () => {
  saving.value = true;
  try {
    await updateSettings(formData.value);
    message.success('保存成功');
  } catch (error) {
    message.error('保存失败');
  } finally {
    saving.value = false;
  }
};

onMounted(() => {
  fetchSettings();
});
</script>
