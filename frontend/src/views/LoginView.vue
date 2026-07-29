<template>
  <div class="login-page">
    <section class="login-panel">
      <div class="identity-band">
        <div class="stamp">SF</div>
        <p>{{ app.copy }}</p>
        <h1>{{ app.title }}</h1>
        <div class="lane-map">
          <span>{{ app.steps[0] }}</span><i></i><span>{{ app.steps[1] }}</span><i></i><span>{{ app.steps[2] }}</span>
        </div>
      </div>
      <el-card class="auth-card" shadow="never">
        <div class="app-login-title">
          <strong>{{ app.name }}</strong>
        </div>
        <el-tabs v-model="mode" stretch @tab-change="clearFeedback">
          <el-tab-pane label="登录" name="login" />
          <el-tab-pane label="注册" name="register" />
        </el-tabs>
        <el-alert v-if="feedback" :title="feedback" :type="feedbackType" class="login-alert" show-icon :closable="false" />
        <el-form :model="form" label-position="top" @submit.prevent>
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
            {{ mode === "login" ? `登录${app.name}` : `注册${app.name}` }}
          </el-button>
          <el-button class="full-button ghost" @click="preview">本地预览进入{{ app.name }}</el-button>
        </el-form>
      </el-card>
    </section>
  </div>
</template>

<script setup lang="ts">
import { computed, reactive, ref } from "vue";
import { ElMessage } from "element-plus";
import { useRoute, useRouter } from "vue-router";
import { roleHome, useAuthStore } from "@/stores/auth";
import { userApi } from "@/api/user";
import type { Role } from "@/types";

type FeedbackType = "success" | "warning" | "error" | "info";

const router = useRouter();
const route = useRoute();
const auth = useAuthStore();
const mode = ref<"login" | "register">("login");
const loading = ref(false);
const codeLoading = ref(false);
const feedback = ref("");
const feedbackType = ref<FeedbackType>("info");
const form = reactive({ account: "", password: "", verifyCode: "" });

const role = computed(() => (route.meta.loginRole || "user") as Role);
const appMap: Record<Role, { name: string; title: string; copy: string; steps: string[] }> = {
  user: {
    name: "用户端 App",
    title: "顺丰用户寄件 App",
    copy: "寄快递、查快递、地址簿",
    steps: ["寄件", "查件", "我的"],
  },
  courier: {
    name: "快递员端 App",
    title: "顺丰快递员作业 App",
    copy: "派件、揽件、异常上报",
    steps: ["任务", "导航", "签收"],
  },
  merchant: {
    name: "商家端 App",
    title: "顺丰驿站商家 App",
    copy: "入库、取件、财务对账",
    steps: ["入库", "核验", "对账"],
  },
  admin: {
    name: "运营端 App",
    title: "顺丰运营管理 App",
    copy: "数据总计、明细查询、导出",
    steps: ["统计", "筛选", "导出"],
  },
};
const app = computed(() => appMap[role.value]);

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
    setFeedback("密码需要 8-32 位，和后端注册/登录规则保持一致。", "warning");
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
    return setFeedback(isEmail ? "邮箱格式不正确，请检查后再获取验证码。" : "手机号格式不正确，请输入 11 位手机号。", "warning");
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
    if (auth.session?.role !== role.value) {
      ElMessage.warning(`当前账号不是${app.value.name}角色，已进入账号对应应用。`);
    }
    setFeedback(`${mode.value === "login" ? "登录" : "注册"}成功，正在进入${app.value.name}。`, "success", false);
    router.push(String(route.query.redirect || roleHome[auth.session?.role || role.value]));
  } catch (error) {
    setFeedback(`${mode.value === "login" ? "登录" : "注册"}失败：${error instanceof Error ? error.message : "请检查账号、密码和后端服务。"}`);
  } finally {
    loading.value = false;
  }
}

function preview() {
  auth.switchRole(role.value);
  setFeedback(`已进入${app.value.name}本地预览。`, "success");
  router.push(roleHome[role.value]);
}
</script>
