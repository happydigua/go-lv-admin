<template>
  <n-upload
    :action="uploadUrl"
    :headers="headers"
    :accept="accept"
    :show-file-list="showFileList"
    :max="max"
    :default-file-list="defaultFileList"
    :list-type="listType"
    @finish="handleFinish"
    @remove="handleRemove"
    @before-upload="handleBeforeUpload"
  >
    <slot>
      <n-button v-if="listType !== 'image-card'">
        <template #icon>
          <n-icon :component="CloudUploadOutline" />
        </template>
        {{ buttonText }}
      </n-button>
      <n-upload-dragger v-else>
        <n-icon :component="ImageOutline" size="48" color="#999" />
        <p style="margin-top: 8px; color: #999;">点击或拖拽上传</p>
      </n-upload-dragger>
    </slot>
  </n-upload>
</template>

<script setup lang="ts">
import { computed, ref, watch } from 'vue';
import { CloudUploadOutline, ImageOutline } from '@vicons/ionicons5';
import type { UploadFileInfo } from 'naive-ui';

interface Props {
  type?: 'image' | 'file';
  value?: string;
  showFileList?: boolean;
  max?: number;
  listType?: 'text' | 'image' | 'image-card';
  buttonText?: string;
}

const props = withDefaults(defineProps<Props>(), {
  type: 'image',
  showFileList: true,
  max: 1,
  listType: 'text',
  buttonText: '上传文件'
});

const emit = defineEmits<{
  (e: 'update:value', value: string): void;
  (e: 'success', data: any): void;
  (e: 'remove', file: UploadFileInfo): void;
}>();

const token = localStorage.getItem('token') || '';

const uploadUrl = computed(() => {
  const baseUrl = import.meta.env.VITE_API_BASE_URL || '';
  return props.type === 'image' ? `${baseUrl}/upload/image` : `${baseUrl}/upload/file`;
});

const headers = computed(() => ({
  Authorization: `Bearer ${token}`
}));

const accept = computed(() => {
  if (props.type === 'image') {
    return 'image/png,image/jpeg,image/gif,image/webp,image/svg+xml';
  }
  return '*';
});

const defaultFileList = ref<UploadFileInfo[]>([]);

// 监听 value 变化，初始化文件列表
watch(
  () => props.value,
  (val) => {
    if (val && defaultFileList.value.length === 0) {
      defaultFileList.value = [{
        id: '1',
        name: val.split('/').pop() || 'file',
        status: 'finished',
        url: val
      }];
    }
  },
  { immediate: true }
);

const handleBeforeUpload = ({ file }: { file: UploadFileInfo }) => {
  // 图片大小限制 5MB，文件限制 20MB
  const maxSize = props.type === 'image' ? 5 * 1024 * 1024 : 20 * 1024 * 1024;
  if (file.file && file.file.size > maxSize) {
    return false;
  }
  return true;
};

const handleFinish = ({ file, event }: { file: UploadFileInfo; event?: ProgressEvent }) => {
  try {
    const response = JSON.parse((event?.target as XMLHttpRequest)?.response || '{}');
    if (response.code === 0 && response.data) {
      file.url = response.data.url;
      emit('update:value', response.data.url);
      emit('success', response.data);
    }
  } catch (e) {
    console.error('Upload response parse error:', e);
  }
  return file;
};

const handleRemove = ({ file }: { file: UploadFileInfo }) => {
  emit('update:value', '');
  emit('remove', file);
  return true;
};
</script>
