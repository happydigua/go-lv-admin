<template>
  <div class="table-demo">
    <ProTable
      ref="tableRef"
      title="用户列表 (Demo)"
      :columns="columns"
      :request="loadData"
      :search-schema="searchSchema"
      row-key="ID"
      @update:checked-row-keys="keys => checkedRowKeys = keys"
    >
      <template #toolbar>
        <n-button type="primary" @click="handleAdd">
          <template #icon><n-icon :component="AddOutline" /></template>
          新建数据
        </n-button>
        <n-button type="error" dashed @click="handleBatchDelete">
          批量删除
        </n-button>
      </template>
    </ProTable>

    <!-- 编辑/新增弹窗 -->
    <n-modal v-model:show="showModal" preset="card" :title="modalType === 'create' ? '新建数据' : '编辑数据'" style="width: 600px">
      <n-form
        ref="formRef"
        :model="formParams"
        :rules="rules"
        label-placement="left"
        label-width="80"
      >
        <n-form-item label="编号" path="code">
          <n-input v-model:value="formParams.code" placeholder="请输入编号" />
        </n-form-item>
        <n-form-item label="名称" path="name">
          <n-input v-model:value="formParams.name" placeholder="请输入名称" />
        </n-form-item>
        <n-form-item label="分类" path="category">
          <n-select
            v-model:value="formParams.category"
            :options="[
              { label: '电子产品', value: '电子产品' },
              { label: '家居用品', value: '家居用品' },
              { label: '办公文具', value: '办公文具' }
            ]"
            placeholder="请选择分类"
          />
        </n-form-item>
        <n-form-item label="金额" path="amount">
          <n-input-number v-model:value="formParams.amount" placeholder="请输入金额" style="width: 100%" />
        </n-form-item>
        <n-form-item label="状态" path="status">
          <n-select
            v-model:value="formParams.status"
            :options="[
              { label: '正常', value: 1 },
              { label: '停用', value: 2 }
            ]"
          />
        </n-form-item>
        <n-form-item label="描述" path="description">
          <n-input v-model:value="formParams.description" type="textarea" placeholder="请输入描述" />
        </n-form-item>
        <n-space justify="end">
          <n-button @click="showModal = false">取消</n-button>
          <n-button type="primary" @click="handleSubmit">确定</n-button>
        </n-space>
      </n-form>
    </n-modal>
  </div>
</template>

<script setup lang="ts">
import { ref, h, reactive } from 'vue';
import { 
  NTag, NButton, NSpace, useMessage, useDialog, 
  NModal, NCard, NForm, NFormItem, NInput, NSelect, NInputNumber,
  type FormInst, type FormRules
} from 'naive-ui';
import { AddOutline, PencilOutline, TrashOutline } from '@vicons/ionicons5';
import ProTable from '@/components/ProTable/index.vue';
import type { ProTableProps, SearchColumn } from '@/components/ProTable/types';
import { getDemoList, createDemo, updateDemo, deleteDemo } from '@/api/demo/table';

const message = useMessage();
const dialog = useDialog();
const tableRef = ref();

// 状态
const showModal = ref(false);
const modalType = ref<'create' | 'edit'>('create');
const formRef = ref<FormInst | null>(null);
const checkedRowKeys = ref<any[]>([]);

// 表单数据
const formParams = reactive({
  ID: 0,
  code: '',
  name: '',
  category: null as string | null,
  amount: 0,
  status: 1,
  description: ''
});

// 表单规则
const rules: FormRules = {
  code: { required: true, message: '请输入编号', trigger: 'blur' },
  name: { required: true, message: '请输入名称', trigger: 'blur' },
  category: { required: true, message: '请选择分类', trigger: ['blur', 'change'] },
  amount: { type: 'number', required: true, message: '请输入金额', trigger: ['blur', 'change'] }
};

// 加载数据
const loadData = async (params: any) => {
  const res: any = await getDemoList(params);
  return {
    list: res.list,
    total: res.total
  };
};

