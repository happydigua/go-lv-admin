<template>
  <n-card title="操作日志">
    <!-- 搜索区域 -->
    <n-space style="margin-bottom: 16px;">
      <n-input v-model:value="searchForm.username" placeholder="用户名" clearable style="width: 150px;" />
      <n-select
        v-model:value="searchForm.module"
        placeholder="选择模块"
        :options="moduleOptions"
        clearable
        style="width: 150px;"
      />
      <n-select
        v-model:value="searchForm.action"
        placeholder="选择操作"
        :options="actionOptions"
        clearable
        style="width: 120px;"
      />
      <n-button type="primary" @click="fetchData">
        <template #icon><n-icon :component="SearchOutline" /></template>
        搜索
      </n-button>
      <n-button @click="handleReset">重置</n-button>
      <n-popconfirm @positive-click="handleClear">
        <template #trigger>
          <n-button type="error">清空日志</n-button>
        </template>
        确定要清空所有操作日志吗？此操作不可恢复！
      </n-popconfirm>
    </n-space>

    <n-data-table
      :columns="columns"
      :data="tableData"
      :pagination="pagination"
      :loading="loading"
      :bordered="false"
      :row-key="(row: any) => row.ID"
      @update:page="handlePageChange"
      @update:page-size="handlePageSizeChange"
    />
  </n-card>

  <!-- 详情弹窗 -->
  <n-modal v-model:show="showDetailModal" preset="card" title="日志详情" style="width: 700px; max-height: 80vh;">
    <n-scrollbar style="max-height: calc(80vh - 100px);">
      <n-descriptions :column="2" label-placement="left" bordered>
        <n-descriptions-item label="用户">{{ currentLog?.username }}</n-descriptions-item>
        <n-descriptions-item label="IP">{{ currentLog?.ip }}</n-descriptions-item>
        <n-descriptions-item label="模块">{{ currentLog?.module }}</n-descriptions-item>
        <n-descriptions-item label="操作">{{ currentLog?.action }}</n-descriptions-item>
        <n-descriptions-item label="请求方式">{{ currentLog?.method }}</n-descriptions-item>
        <n-descriptions-item label="耗时">{{ currentLog?.latency }}ms</n-descriptions-item>
        <n-descriptions-item label="状态码">
          <n-tag :type="currentLog?.status === 200 ? 'success' : 'error'">{{ currentLog?.status }}</n-tag>
        </n-descriptions-item>
        <n-descriptions-item label="时间">{{ formatTime(currentLog?.CreatedAt) }}</n-descriptions-item>
        <n-descriptions-item label="请求路径" :span="2">{{ currentLog?.path }}</n-descriptions-item>
        <n-descriptions-item label="User-Agent" :span="2">
          <n-ellipsis style="max-width: 500px;">{{ currentLog?.userAgent }}</n-ellipsis>
        </n-descriptions-item>
      </n-descriptions>
      <n-divider>请求参数</n-divider>
      <n-code :code="currentLog?.body || '无'" language="json" style="max-height: 150px; overflow: auto;" />
      <n-divider>响应内容</n-divider>
      <div style="max-height: 200px; overflow: auto; border: 1px solid #e0e0e6; border-radius: 4px; padding: 8px;">
        <n-code :code="formatResponse(currentLog?.response)" language="json" word-wrap />
      </div>
    </n-scrollbar>
  </n-modal>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted, h } from 'vue';
import { NButton, NTag, useMessage, useDialog } from 'naive-ui';
import { SearchOutline } from '@vicons/ionicons5';
import { getOperationLogList, clearOperationLogs } from '@/api/system/log';

const message = useMessage();
const dialog = useDialog();

const loading = ref(false);
const tableData = ref<any[]>([]);
const showDetailModal = ref(false);
const currentLog = ref<any>(null);

const searchForm = reactive({
  username: '',
  module: null as string | null,
  action: null as string | null
});

const pagination = reactive({
  page: 1,
  pageSize: 10,
  itemCount: 0,
  showSizePicker: true,
  pageSizes: [10, 20, 50]
});

const moduleOptions = [
  { label: '登录', value: '登录' },
  { label: '仪表盘', value: '仪表盘' },
  { label: '用户管理', value: '用户管理' },
  { label: '角色管理', value: '角色管理' },
  { label: '菜单管理', value: '菜单管理' },
  { label: '个人中心', value: '个人中心' },
];

const actionOptions = [
  { label: '查询', value: '查询' },
  { label: '新增', value: '新增' },
  { label: '修改', value: '修改' },
  { label: '删除', value: '删除' },
  { label: '登录', value: '登录' },
];

const formatTime = (dateStr: string) => {
  if (!dateStr) return '-';
  return new Date(dateStr).toLocaleString('zh-CN');
};

const formatResponse = (response: string) => {
  if (!response) return '无';
  try {
    return JSON.stringify(JSON.parse(response), null, 2);
  } catch {
    return response;
  }
};

const columns = [
  { title: '用户', key: 'username', width: 100 },
  { title: 'IP', key: 'ip', width: 120 },
  { title: '模块', key: 'module', width: 100 },
  { title: '操作', key: 'action', width: 80 },
  { title: '方法', key: 'method', width: 80 },
  { title: '路径', key: 'path', ellipsis: { tooltip: true } },
  {
    title: '状态',
    key: 'status',
    width: 80,
    render: (row: any) => h(NTag, { type: row.status === 200 ? 'success' : 'error', size: 'small' }, { default: () => row.status })
  },
  { title: '耗时', key: 'latency', width: 80, render: (row: any) => `${row.latency}ms` },
  { title: '时间', key: 'CreatedAt', width: 160, render: (row: any) => formatTime(row.CreatedAt) },
  {
    title: '操作',
    key: 'actions',
    width: 80,
    render: (row: any) => h(NButton, { size: 'small', tertiary: true, type: 'info', onClick: () => handleDetail(row) }, { default: () => '详情' })
  }
];

const fetchData = async () => {
  loading.value = true;
  try {
    const res: any = await getOperationLogList({
      page: pagination.page,
      pageSize: pagination.pageSize,
      username: searchForm.username,
      module: searchForm.module || '',
      action: searchForm.action || ''
    });
    tableData.value = res.list || [];
    pagination.itemCount = res.total || 0;
  } catch (error) {
    console.error('Failed to fetch logs:', error);
  } finally {
    loading.value = false;
  }
};

const handlePageChange = (page: number) => {
  pagination.page = page;
  fetchData();
};

const handlePageSizeChange = (pageSize: number) => {
  pagination.pageSize = pageSize;
  pagination.page = 1;
  fetchData();
};

const handleReset = () => {
  searchForm.username = '';
  searchForm.module = null;
  searchForm.action = null;
  pagination.page = 1;
  fetchData();
};

const handleDetail = (row: any) => {
  currentLog.value = row;
  showDetailModal.value = true;
};

const handleClear = async () => {
  try {
    await clearOperationLogs();
    message.success('清空成功');
    fetchData();
  } catch (error) {
    message.error('清空失败');
  }
};

onMounted(() => {
  fetchData();
});
</script>
