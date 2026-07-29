<template>
  <PortalShell
    title="快递员端"
    subtitle="揽派作业移动工作台"
    eyebrow="Courier operations"
    login-path="/courier/login"
    :active="active"
    :menu="menu"
    @select="active = $event"
  >
    <section class="content-grid">
      <el-card shadow="never" class="span-2">
        <template #header>
          <div class="card-header">
            <strong>待派件清单</strong>
            <div class="toolbar">
              <el-select v-model="filter" size="small">
                <el-option label="全部" value="全部" />
                <el-option label="货到付款" value="货到付款" />
                <el-option label="需签收" value="需签收" />
                <el-option label="易碎品" value="易碎品" />
              </el-select>
              <el-button size="small" :icon="Sort">排序</el-button>
            </div>
          </div>
        </template>
        <div class="task-list">
          <article v-for="task in tasks" :key="task.no" class="task-card">
            <div>
              <p class="muted">{{ task.no }} · {{ task.time }}</p>
              <h3>{{ task.name }} {{ task.phone }}</h3>
              <p>{{ task.address }}</p>
              <div class="tag-line">
                <el-tag v-for="tag in task.tags" :key="tag" size="small">{{ tag }}</el-tag>
              </div>
            </div>
            <div class="quick-actions">
              <el-button circle :icon="Phone" />
              <el-button circle :icon="Location" />
              <el-button type="primary">签收</el-button>
            </div>
          </article>
        </div>
      </el-card>
      <el-card shadow="never">
        <template #header><strong>服务数据</strong></template>
        <div class="metric-row vertical">
          <div><b>86</b><span>今日取派件</span></div>
          <div><b>97.3%</b><span>时效达标</span></div>
          <div><b>4.8</b><span>客户评分</span></div>
        </div>
      </el-card>
      <el-card shadow="never">
        <template #header><strong>业务实现逻辑</strong></template>
        <el-steps direction="vertical" :active="3">
          <el-step title="同步任务" description="按配送区域加载本人派件与揽收任务" />
          <el-step title="现场操作" description="电话、导航、投柜、签收和异常上报" />
          <el-step title="状态回写" description="完成后更新快件状态并同步考核数据" />
        </el-steps>
      </el-card>
    </section>
  </PortalShell>
</template>

<script setup lang="ts">
import { computed, ref } from "vue";
import { Bell, Box, Location, Phone, Sort, Tools } from "@element-plus/icons-vue";
import PortalShell from "@/components/PortalShell.vue";
import { courierTasks } from "@/data/mock";

const active = ref("delivery");
const filter = ref("全部");
const menu = [
  { key: "delivery", label: "派件", icon: Box },
  { key: "pickup", label: "取件", icon: Location },
  { key: "tools", label: "工具设备", icon: Tools },
  { key: "message", label: "设置与申诉", icon: Bell },
];
const tasks = computed(() => {
  if (filter.value === "全部") return courierTasks;
  return courierTasks.filter((task) => task.tags.includes(filter.value));
});
</script>
