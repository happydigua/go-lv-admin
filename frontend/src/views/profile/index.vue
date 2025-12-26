<template>
  <div class="profile-page">
    <n-grid :cols="24" :x-gap="16">
      <!-- 左侧用户信息卡片 -->
      <n-gi :span="8">
        <n-card>
          <div class="user-avatar-section">
            <n-avatar
              round
              :size="100"
              :src="userStore.userInfo?.avatar"
              fallback-src="https://via.placeholder.com/200"
            />
            <h2 class="user-name">{{ userStore.userInfo?.nickname || 'Admin' }}</h2>
            <p class="user-role">{{ userStore.userInfo?.Role?.name || '管理员' }}</p>
          </div>
          <n-divider />
          <n-descriptions :column="1" label-placement="left">
            <n-descriptions-item label="用户名">
              {{ userStore.userInfo?.username }}
            </n-descriptions-item>
            <n-descriptions-item label="邮箱">
              {{ userStore.userInfo?.email || '未设置' }}
            </n-descriptions-item>
            <n-descriptions-item label="手机号">
              {{ userStore.userInfo?.phone || '未设置' }}
            </n-descriptions-item>
            <n-descriptions-item label="注册时间">
              {{ formatDate(userStore.userInfo?.CreatedAt) }}
            </n-descriptions-item>
          </n-descriptions>
        </n-card>
      </n-gi>
      
      <!-- 右侧编辑区域 -->
      <n-gi :span="16">
        <n-card title="基本信息">
          <n-tabs type="line" animated>
            <n-tab-pane name="info" tab="个人资料">
              <n-form
                ref="infoFormRef"
                :model="infoForm"
                :rules="infoRules"
                label-placement="left"
                label-width="80"
              >
                <n-form-item label="昵称" path="nickname">
                  <n-input v-model:value="infoForm.nickname" placeholder="请输入昵称" />
                </n-form-item>
                <n-form-item label="邮箱" path="email">
                  <n-input v-model:value="infoForm.email" placeholder="请输入邮箱" />
                </n-form-item>
                <n-form-item label="手机号" path="phone">
                  <n-input v-model:value="infoForm.phone" placeholder="请输入手机号" />
                </n-form-item>
                <n-form-item>
                  <n-button type="primary" :loading="infoLoading" @click="handleUpdateInfo">保存修改</n-button>
                </n-form-item>
              </n-form>
            </n-tab-pane>
            
            <n-tab-pane name="password" tab="修改密码">
              <n-form
                ref="passwordFormRef"
                :model="passwordForm"
                :rules="passwordRules"
                label-placement="left"
                label-width="100"
              >
                <n-form-item label="原密码" path="oldPassword">
                  <n-input
                    v-model:value="passwordForm.oldPassword"
                    type="password"
                    placeholder="请输入原密码"
                    show-password-on="click"
                  />
                </n-form-item>
                <n-form-item label="新密码" path="newPassword">
                  <n-input
                    v-model:value="passwordForm.newPassword"
                    type="password"
                    placeholder="请输入新密码（至少6位）"
                    show-password-on="click"
                  />
                </n-form-item>
                <n-form-item label="确认新密码" path="confirmPassword">
                  <n-input
                    v-model:value="passwordForm.confirmPassword"
                    type="password"
                    placeholder="请再次输入新密码"
                    show-password-on="click"
                  />
                </n-form-item>
                <n-form-item>
                  <n-button type="primary" :loading="passwordLoading" @click="handleUpdatePassword">修改密码</n-button>
                </n-form-item>
              </n-form>
            </n-tab-pane>
          </n-tabs>
        </n-card>
      </n-gi>
    </n-grid>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue';
import { useUserStore } from '@/store/user';
import { useMessage } from 'naive-ui';
import { updateProfile, changePassword } from '@/api/profile';

const userStore = useUserStore();
const message = useMessage();

const infoFormRef = ref();
const passwordFormRef = ref();
const infoLoading = ref(false);
const passwordLoading = ref(false);

const infoForm = ref({
  nickname: '',
  email: '',
  phone: ''
});

const passwordForm = ref({
  oldPassword: '',
  newPassword: '',
  confirmPassword: ''
});

const infoRules = {
  nickname: { required: true, message: '请输入昵称', trigger: 'blur' },
  email: { type: 'email', message: '请输入正确的邮箱格式', trigger: 'blur' }
};

const passwordRules = {
  oldPassword: { required: true, message: '请输入原密码', trigger: 'blur' },
  newPassword: [
    { required: true, message: '请输入新密码', trigger: 'blur' },
    { min: 6, message: '密码长度不能小于6位', trigger: 'blur' }
  ],
  confirmPassword: [
    { required: true, message: '请再次输入新密码', trigger: 'blur' },
    {
      validator: (_rule: any, value: string) => value === passwordForm.value.newPassword,
      message: '两次密码输入不一致',
      trigger: 'blur'
    }
  ]
};

const formatDate = (dateStr: string) => {
  if (!dateStr) return '-';
  return new Date(dateStr).toLocaleDateString('zh-CN');
};

const handleUpdateInfo = () => {
  infoFormRef.value?.validate(async (errors: any) => {
    if (!errors) {
      infoLoading.value = true;
      try {
        await updateProfile(infoForm.value);
        message.success('个人信息更新成功');
        // 更新本地存储的用户信息
        if (userStore.userInfo) {
          userStore.userInfo.nickname = infoForm.value.nickname;
          userStore.userInfo.email = infoForm.value.email;
          userStore.userInfo.phone = infoForm.value.phone;
        }
      } catch (error) {
        message.error('更新失败');
      } finally {
        infoLoading.value = false;
      }
    }
  });
};

const handleUpdatePassword = () => {
  passwordFormRef.value?.validate(async (errors: any) => {
    if (!errors) {
      passwordLoading.value = true;
      try {
        await changePassword({
          oldPassword: passwordForm.value.oldPassword,
          newPassword: passwordForm.value.newPassword
        });
        message.success('密码修改成功');
        passwordForm.value = { oldPassword: '', newPassword: '', confirmPassword: '' };
      } catch (error) {
        message.error('密码修改失败');
      } finally {
        passwordLoading.value = false;
      }
    }
  });
};

onMounted(() => {
  if (userStore.userInfo) {
    infoForm.value = {
      nickname: userStore.userInfo.nickname || '',
      email: userStore.userInfo.email || '',
      phone: userStore.userInfo.phone || ''
    };
  }
});
</script>

<style scoped>
.profile-page {
  min-height: calc(100vh - 64px - 48px);
}

.user-avatar-section {
  text-align: center;
  padding: 20px 0;
}

.user-name {
  margin: 16px 0 8px;
  font-size: 20px;
}

.user-role {
  margin: 0;
  color: #999;
}
</style>
