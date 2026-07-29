<template>
  <div class="portal-shell">
    <aside class="sidebar">
      <div class="brand-mark">
        <div class="brand-block">SF</div>
        <div>
          <strong>顺丰速运</strong>
          <span>{{ subtitle }}</span>
        </div>
      </div>
      <el-menu :default-active="active" class="side-menu" @select="$emit('select', $event)">
        <el-menu-item v-for="item in menu" :key="item.key" :index="item.key">
          <el-icon><component :is="item.icon" /></el-icon>
          <span>{{ item.label }}</span>
        </el-menu-item>
      </el-menu>
    </aside>
    <main class="workbench">
      <header class="topbar">
        <div>
          <p class="eyebrow">{{ eyebrow }}</p>
          <h1>{{ title }}</h1>
        </div>
        <div class="top-actions">
          <el-tag effect="plain">{{ today }}</el-tag>
          <el-button :icon="SwitchButton" @click="logout">退出</el-button>
        </div>
      </header>
      <slot />
    </main>
  </div>
</template>

<script setup lang="ts">
import { computed } from "vue";
import type { Component } from "vue";
import { SwitchButton } from "@element-plus/icons-vue";
import { useAuthStore } from "@/stores/auth";

const props = defineProps<{
  title: string;
  subtitle: string;
  eyebrow: string;
  loginPath: string;
  active: string;
  menu: Array<{ key: string; label: string; icon: Component }>;
}>();

defineEmits<{ select: [key: string] }>();

const auth = useAuthStore();
const today = computed(() => new Intl.DateTimeFormat("zh-CN", { dateStyle: "medium" }).format(new Date()));

async function logout() {
  await auth.logout();
  window.location.hash = `#${props.loginPath}`;
}
</script>
