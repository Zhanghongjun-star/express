<template>
  <PortalShell
    title="商家端"
    subtitle="驿站收发与财务台账"
    eyebrow="Station console"
    login-path="/merchant/login"
    :active="active"
    :menu="menu"
    @select="active = $event"
  >
    <section class="content-grid">
      <el-card v-for="job in merchantJobs" :key="job.title" shadow="never" class="job-card">
        <p class="muted">{{ job.note }}</p>
        <h2>{{ job.count }}</h2>
        <strong>{{ job.title }}</strong>
        <el-button text type="primary">进入处理</el-button>
      </el-card>
      <el-card shadow="never" class="span-2">
        <template #header><strong>快件收发管理</strong></template>
        <div class="flow-board">
          <div>
            <h3>寄件信息录入</h3>
            <p>录入寄收件人地址、手机号核验、物品类型和寄件渠道。</p>
          </div>
          <div>
            <h3>计费收款</h3>
            <p>称重计价、现结/到付/月结登记、优惠券核销。</p>
          </div>
          <div>
            <h3>派件入库</h3>
            <p>批量扫码入库，区分派送、自提、异常件和货架编号。</p>
          </div>
          <div>
            <h3>客户取件</h3>
            <p>手机号、取件码、运单号核验，扫码出库并更新状态。</p>
          </div>
        </div>
      </el-card>
      <el-card shadow="never" class="span-2">
        <template #header><strong>站点数据看板</strong></template>
        <div class="station-grid">
          <div><b>1,286</b><span>每日揽件总量</span></div>
          <div><b>94.6%</b><span>取件完成率</span></div>
          <div><b>73</b><span>滞留件数量</span></div>
          <div><b>8</b><span>售后异常工单</span></div>
        </div>
      </el-card>
    </section>
  </PortalShell>
</template>

<script setup lang="ts">
import { ref } from "vue";
import { ChatDotRound, DataBoard, Money, Tickets, Van } from "@element-plus/icons-vue";
import PortalShell from "@/components/PortalShell.vue";
import { merchantJobs } from "@/data/mock";

const active = ref("receive");
const menu = [
  { key: "receive", label: "快件收发", icon: Van },
  { key: "tasks", label: "订单任务", icon: Tickets },
  { key: "finance", label: "财务对账", icon: Money },
  { key: "message", label: "通知联系", icon: ChatDotRound },
  { key: "dashboard", label: "数据看板", icon: DataBoard },
];
</script>
