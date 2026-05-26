<template>
  <div class="login-container">
    <div class="login-card">
      <div class="login-header">
        <h1>🚀 DevOpsStar</h1>
        <p>国内开箱即用的 DevOps 平台</p>
      </div>
      <el-form
        ref="loginFormRef"
        :model="loginForm"
        :rules="loginRules"
        class="login-form"
        @keyup.enter="handleLogin"
      >
        <el-form-item prop="username">
          <el-input
            v-model="loginForm.username"
            placeholder="用户名 / 邮箱"
            size="large"
            :prefix-icon="User"
          />
        </el-form-item>
        <el-form-item prop="password">
          <el-input
            v-model="loginForm.password"
            type="password"
            placeholder="密码"
            size="large"
            :prefix-icon="Lock"
            show-password
          />
        </el-form-item>
        <el-form-item>
          <el-button
            type="primary"
            size="large"
            :loading="loading"
            class="login-btn"
            @click="handleLogin"
          >
            登 录
          </el-button>
        </el-form-item>
      </el-form>
      <div class="login-footer">
        <span>默认账号：<code>admin / admin123</code></span>
        <span>还没账号？<a @click="showRegister = true">立即注册</a></span>
      </div>
    </div>

    <!-- 注册对话框 -->
    <el-dialog v-model="showRegister" title="注册新账号" width="420px">
      <el-form :model="registerForm" :rules="registerRules" label-width="80px">
        <el-form-item label="用户名" prop="username">
          <el-input v-model="registerForm.username" />
        </el-form-item>
        <el-form-item label="邮箱" prop="email">
          <el-input v-model="registerForm.email" />
        </el-form-item>
        <el-form-item label="密码" prop="password">
          <el-input v-model="registerForm.password" type="password" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showRegister = false">取消</el-button>
        <el-button type="primary" :loading="regLoading" @click="handleRegister">注册</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive } from "vue";
import { useRouter } from "vue-router";
import { useUserStore } from "@/stores/user";
import { User, Lock } from "@element-plus/icons-vue";
import { ElMessage } from "element-plus";
import { register } from "@/api/auth";

const router = useRouter();
const userStore = useUserStore();
const loginFormRef = ref();
const loading = ref(false);
const regLoading = ref(false);
const showRegister = ref(false);

const loginForm = reactive({
  username: "admin",
  password: "admin123",
});

const registerForm = reactive({
  username: "",
  email: "",
  password: "",
});

const loginRules = {
  username: [{ required: true, message: "请输入用户名", trigger: "blur" }],
  password: [{ required: true, message: "请输入密码", trigger: "blur" }],
};

const registerRules = {
  username: [{ required: true, message: "请输入用户名", trigger: "blur" }],
  email: [{ required: true, type: "email", message: "请输入有效邮箱", trigger: "blur" }],
  password: [{ required: true, min: 6, message: "密码至少6位", trigger: "blur" }],
};

const handleLogin = async () => {
  await loginFormRef.value?.validate();
  loading.value = true;
  try {
    await userStore.handleLogin(loginForm.username, loginForm.password);
  } finally {
    loading.value = false;
  }
};

const handleRegister = async () => {
  regLoading.value = true;
  try {
    await register(registerForm);
    ElMessage.success("注册成功，请登录");
    showRegister.value = false;
  } catch (err: any) {
    ElMessage.error(err.message || "注册失败");
  } finally {
    regLoading.value = false;
  }
};
</script>

<style scoped>
.login-container {
  height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  background: linear-gradient(135deg, #1a1a2e 0%, #16213e 50%, #0f3460 100%);
}
.login-card {
  width: 420px;
  padding: 40px;
  background: #fff;
  border-radius: 12px;
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.2);
}
.login-header {
  text-align: center;
  margin-bottom: 32px;
}
.login-header h1 {
  font-size: 28px;
  color: #409eff;
  margin-bottom: 8px;
}
.login-header p {
  color: #909399;
  font-size: 14px;
}
.login-form {
  margin-top: 20px;
}
.login-btn {
  width: 100%;
}
.login-footer {
  margin-top: 16px;
  display: flex;
  justify-content: space-between;
  font-size: 13px;
  color: #909399;
}
.login-footer a {
  color: #409eff;
  cursor: pointer;
}
</style>
