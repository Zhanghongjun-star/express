<template>
  <div class="login-page">
    <div v-if="showSplash" class="sf-splash-screen">
      <img src="/sf-splash.png" alt="顺丰速运启动页" />
    </div>

    <main v-else class="login-stage">
      <section class="login-visual">
        <div class="login-visual-frame">
          <img src="/sf-splash.png" alt="顺丰速运视觉图" />
        </div>
      </section>

      <section class="login-panel">
        <div class="login-brand">
          <span>{{ app.name }}</span>
          <strong>{{ app.title }}</strong>
          <p>{{ app.copy }}</p>
        </div>

        <el-card class="auth-card" shadow="never">
          <el-tabs v-model="mode" stretch @tab-change="clearFeedback">
            <el-tab-pane label="登录" name="login" />
            <el-tab-pane label="注册" name="register" />
          </el-tabs>

          <el-alert
            v-if="feedback"
            :title="feedback"
            :type="feedbackType"
            class="login-alert"
            show-icon
            :closable="false"
          />

          <el-form class="auth-form" :model="form" label-position="top" @submit.prevent>
            <el-form-item label="账号">
              <el-input v-model.trim="form.account" placeholder="手机号或邮箱" @keyup.enter="submit" />
            </el-form-item>
            <el-form-item v-if="mode === 'register'" label="验证码">
              <div class="inline-field">
                <el-input v-model.trim="form.verifyCode" placeholder="6 位验证码" @keyup.enter="submit" />
                <el-button :loading="codeLoading" @click="sendCode">获取</el-button>
              </div>
            </el-form-item>
            <el-form-item label="密码">
              <el-input v-model="form.password" placeholder="8-32 位密码" show-password @keyup.enter="submit" />
            </el-form-item>
            <el-button type="primary" class="full-button" :loading="loading" @click="submit">
              {{ mode === "login" ? "立即登录" : "立即注册" }}
            </el-button>
            <el-button class="full-button ghost" @click="preview">本地预览</el-button>
          </el-form>
        </el-card>
      </section>
    </main>
  </div>
</template>

<script setup lang="ts">
import { onMounted, reactive, ref } from "vue";
import { ElMessage } from "element-plus";
import { useRoute, useRouter } from "vue-router";
import { roleHome, useAuthStore } from "@/stores/auth";
import { userApi } from "@/api/user";

type FeedbackType = "success" | "warning" | "error" | "info";

const router = useRouter();
const route = useRoute();
const auth = useAuthStore();
const mode = ref<"login" | "register">("login");
const loading = ref(false);
const codeLoading = ref(false);
const showSplash = ref(true);
const feedback = ref("");
const feedbackType = ref<FeedbackType>("info");
const form = reactive({ account: "", password: "", verifyCode: "" });

const app = {
  name: "顺丰速运",
  title: "欢迎回来",
  hero: "寄件、查件、地址簿、会员福利，一个界面全搞定",
  badge: "SF EXPRESS",
  subtitle: "像手机应用一样清爽的顺丰服务入口",
  copy: "登录后进入首页，已完成接口的功能会直接连接后端与数据库。",
};

onMounted(() => {
  window.setTimeout(() => {
    showSplash.value = false;
  }, 2600);
});

function accountType() {
  return form.account.includes("@") ? "email" : "phone";
}

function showFeedbackMessage(message: string, type: FeedbackType) {
  if (type === "success") {
    ElMessage.success(message);
  } else if (type === "warning") {
    ElMessage.warning(message);
  } else if (type === "error") {
    ElMessage.error(message);
  }
}

function setFeedback(message: string, type: FeedbackType = "error", showToast = true) {
  feedback.value = message;
  feedbackType.value = type;
  if (showToast) showFeedbackMessage(message, type);
}

function clearFeedback() {
  feedback.value = "";
}

function validateBase() {
  clearFeedback();
  if (!form.account) {
    setFeedback("请输入手机号或邮箱。", "warning");
    return false;
  }
  const isEmail = form.account.includes("@");
  const validAccount = isEmail ? /^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(form.account) : /^1\d{10}$/.test(form.account);
  if (!validAccount) {
    setFeedback(isEmail ? "邮箱格式不正确，请检查后再试。" : "手机号格式不正确，请输入 11 位手机号。", "warning");
    return false;
  }
  if (mode.value === "register" && !/^\d{6}$/.test(form.verifyCode)) {
    setFeedback("请输入 6 位验证码。", "warning");
    return false;
  }
  if (form.password.length < 8 || form.password.length > 32) {
    setFeedback("密码需要 8-32 位。", "warning");
    return false;
  }
  return true;
}

async function sendCode() {
  clearFeedback();
  if (!form.account) return setFeedback("请先输入手机号或邮箱。", "warning");
  const isEmail = form.account.includes("@");
  const validAccount = isEmail ? /^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(form.account) : /^1\d{10}$/.test(form.account);
  if (!validAccount) {
    return setFeedback(isEmail ? "邮箱格式不正确，请先检查。" : "手机号格式不正确，请输入 11 位手机号。", "warning");
  }
  codeLoading.value = true;
  try {
    await userApi.sendCode(form.account, accountType());
    setFeedback("验证码已发送，请查看短信或邮箱。", "success");
  } catch (error) {
    setFeedback(`验证码发送失败：${error instanceof Error ? error.message : "请检查后端服务是否启动。"}`);
  } finally {
    codeLoading.value = false;
  }
}

async function submit() {
  if (!validateBase()) return;
  loading.value = true;
  try {
    clearFeedback();
    if (mode.value === "login") {
      await auth.login(form.account, form.password);
    } else {
      await auth.register({
        phone: accountType() === "phone" ? form.account : "",
        email: accountType() === "email" ? form.account : "",
        password: form.password,
        verifyCode: form.verifyCode,
      });
    }
    setFeedback(`${mode.value === "login" ? "登录" : "注册"}成功，正在进入页面。`, "success", false);
    router.push(String(route.query.redirect || roleHome.user));
  } catch (error) {
    const message = error instanceof Error ? error.message : "请检查账号、密码和后端服务。";
    if (mode.value === "login" && message.includes("还没有注册")) {
      mode.value = "register";
      form.verifyCode = "";
      setFeedback(message, "warning");
      return;
    }
    setFeedback(`${mode.value === "login" ? "登录" : "注册"}失败：${message}`);
  } finally {
    loading.value = false;
  }
}

function preview() {
  auth.switchRole("user");
  setFeedback("已进入本地预览。", "success");
  router.push(roleHome.user);
}
</script>
