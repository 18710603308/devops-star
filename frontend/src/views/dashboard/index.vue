<template>
  <div class="dashboard">
    <!-- 统计卡片 -->
    <el-row :gutter="20" class="stat-cards">
      <el-col :span="6" v-for="item in statCards" :key="item.title">
        <el-card shadow="hover" class="stat-card">
          <div class="stat-card-body">
            <div class="stat-info">
              <div class="stat-title">{{ item.title }}</div>
              <div class="stat-value" :style="{ color: item.color }">{{ item.value }}</div>
              <div class="stat-trend">
                较昨日
                <span :style="{ color: item.trendUp ? '#67C23A' : '#F56C6C' }">
                  {{ item.trend }}
                </span>
              </div>
            </div>
            <el-icon class="stat-icon" :style="{ color: item.color }">
              <component :is="item.icon" />
            </el-icon>
          </div>
        </el-card>
      </el-col>
    </el-row>

    <!-- 图表区域 -->
    <el-row :gutter="20" class="chart-row">
      <el-col :span="16">
        <el-card header="构建趋势（近 7 天）" class="chart-card">
          <div ref="buildChartRef" style="height: 320px"></div>
        </el-card>
      </el-col>
      <el-col :span="8">
        <el-card header="构建结果分布" class="chart-card">
          <div ref="pieChartRef" style="height: 320px"></div>
        </el-card>
      </el-col>
    </el-row>

    <!-- 最近构建记录 -->
    <el-card header="最近构建记录" class="table-card">
      <el-table :data="buildList" stripe style="width: 100%">
        <el-table-column prop="id" label="构建 ID" width="120" />
        <el-table-column prop="project" label="项目" />
        <el-table-column prop="pipeline" label="流水线" />
        <el-table-column prop="status" label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="statusType(row.status)">{{ statusText(row.status) }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="branch" label="分支" width="120" />
        <el-table-column prop="trigger" label="触发者" width="120" />
        <el-table-column prop="duration" label="耗时" width="100" />
        <el-table-column prop="time" label="时间" width="180" />
        <el-table-column label="操作" width="100" fixed="right">
          <template #default="{ row }">
            <el-button link type="primary" size="small" @click="viewLog(row)">日志</el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, nextTick } from "vue";
import { getPipelineStats, getBuildList } from "@/api/monitor";
import { Pipeline, Share, Box, Promotion } from "@element-plus/icons-vue";
import * as echarts from "echarts";
import type { ECharts } from "echarts";

const statCards = ref([
  { title: "今日构建", value: 15, trend: "+3", trendUp: true, color: "#409EFF", icon: Pipeline },
  { title: "项目总数", value: 8, trend: "+1", trendUp: true, color: "#67C23A", icon: Share },
  { title: "镜像总数", value: 42, trend: "+5", trendUp: true, color: "#E6A23C", icon: Box },
  { title: "部署成功", value: 12, trend: "-1", trendUp: false, color: "#F56C6C", icon: Promotion },
]);

const buildList = ref([
  { id: "build-001", project: "devops-star", pipeline: "前端构建", status: "success", branch: "main", trigger: "admin", duration: "2m 30s", time: "2026-05-27 00:30:00" },
  { id: "build-002", project: "devops-star", pipeline: "后端构建", status: "running", branch: "main", trigger: "admin", duration: "1m 15s", time: "2026-05-27 00:35:00" },
  { id: "build-003", project: "devops-star", pipeline: "镜像构建", status: "failed", branch: "main", trigger: "admin", duration: "45s", time: "2026-05-27 00:20:00" },
  { id: "build-004", project: "devops-star", pipeline: "部署生产", status: "success", branch: "main", trigger: "admin", duration: "3m 10s", time: "2026-05-26 22:00:00" },
]);

const buildChartRef = ref<HTMLElement>();
const pieChartRef = ref<HTMLElement>();
let buildChart: ECharts | null = null;
let pieChart: ECharts | null = null;

const statusType = (status: string) => {
  const map: Record<string, string> = { success: "success", running: "primary", failed: "danger", pending: "warning" };
  return map[status] || "info";
};

const statusText = (status: string) => {
  const map: Record<string, string> = { success: "成功", running: "构建中", failed: "失败", pending: "等待中" };
  return map[status] || status;
};

const viewLog = (row: any) => {
  // 打开日志抽屉
  console.log("查看日志", row.id);
};

const initCharts = () => {
  nextTick(() => {
    // 构建趋势图
    if (buildChartRef.value) {
      buildChart = echarts.init(buildChartRef.value);
      buildChart.setOption({
        tooltip: { trigger: "axis" },
        legend: { data: ["成功", "失败"] },
        xAxis: { type: "category", data: ["05-21", "05-22", "05-23", "05-24", "05-25", "05-26", "05-27"] },
        yAxis: { type: "value" },
        series: [
          { name: "成功", type: "bar", data: [12, 15, 10, 18, 14, 16, 12], itemStyle: { color: "#67C23A" } },
          { name: "失败", type: "bar", data: [2, 1, 3, 1, 2, 1, 3], itemStyle: { color: "#F56C6C" } },
        ],
      });
    }

    // 饼图
    if (pieChartRef.value) {
      pieChart = echarts.init(pieChartRef.value);
      pieChart.setOption({
        tooltip: { trigger: "item" },
        legend: { bottom: 0 },
        series: [
          {
            type: "pie",
            radius: ["40%", "70%"],
            data: [
              { value: 98, name: "成功", itemStyle: { color: "#67C23A" } },
              { value: 12, name: "失败", itemStyle: { color: "#F56C6C" } },
              { value: 5, name: "进行中", itemStyle: { color: "#409EFF" } },
            ],
          },
        ],
      });
    }
  });
};

onMounted(() => {
  initCharts();
  window.addEventListener("resize", () => {
    buildChart?.resize();
    pieChart?.resize();
  });
});
</script>

<style scoped>
.dashboard {
  min-height: 100%;
}
.stat-cards {
  margin-bottom: 20px;
}
.stat-card {
  border-radius: 8px;
}
.stat-card-body {
  display: flex;
  align-items: center;
  justify-content: space-between;
}
.stat-info {
  flex: 1;
}
.stat-title {
  font-size: 13px;
  color: #909399;
  margin-bottom: 8px;
}
.stat-value {
  font-size: 28px;
  font-weight: bold;
  margin-bottom: 4px;
}
.stat-trend {
  font-size: 12px;
  color: #909399;
}
.stat-icon {
  font-size: 48px;
  opacity: 0.15;
}
.chart-row {
  margin-bottom: 20px;
}
.chart-card {
  border-radius: 8px;
}
.table-card {
  border-radius: 8px;
}
</style>
