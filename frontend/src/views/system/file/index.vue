<template>
  <n-card title="文件管理">
    <template #header-extra>
      <n-space>
        <n-tag :type="storageDriverType" size="small">
          {{ storageDriverLabel }}
        </n-tag>
      </n-space>
    </template>

    <!-- 存储配置说明 -->
    <n-alert type="info" style="margin-bottom: 16px;">
      <template #header>存储驱动配置</template>
      在后端 <n-text code>config/config.yaml</n-text> 中配置 <n-text code>storage.driver</n-text> 切换存储后端
    </n-alert>

    <!-- 支持的存储驱动 -->
    <n-grid :cols="4" :x-gap="16" :y-gap="16" style="margin-bottom: 24px;">
      <n-gi>
        <n-card size="small" :class="{ 'driver-card-active': currentDriver === 'local' }">
          <n-space vertical align="center">
            <n-icon :component="ServerOutline" size="32" color="#18a058" />
            <n-text strong>本地存储</n-text>
            <n-text depth="3" style="font-size: 12px;">免费 · 内置</n-text>
          </n-space>
        </n-card>
      </n-gi>
      <n-gi>
        <n-card size="small" :class="{ 'driver-card-active': currentDriver === 'oss' }">
          <n-space vertical align="center">
            <n-icon :component="CloudOutline" size="32" color="#ff6a00" />
            <n-text strong>阿里云 OSS</n-text>
            <n-text depth="3" style="font-size: 12px;">高可用 · 全球加速</n-text>
          </n-space>
        </n-card>
      </n-gi>
      <n-gi>
        <n-card size="small" :class="{ 'driver-card-active': currentDriver === 'cos' }">
          <n-space vertical align="center">
            <n-icon :component="CloudOutline" size="32" color="#006eff" />
            <n-text strong>腾讯云 COS</n-text>
            <n-text depth="3" style="font-size: 12px;">高可用 · CDN 加速</n-text>
          </n-space>
        </n-card>
      </n-gi>
      <n-gi>
        <n-card size="small" :class="{ 'driver-card-active': currentDriver === 'r2' }">
          <n-space vertical align="center">
            <n-icon :component="CloudOutline" size="32" color="#f48120" />
            <n-text strong>Cloudflare R2</n-text>
            <n-text depth="3" style="font-size: 12px;">零出口费用 · S3兼容</n-text>
          </n-space>
        </n-card>
      </n-gi>
    </n-grid>

    <n-divider />

    <!-- 上传区域 -->
    <n-grid :cols="2" :x-gap="24">
      <n-gi>
        <n-card title="图片上传" size="small">
          <n-upload
            :action="uploadImageUrl"
            :headers="headers"
            accept="image/*"
            list-type="image-card"
            :max="5"
            @finish="handleUploadFinish"
            @error="handleUploadError"
          >
            点击或拖拽上传图片
          </n-upload>
          <n-text depth="3" style="font-size: 12px; margin-top: 8px; display: block;">
            支持 jpg/png/gif/webp，单个文件最大 5MB
          </n-text>
        </n-card>
      </n-gi>
      <n-gi>
        <n-card title="文件上传" size="small">
          <n-upload
            :action="uploadFileUrl"
            :headers="headers"
            :max="5"
            @finish="handleUploadFinish"
            @error="handleUploadError"
          >
            <n-upload-dragger>
              <n-space vertical align="center">
                <n-icon :component="CloudUploadOutline" size="48" color="#999" />
                <n-text>点击或拖拽文件到此区域上传</n-text>
                <n-text depth="3" style="font-size: 12px;">
                  支持 pdf/doc/xls/ppt/txt/zip 等，最大 20MB
                </n-text>
              </n-space>
            </n-upload-dragger>
          </n-upload>
        </n-card>
      </n-gi>
    </n-grid>

    <n-divider />

    <!-- 已上传文件列表 -->
    <n-card title="已上传文件" size="small">
      <n-empty v-if="uploadedFiles.length === 0" description="暂无已上传文件" />
      <n-list v-else bordered>
        <n-list-item v-for="(file, index) in uploadedFiles" :key="index">
          <template #prefix>
            <n-avatar v-if="isImage(file.mimeType)" :src="file.url" size="small" />
            <n-icon v-else :component="DocumentOutline" size="24" />
          </template>
          <n-thing>
            <template #header>
              <n-ellipsis style="max-width: 300px;">{{ file.filename }}</n-ellipsis>
            </template>
            <template #description>
              <n-space>
                <n-text depth="3">{{ formatSize(file.size) }}</n-text>
                <n-text depth="3">{{ file.mimeType }}</n-text>
              </n-space>
            </template>
          </n-thing>
          <template #suffix>
            <n-space>
              <n-button size="small" tertiary @click="copyUrl(file.url)">
                <template #icon><n-icon :component="CopyOutline" /></template>
                复制链接
              </n-button>
              <n-button size="small" tertiary type="primary" tag="a" :href="file.url" target="_blank">
                <template #icon><n-icon :component="OpenOutline" /></template>
                预览
              </n-button>
            </n-space>
          </template>
        </n-list-item>
      </n-list>
    </n-card>
  </n-card>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue';