// 表格列定义
const columns = [
  { type: 'selection' },
  { title: 'ID', key: 'ID', width: 80, sorter: true },
  { title: '编号', key: 'code', width: 120 },
  { title: '名称', key: 'name', width: 150 },
  { title: '分类', key: 'category', width: 100 },
  { title: '金额', key: 'amount', width: 100, sorter: true },
  {
    title: '状态',
    key: 'status',
    width: 100,
    render(row: any) {
      return h(
        NTag,
        { type: row.status === 1 ? 'success' : 'error', bordered: false },
        { default: () => (row.status === 1 ? '正常' : '停用') }
      );
    }
  },
  { title: '创建时间', key: 'CreatedAt', width: 180 },
  {
    title: '操作',
    key: 'actions',
    width: 150,
    render(row: any) {
      return h(NSpace, null, {
        default: () => [
          h(
            NButton,
            {
              size: 'small',
              type: 'primary',
              secondary: true,
              onClick: () => handleEdit(row)
            },
            { icon: () => h(PencilOutline) }
          ),
          h(
            NButton,
            {
              size: 'small',
              type: 'error',
              secondary: true,
              onClick: () => handleDelete(row)
            },
            { icon: () => h(TrashOutline) }
          )
        ]
      });
    }
  }
];

// 搜索表单配置
const searchSchema: Record<string, SearchColumn> = {
  name: { type: 'input', label: '名称', placeholder: '请输入名称' },
  category: {
    type: 'select',
    label: '分类',
    options: [
      { label: '电子产品', value: '电子产品' },
      { label: '家居用品', value: '家居用品' },
      { label: '办公文具', value: '办公文具' }
    ]
  },
  status: {
    type: 'select',
    label: '状态',
    options: [
      { label: '正常', value: 1 },
      { label: '停用', value: 2 }
    ]
  }
};

const handleAdd = () => {
  modalType.value = 'create';
  formParams.ID = 0;
  formParams.code = `DEMO${Math.floor(Math.random() * 100000).toString().padStart(5, '0')}`; // 自动生成一个编号
  formParams.name = '';
  formParams.category = null;
  formParams.amount = 0;
  formParams.status = 1;
  formParams.description = '';
  showModal.value = true;
};

const handleEdit = (row: any) => {
  modalType.value = 'edit';
  Object.assign(formParams, row);
  showModal.value = true;
};

const handleSubmit = (e: MouseEvent) => {
  e.preventDefault();
  formRef.value?.validate(async (errors) => {
    if (!errors) {
      try {
        if (modalType.value === 'create') {
          await createDemo(formParams);
          message.success('创建成功');
        } else {
          await updateDemo(formParams);
          message.success('更新成功');
        }
        showModal.value = false;
        tableRef.value?.refresh();
      } catch (error) {
        // error handled by interceptor
      }
    }
  });
};

const handleDelete = (row: any) => {
  dialog.warning({
    title: '确认删除',
    content: `确定要删除 ${row.name} 吗？`,
    positiveText: '确定',
    negativeText: '取消',
    onPositiveClick: async () => {
      await deleteDemo(row.ID);
      message.success('删除成功');
      tableRef.value?.refresh();
    }
  });
};

const handleBatchDelete = () => {
  if (checkedRowKeys.value.length === 0) {
    message.warning('请至少选择一项');
    return;
  }
  dialog.warning({
    title: '批量删除',
    content: `确定要删除选中的 ${checkedRowKeys.value.length} 项吗？`,
    positiveText: '确定',
    negativeText: '取消',
    onPositiveClick: async () => {
      // 循环删除，实际项目建议添加后端批量接口
      for (const id of checkedRowKeys.value) {
        await deleteDemo(id);
      }
      message.success('批量删除成功');
      checkedRowKeys.value = [];
      tableRef.value?.refresh();
    }
  });
};
</script>
