<template>
  <div class="monitor-page">
    <!-- 统计卡片 -->
    <el-row :gutter="20" class="stat-cards">
      <el-col :span="6" v-for="item in stats" :key="item.label">
        <el-card shadow="hover" class="stat-card">
          <div class="stat-value" :style="{ color: item.color }">{{ item.value }}</div>
          <div class="stat-label">{{ item.label }}</div>
        </el-card>
      </el-col>
    </el-row>

    <!-- Gitea 服务状态 + 系统资源 -->
    <el-row :gutter="20" style="margin-top: 20px">
      <el-col :span="12">
        <el-card header="Gitea 服务状态">
          <el-descriptions :column="3" border>
            <el-descriptions-item label="服务状态">
              <el-tag :type="giteaStatus.status === 'running' ? 'success' : 'danger'">
                {{ giteaStatus.status === 'running' ? '运行中' : '异常' }}
              </el-tag>
            </el-descriptions-item>
            <el-descriptions-item label="仓库数量">{{ giteaStatus.repos }}</el-descriptions-item>
            <el-descriptions-item label="活跃用户">{{ giteaStatus.users }}</el-descriptions-item>
            <el-descriptions-item label="今日推送">{{ giteaStatus.pushes }} 次</el-descriptions-item>
            <el-descriptions-item label="存储占用">{{ giteaStatus.storage }}</el-descriptions-item>
            <el-descriptions-item label="运行时间">{{ giteaStatus.uptime }}</el-descriptions-item>
          </el-descriptions>
        </el-card>
      </el-col>

      <el-col :span="12">
        <el-card header="系统资源（实时）">
          <div style="padding: 20px 0">
            <p style="margin-bottom: 8px">CPU 使用率 <span style="float:right">{{ systemResource.cpu }}%</span></p>
            <el-progress :percentage="systemResource.cpu" :color="'#409EFF'" />
            <p style="margin: 16px 0 8px">内存使用率 <span style="float:right">{{ systemResource.memory }}%</span></p>
            <el-progress :percentage="systemResource.memory" :color="'#67C23A'" />
            <p style="margin: 16px 0 8px">磁盘使用率 <span style="float:right">{{ systemResource.disk }}%</span></p>
            <el-progress :percentage="systemResource.disk" :color="'#E6A23C'" />
          </div>
        </el-card>
      </el-col>
    </el-row>

    <!-- 部署统计 -->
    <el-row :gutter="20" style="margin-top: 20px">
      <el-col :span="24">
        <el-card header="今日部署统计">
          <el-descriptions :column="4" border>
            <el-descriptions-item label="今日部署">{{ deployStats.total_today }}</el-descriptions-item>
            <el-descriptions-item label="成功">
              <span style="color: #67C23A">{{ deployStats.success_today }}</span>
            </el-descriptions-item>
            <el-descriptions-item label="失败">
              <span style="color: #F56C6C">{{ deployStats.failed_today }}</span>
            </el-descriptions-item>
            <el-descriptions-item label="最近部署">{{ deployStats.last_deploy }}</el-descriptions-item>
          </el-descriptions>
        </el-card>
      </el-col>
    </el-row>

    <!-- 最近告警 -->
    <el-row :gutter="20" style="margin-top: 20px">
      <el-col :span="24">
        <el-card header="最近告警">
          <el-timeline>
            <el-timeline-item
              v-for="(alert, index) in recentAlerts"
              :key="index"
              :timestamp="alert.time"
              :type="alert.type"
              :placement="'top'"
            >
              <p v-html="alert.content"></p>
            </el-timeline-item>
          </el-timeline>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { monitorAPI } from '@/api/monitor'
import { formatDistanceToNow } from 'date-fns'
import { zhCN } from 'date-fns/locale'

// 统计卡片数据
const stats = ref([
  { label: '今日构建', value: 0, color: '#409EFF' },
  { label: '成功', value: 0, color: '#67C23A' },
  { label: '失败', value: 0, color: '#F56C6C' },
  { label: '部署次数', value: 0, color: '#E6A23C' },
])

// Gitea 服务状态
const giteaStatus = ref({
  status: 'running',
  repos: 0,
  users: 0,
  pushes: 0,
  storage: '0 GB',
  uptime: '0 天',
})

// 系统资源
const systemResource = ref({
  cpu: 0,
  memory: 0,
  disk: 0,
})

// 部署统计
const deployStats = ref({
  total_today: 0,
  success_today: 0,
  failed_today: 0,
  last_deploy: '-',
})

// 最近告警（模拟数据，实际应从 API 获取）
const recentAlerts = ref([
  { time: '2026-05-27 01:35:00', type: 'success', content: '流水线 <strong>前端构建</strong> 执行成功' },
  { time: '2026-05-27 01:20:00', type: 'danger', content: '流水线 <strong>镜像构建</strong> 执行失败' },
  { time: '2026-05-26 22:00:00', type: 'success', content: '部署 <strong>生产环境</strong> 成功' },
])

// 获取监控数据
const fetchMonitorData = async () => {
  try {
    // 获取流水线统计
    const pipelineRes = await monitorAPI.getPipelineStats()
    if (pipelineRes.data) {
      const data = pipelineRes.data
      stats.value[0].value = data.total_runs || 0
      stats.value[1].value = data.success_runs || 0
      stats.value[2].value = data.failed_runs || 0
    }

    // 获取部署统计
    const deployRes = await monitorAPI.getDeployStats()
    if (deployRes.data) {
      const data = deployRes.data
      deployStats.value.total_today = data.total_today || 0
      deployStats.value.success_today = data.success_today || 0
      deployStats.value.failed_today = data.failed_today || 0
      deployStats.value.last_deploy = data.last_deploy
        ? formatDistanceToNow(new Date(data.last_deploy), { addSuffix: true, locale: zhCN })
        : '-'
      // 更新统计卡片中的部署次数
      stats.value[3].value = data.total || 0
    }

    // 获取系统资源（从 Prometheus 获取）
    const statsRes = await monitorAPI.getStats()
    if (statsRes.data) {
      const data = statsRes.data
      // 系统资源
      systemResource.value.cpu = Math.round(data.cpu?.usage || 0)
      systemResource.value.memory = Math.round(data.memory?.usage || 0)
      systemResource.value.disk = Math.round(data.disk?.usage || 0)

      // Gitea 状态
      giteaStatus.value.status = data.gitea_status || 'unknown'
      giteaStatus.value.repos = data.gitea_repos || 0
      giteaStatus.value.users = data.gitea_users || 0
      giteaStatus.value.pushes = data.gitea_pushes || 0
    }
  } catch (error) {
    console.error('获取监控数据失败：', error)
  }
}

onMounted(() => {
  fetchMonitorData()
  // 每 30 秒刷新一次
  const timer = setInterval(fetchMonitorData, 30000)
  return () => clearInterval(timer)
})
</script>

<style scoped>
.monitor-page { padding: 0; }
.stat-cards { margin-bottom: 0; }
.stat-card { border-radius: 8px; text-align: center; }
.stat-value { font-size: 32px; font-weight: bold; }
.stat-label { font-size: 13px; color: #909399; margin-top: 4px; }
</style>