import { useMessage } from 'naive-ui';
import { 
  ServerOutline, 
  CloudOutline, 
  CloudUploadOutline, 
  DocumentOutline,
  CopyOutline,
  OpenOutline
} from '@vicons/ionicons5';

const message = useMessage();

// 当前存储驱动（从后端配置获取，这里暂时模拟）
const currentDriver = ref('local');

const storageDriverLabel = computed(() => {
  const labels: Record<string, string> = {
    local: '本地存储',
    oss: '阿里云 OSS',
    cos: '腾讯云 COS',
    r2: 'Cloudflare R2'
  };
  return labels[currentDriver.value] || '未知';
});

const storageDriverType = computed(() => {
  return currentDriver.value === 'local' ? 'success' : 'info';
});

const token = localStorage.getItem('token') || '';
const baseUrl = import.meta.env.VITE_API_BASE_URL || '';

const uploadImageUrl = computed(() => `${baseUrl}/upload/image`);
const uploadFileUrl = computed(() => `${baseUrl}/upload/file`);

const headers = computed(() => ({
  Authorization: `Bearer ${token}`
}));

interface UploadedFile {
  url: string;
  key: string;
  filename: string;
  size: number;
  mimeType: string;
}

const uploadedFiles = ref<UploadedFile[]>([]);

const handleUploadFinish = ({ file, event }: any) => {
  try {
    const response = JSON.parse((event?.target as XMLHttpRequest)?.response || '{}');
    if (response.code === 0 && response.data) {
      uploadedFiles.value.unshift(response.data);
      message.success(`${file.name} 上传成功`);
    } else {
      message.error(response.msg || '上传失败');
    }
  } catch (e) {
    message.error('上传响应解析失败');
  }
  return file;
};

const handleUploadError = ({ file }: any) => {
  message.error(`${file.name} 上传失败`);
};

const isImage = (mimeType: string) => {
  return mimeType?.startsWith('image/');
};

const formatSize = (bytes: number) => {
  if (bytes < 1024) return bytes + ' B';
  if (bytes < 1024 * 1024) return (bytes / 1024).toFixed(1) + ' KB';
  return (bytes / 1024 / 1024).toFixed(1) + ' MB';
};

const copyUrl = async (url: string) => {
  try {
    await navigator.clipboard.writeText(url);
    message.success('链接已复制');
  } catch {
    message.error('复制失败');
  }
};
</script>

<style scoped>
.driver-card-active {
  border-color: #18a058;
  background: rgba(24, 160, 88, 0.05);
}

:deep(.n-card) {
  transition: all 0.3s;
}

:deep(.n-card:hover) {
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.1);
}
</style>
