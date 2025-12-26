<template>
  <div class="login-container">
    <!-- 语言切换器 -->
    <div class="locale-switch">
      <LocaleSwitcher />
    </div>

    <!-- 左侧品牌区域 -->
    <div class="login-left">
      <div class="brand-content">
        <h1 class="brand-title">Go Lv Vue Admin</h1>
        <p class="brand-desc">{{ t('login.brandDesc') }}</p>
        <div class="feature-list">
          <div class="feature-item">
            <n-icon size="24" color="#fff"><CheckmarkCircleOutline /></n-icon>
            <span>{{ t('login.feature1') }}</span>
          </div>
          <div class="feature-item">
            <n-icon size="24" color="#fff"><CheckmarkCircleOutline /></n-icon>
            <span>{{ t('login.feature2') }}</span>
          </div>
          <div class="feature-item">
            <n-icon size="24" color="#fff"><CheckmarkCircleOutline /></n-icon>
            <span>{{ t('login.feature3') }}</span>
          </div>
        </div>
      </div>
    </div>
    
    <!-- 右侧登录区域 -->
    <div class="login-right">
      <div class="login-box">
        <div class="login-header">
          <h2 class="title">{{ t('login.welcome') }}</h2>
          <p class="subtitle">{{ t('login.subtitle') }}</p>
        </div>
        <n-form ref="formRef" :model="formValue" :rules="rules" size="large">
          <n-form-item path="username">
            <n-input v-model:value="formValue.username" :placeholder="t('login.usernamePlaceholder')">
              <template #prefix>
                <n-icon :component="PersonOutline" />
              </template>
            </n-input>
          </n-form-item>
          <n-form-item path="password">
            <n-input
              v-model:value="formValue.password"
              type="password"
              show-password-on="mousedown"
              :placeholder="t('login.passwordPlaceholder')"
              @keyup.enter="handleLoginClick"
            >
              <template #prefix>
                <n-icon :component="LockClosedOutline" />
              </template>
            </n-input>
          </n-form-item>
          <n-form-item>
            <n-button type="primary" block :loading="loading" @click="handleLoginClick">
              {{ t('login.login') }}
            </n-button>
          </n-form-item>
        </n-form>
        <div class="login-footer">
          <span>© 2024 Go Lv Vue Admin</span>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue';
import { useRouter } from 'vue-router';
import { useI18n } from 'vue-i18n';
import { useUserStore } from '@/store/user';
import { type FormInst, useMessage } from 'naive-ui';
import { PersonOutline, LockClosedOutline, CheckmarkCircleOutline } from '@vicons/ionicons5';
import LocaleSwitcher from '@/components/LocaleSwitcher.vue';

const { t } = useI18n();
const router = useRouter();
const userStore = useUserStore();
const message = useMessage();

const formRef = ref<FormInst | null>(null);
const loading = ref(false);

const formValue = ref({
  username: 'admin',
  password: 'password'
});

const rules = computed(() => ({
  username: {
    required: true,
    message: t('login.usernameRequired'),
    trigger: 'blur'
  },
  password: {
    required: true,
    message: t('login.passwordRequired'),
    trigger: 'blur'
  }
}));

const handleLoginClick = (e: MouseEvent | KeyboardEvent) => {
  e.preventDefault();
  formRef.value?.validate(async (errors) => {
    if (!errors) {
      loading.value = true;
      const success = await userStore.handleLogin(formValue.value);
      loading.value = false;
      if (success) {
        message.success(t('login.loginSuccess'));
        router.push('/');
      }
    }
  });
};
</script>

<style scoped>
.login-container {
  display: flex;
  min-height: 100vh;
  width: 100%;
  position: relative;
}

.locale-switch {
  position: absolute;
  top: 20px;
  right: 20px;
  z-index: 100;
}

/* 左侧品牌区域 */
.login-left {
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: center;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  padding: 40px;
}

.brand-content {
  max-width: 480px;
  color: white;
}

.brand-title {
  font-size: 42px;
  font-weight: 700;
  margin: 0 0 16px 0;
  letter-spacing: 2px;
}

.brand-desc {
  font-size: 18px;
  opacity: 0.9;
  margin-bottom: 40px;
  line-height: 1.6;
}

.feature-list {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.feature-item {
  display: flex;
  align-items: center;
  gap: 12px;
  font-size: 16px;
}

/* 右侧登录区域 */
.login-right {
  width: 520px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: #f8f9fa;
  padding: 40px;
}

.login-box {
  width: 100%;
  max-width: 400px;
}

.login-header {
  text-align: center;
  margin-bottom: 40px;
}

.title {
  margin: 0 0 8px 0;
  font-size: 28px;
  font-weight: 600;
  color: #333;
}

.subtitle {
  margin: 0;
  font-size: 14px;
  color: #666;
}

.login-footer {
  text-align: center;
  margin-top: 24px;
  font-size: 12px;
  color: #999;
}

/* 响应式 - 平板及以下隐藏左侧 */
@media screen and (max-width: 900px) {
  .login-left {
    display: none;
  }
  
  .login-right {
    width: 100%;
    background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  }
  
  .login-box {
    background: white;
    padding: 40px;
    border-radius: 16px;
    box-shadow: 0 8px 32px rgba(0, 0, 0, 0.15);
  }
}

@media screen and (max-width: 480px) {
  .login-right {
    padding: 20px;
  }
  
  .login-box {
    padding: 24px;
  }
  
  .title {
    font-size: 24px;
  }
}
</style>
