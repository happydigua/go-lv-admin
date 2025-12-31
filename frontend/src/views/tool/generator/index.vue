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
          
          <n-divider>菜单配置</n-divider>
          
          <n-form-item label="父菜单">
            <n-tree-select
              v-model:value="config.parentMenuId"
              :options="menuOptions"
              placeholder="选择父菜单（可选）"
              clearable
              default-expand-all
            />
          </n-form-item>
          <n-form-item label="菜单图标">
            <n-input v-model:value="config.menuIcon" placeholder="如：DocumentOutline">
              <template #suffix>
                <n-tooltip trigger="hover">
                  <template #trigger>
                    <n-icon><information-circle-outline /></n-icon>
                  </template>
                  使用 ionicons5 图标名称
                </n-tooltip>
              </template>
            </n-input>
          </n-form-item>
          <n-form-item label="覆盖已存在">
            <n-switch v-model:value="config.overwrite" />
            <span style="margin-left: 8px; color: #999;">如果文件已存在，是否覆盖</span>
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
          <n-button type="success" @click="handleGenerate" :loading="generating">
            <template #icon><n-icon><save-outline /></n-icon></template>
            生成并写入文件
          </n-button>
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

      <n-tab-pane name="result" tab="生成结果" :disabled="!generateResult">
        <n-result
          :status="generateResult?.success ? 'success' : 'error'"
          :title="generateResult?.success ? '代码生成成功！' : '生成失败'"
          :description="generateResult?.message"
        >
          <template #footer v-if="generateResult?.success">
            <n-space vertical>
              <n-alert type="info" title="已创建文件">
                <n-ul>
                  <n-li v-for="file in generateResult?.files" :key="file">{{ file }}</n-li>
                </n-ul>
              </n-alert>
              <n-alert type="warning" title="后续步骤">
                <n-ol>
                  <n-li>重启后端服务以加载新路由</n-li>
                  <n-li>刷新前端页面查看新菜单</n-li>
                  <n-li>如需修改，可手动编辑生成的文件</n-li>
                </n-ol>
              </n-alert>
            </n-space>
          </template>
        </n-result>
      </n-tab-pane>
    </n-tabs>
  </n-card>
</template>

<script setup lang="ts">
import { h, ref, reactive, onMounted } from 'vue';
import { NButton, NSwitch, NSelect, NIcon, useMessage, useDialog } from 'naive-ui';
import { InformationCircleOutline, SaveOutline } from '@vicons/ionicons5';
import { getTables, getTableColumns, previewCode as previewCodeApi, generateCode } from '@/api/generator';
import { getMenuList } from '@/api/system/menu';
import hljs from 'highlight.js/lib/core';
import go from 'highlight.js/lib/languages/go';
import typescript from 'highlight.js/lib/languages/typescript';
import xml from 'highlight.js/lib/languages/xml';

hljs.registerLanguage('go', go);
hljs.registerLanguage('typescript', typescript);
hljs.registerLanguage('vue', xml);

const message = useMessage();
const dialog = useDialog();

const activeTab = ref('tables');
const loadingTables = ref(false);
const previewing = ref(false);
const generating = ref(false);
const tables = ref<any[]>([]);
const selectedTable = ref('');
const previewCode = ref<Record<string, string> | null>(null);
const generateResult = ref<any>(null);
const menuOptions = ref<any[]>([]);

const config = reactive({
  tableName: '',
  tableComment: '',
  moduleName: '',
  packageName: '',
  structName: '',
  parentMenuId: null as number | null,
  menuIcon: 'DocumentOutline',
  overwrite: false,
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

// 转换菜单为树形选项
function convertMenuToOptions(menus: any[]): any[] {
  return menus.map((menu: any) => ({
    key: menu.ID,
    label: menu.title,
    children: menu.children?.length ? convertMenuToOptions(menu.children) : undefined
  }));
}

const fetchMenus = async () => {
  try {
    const data: any = await getMenuList();
    menuOptions.value = convertMenuToOptions(data || []);
  } catch (error) {
    console.error('Failed to fetch menus:', error);
  }
};

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
    const data: any = await previewCodeApi({
      tableName: config.tableName,
      tableComment: config.tableComment,
      moduleName: config.moduleName,
      packageName: config.packageName,
      structName: config.structName,
      columns: config.columns
    });
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
  dialog.warning({
    title: '确认生成',
    content: '将生成代码文件并写入到项目目录，同时创建菜单记录。确定继续？',
    positiveText: '确定生成',
    negativeText: '取消',
    onPositiveClick: async () => {
      generating.value = true;
      try {
        const result: any = await generateCode({
          tableName: config.tableName,
          tableComment: config.tableComment,
          moduleName: config.moduleName,
          packageName: config.packageName,
          structName: config.structName,
          columns: config.columns,
          parentMenuId: config.parentMenuId || 0,
          menuIcon: config.menuIcon || 'DocumentOutline',
          overwrite: config.overwrite
        });
        generateResult.value = result;
        activeTab.value = 'result';
        message.success('代码生成成功！');
      } catch (error: any) {
        generateResult.value = {
          success: false,
          message: error.message || '生成失败',
          files: []
        };
        activeTab.value = 'result';
        message.error('生成失败');
      } finally {
        generating.value = false;
      }
    }
  });
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
  fetchMenus();
});
</script>

