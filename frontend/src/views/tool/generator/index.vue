<template>
  <n-card title="代码生成器">
    <n-tabs v-model:value="activeTab">
      <n-tab-pane name="tables" tab="数据库表">
        <n-data-table
          :columns="tableColumns"
          :data="tables"
          :loading="loadingTables"
          :bordered="false"
          :row-key="(row: any) => row.tableName"
        />
      </n-tab-pane>

      <n-tab-pane name="config" tab="生成配置" :disabled="!selectedTable">
        <n-form label-placement="left" label-width="120" style="max-width: 600px; margin-bottom: 24px;">
          <n-form-item label="表名">
            <n-input :value="config.tableName" disabled />
          </n-form-item>
          <n-form-item label="表注释">
            <n-input v-model:value="config.tableComment" placeholder="如：文章" />
          </n-form-item>
          <n-form-item label="模块名">
            <n-input v-model:value="config.moduleName" placeholder="如：article (小写)" />
          </n-form-item>
          <n-form-item label="包名">
            <n-input v-model:value="config.packageName" placeholder="如：blog (小写)" />
          </n-form-item>
          <n-form-item label="结构体名">
            <n-input v-model:value="config.structName" placeholder="如：Article (首字母大写)" />
          </n-form-item>
        </n-form>

        <n-divider>字段配置</n-divider>
        <n-data-table
          :columns="columnConfigColumns"
          :data="config.columns"
          :bordered="false"
          size="small"
        />

        <n-space style="margin-top: 24px;">
          <n-button type="primary" @click="handlePreview" :loading="previewing">预览代码</n-button>
          <n-button type="success" @click="handleGenerate" :loading="generating">生成代码</n-button>
        </n-space>
      </n-tab-pane>

      <n-tab-pane name="preview" tab="代码预览" :disabled="!previewCode">
        <n-tabs type="line">
          <n-tab-pane v-for="(code, name) in previewCode" :key="name" :name="name" :tab="codeTabNames[name]">
            <n-card :bordered="false">
              <n-code :code="code" language="go" :hljs="hljs" />
            </n-card>
            <n-button style="margin-top: 12px;" @click="copyCode(code)">复制代码</n-button>
          </n-tab-pane>
        </n-tabs>
      </n-tab-pane>
    </n-tabs>
  </n-card>
</template>

<script setup lang="ts">
import { h, ref, reactive, onMounted } from 'vue';
import { NButton, NSwitch, NSelect, useMessage } from 'naive-ui';
import { getTables, getTableColumns, previewCode as previewCodeApi } from '@/api/generator';
import hljs from 'highlight.js/lib/core';
import go from 'highlight.js/lib/languages/go';
import typescript from 'highlight.js/lib/languages/typescript';
import xml from 'highlight.js/lib/languages/xml';

hljs.registerLanguage('go', go);
hljs.registerLanguage('typescript', typescript);
hljs.registerLanguage('vue', xml);

const message = useMessage();

const activeTab = ref('tables');
const loadingTables = ref(false);
const previewing = ref(false);
const generating = ref(false);
const tables = ref<any[]>([]);
const selectedTable = ref('');
const previewCode = ref<Record<string, string> | null>(null);

const config = reactive({
  tableName: '',
  tableComment: '',
  moduleName: '',
  packageName: '',
  structName: '',
  columns: [] as any[]
});

const codeTabNames: Record<string, string> = {
  model: 'Model',
  service: 'Service',
  api: 'API',
  router: 'Router',
  vue: 'Vue 页面',
  frontendApi: '前端 API'
};

const tableColumns = [
  { title: '表名', key: 'tableName' },
  { title: '注释', key: 'tableComment' },
  {
    title: '操作',
    key: 'actions',
    width: 120,
    render: (row: any) => h(NButton, { 
      size: 'small', 
      type: 'primary',
      onClick: () => handleSelectTable(row)
    }, { default: () => '配置生成' })
  }
];

const columnConfigColumns = [
  { title: '列名', key: 'columnName', width: 120 },
  { title: '类型', key: 'dataType', width: 100 },
  { title: 'Go字段', key: 'goField', width: 120 },
  { title: 'Go类型', key: 'goType', width: 100 },
  { title: '注释', key: 'columnComment', width: 120 },
  { 
    title: '列表', 
    key: 'isList', 
    width: 70,
    render: (row: any, index: number) => h(NSwitch, { 
      value: row.isList,
      onUpdateValue: (v: boolean) => { config.columns[index].isList = v; }
    })
  },
  { 
    title: '表单', 
    key: 'isForm', 
    width: 70,
    render: (row: any, index: number) => h(NSwitch, { 
      value: row.isForm,
      onUpdateValue: (v: boolean) => { config.columns[index].isForm = v; }
    })
  },
  { 
    title: '查询', 
    key: 'isQuery', 
    width: 70,
    render: (row: any, index: number) => h(NSwitch, { 
      value: row.isQuery,
      onUpdateValue: (v: boolean) => { config.columns[index].isQuery = v; }
    })
  },
  {
    title: '表单类型',
    key: 'formType',
    width: 120,
    render: (row: any, index: number) => h(NSelect, {
      value: row.formType,
      options: [
        { label: '输入框', value: 'input' },
        { label: '数字', value: 'number' },
        { label: '文本域', value: 'textarea' },
        { label: '日期', value: 'date' },
        { label: '下拉框', value: 'select' }
      ],
      size: 'small',
      onUpdateValue: (v: string) => { config.columns[index].formType = v; }
    })
  }
];

const fetchTables = async () => {
  loadingTables.value = true;
  try {
    const data: any = await getTables();
    tables.value = data || [];
  } catch (error) {
    console.error('Failed to fetch tables:', error);
  } finally {
    loadingTables.value = false;
  }
};

const handleSelectTable = async (row: any) => {
  selectedTable.value = row.tableName;
  config.tableName = row.tableName;
  config.tableComment = row.tableComment || toCamelCase(row.tableName, true);
  config.moduleName = toCamelCase(removePrefix(row.tableName), false);
  config.packageName = 'business';
  config.structName = toCamelCase(removePrefix(row.tableName), true);

  try {
    const columns: any = await getTableColumns(row.tableName);
    config.columns = columns || [];
  } catch (error) {
    console.error('Failed to fetch columns:', error);
  }

  activeTab.value = 'config';
};

const handlePreview = async () => {
  previewing.value = true;
  try {
    const data: any = await previewCodeApi(config);
    previewCode.value = data;
    activeTab.value = 'preview';
    message.success('代码预览成功');
  } catch (error) {
    message.error('预览失败');
  } finally {
    previewing.value = false;
  }
};

const handleGenerate = () => {
  message.info('代码已生成，请查看预览并复制到对应文件');
  handlePreview();
};

const copyCode = (code: string) => {
  navigator.clipboard.writeText(code);
  message.success('已复制到剪贴板');
};

// 辅助函数
function toCamelCase(s: string, upperFirst: boolean): string {
  const parts = s.split('_');
  return parts.map((p, i) => {
    if (i === 0 && !upperFirst) return p.toLowerCase();
    return p.charAt(0).toUpperCase() + p.slice(1).toLowerCase();
  }).join('');
}

function removePrefix(s: string): string {
  // 移除常见前缀如 lv_
  return s.replace(/^lv_/, '').replace(/^t_/, '').replace(/^sys_/, '');
}

onMounted(() => {
  fetchTables();
});
</script>
