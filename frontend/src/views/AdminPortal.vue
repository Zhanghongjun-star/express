<template>
  <PortalShell
    title="运营端"
    subtitle="全网快件监控与明细导出"
    eyebrow="Operations center"
    login-path="/admin/login"
    :active="active"
    :menu="menu"
    @select="active = $event"
  >
    <section class="content-grid">
      <el-card v-for="item in stats" :key="item.label" shadow="never" class="stat-card">
        <span>{{ item.label }}</span>
        <b>{{ item.value }}</b>
      </el-card>
      <el-card shadow="never" class="span-2">
        <template #header>
          <div class="card-header">
            <strong>快件明细</strong>
            <div class="toolbar">
              <el-select v-model="status" size="small">
                <el-option label="全部状态" value="全部" />
                <el-option label="已揽收" value="已揽收" />
                <el-option label="运输中" value="运输中" />
                <el-option label="派送中" value="派送中" />
              </el-select>
              <el-input v-model="keyword" size="small" placeholder="顺丰单号" />
              <el-button type="primary" size="small" :icon="Download">导出明细</el-button>
            </div>
          </div>
        </template>
        <el-table :data="rows">
          <el-table-column type="index" label="序号" width="70" />
          <el-table-column prop="no" label="顺丰单号" />
          <el-table-column prop="time" label="寄件时间" />
          <el-table-column prop="status" label="当前状态" />
          <el-table-column prop="pay" label="付款类型" />
          <el-table-column prop="from" label="发件城市" />
          <el-table-column prop="to" label="收件城市" />
        </el-table>
      </el-card>
      <el-card shadow="never">
        <template #header><strong>业务实现逻辑</strong></template>
        <el-timeline>
          <el-timeline-item>订单状态由用户下单、快递员揽派、商家入库共同回写。</el-timeline-item>
          <el-timeline-item>数据总计按状态聚合：已揽收、运输中、派送中、已签收、转寄退回、作废。</el-timeline-item>
          <el-timeline-item>快件明细支持付款类型、状态和单号筛选，导出当前条件结果。</el-timeline-item>
        </el-timeline>
      </el-card>
    </section>
  </PortalShell>
</template>

<script setup lang="ts">
import { computed, ref } from "vue";
import { DataAnalysis, Download, List, Setting } from "@element-plus/icons-vue";
import PortalShell from "@/components/PortalShell.vue";
import { adminRows } from "@/data/mock";

const active = ref("stats");
const status = ref("全部");
const keyword = ref("");
const menu = [
  { key: "stats", label: "数据总计", icon: DataAnalysis },
  { key: "detail", label: "快件明细", icon: List },
  { key: "setting", label: "系统设置", icon: Setting },
];
const stats = [
  { label: "总件量", value: "24,918" },
  { label: "已揽收", value: "8,421" },
  { label: "运送中", value: "9,806" },
  { label: "派送中", value: "3,112" },
  { label: "已签收", value: "3,420" },
  { label: "转寄退回", value: "126" },
  { label: "作废", value: "33" },
];
const rows = computed(() => adminRows.filter((row) => {
  const matchedStatus = status.value === "全部" || row.status === status.value;
  const matchedKeyword = !keyword.value || row.no.includes(keyword.value);
  return matchedStatus && matchedKeyword;
}));
</script>
