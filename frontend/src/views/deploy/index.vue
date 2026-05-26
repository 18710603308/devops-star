<template>
  <div class="deploy-page">
    <el-card class="page-header">
      <div class="header-actions">
        <h2>部署管理</h2>
        <el-button type="primary" @click="showEnv = true">
          <el-icon><Plus /></el-icon>
          新建环境
        </el-button>
      </div>
    </el-card>

    <el-row :gutter="20">
      <el-col :span="8" v-for="env in envList" :key="env.id">
        <el-card shadow="hover" class="env-card">
          <div class="env-header">
            <el-tag :type="env.status === 'running' ? 'success' : 'info'">{{ env.status === 'running' ? '运行中' : '空闲' }}</el-tag>
            <span class="env-name">{{ env.display_name || env.name }}</span>
          </div>
          <div class="env-body">
            <p>类型：{{ env.deploy_type === 'docker' ? 'Docker' : 'Kubernetes' }}</p>
            <p>最近部署：{{ env.last_deploy || '暂无' }}</p>
          </div>
          <div class="env-footer">
            <el-button link type="primary" size="small" @click="handleDeploy(env)">部署</el-button>
            <el-button link type="danger" size="small">删除</el-button>
          </div>
        </el-card>
      </el-col>
    </el-row>

    <el-card header="部署记录" class="table-card">
      <el-table :data="deployList" stripe style="width: 100%">
        <el-table-column prop="id" label="部署 ID" width="120" />
        <el-table-column prop="env_name" label="环境" width="120" />
        <el-table-column prop="image_tag" label="镜像版本" />
        <el-table-column prop="status" label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="row.status === 'success' ? 'success' : row.status === 'failed' ? 'danger' : 'warning'">
              {{ { success: '成功', failed: '失败', running: '部署中' }[row.status] || row.status }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="deployed_by" label="部署者" width="120" />
        <el-table-column prop="created_at" label="时间" width="180" />
      </el-table>
    </el-card>

    <el-dialog v-model="showEnv" title="新建部署环境" width="520px">
      <el-form :model="envForm" label-width="100px">
        <el-form-item label="环境名称">
          <el-input v-model="envForm.name" />
        </el-form-item>
        <el-form-item label="显示名称">
          <el-input v-model="envForm.display_name" />
        </el-form-item>
        <el-form-item label="部署类型">
          <el-select v-model="envForm.deploy_type" style="width: 100%">
            <el-option label="Docker" value="docker" />
            <el-option label="Kubernetes" value="k8s" />
          </el-select>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showEnv = false">取消</el-button>
        <el-button type="primary" @click="saveEnv">保存</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { Plus } from '@element-plus/icons-vue'
import { getEnvironments, createEnvironment, getDeployHistory } from '@/api/pipeline'

const envList = ref([
  { id: 1, name: 'dev', display_name: '开发环境', deploy_type: 'docker', status: 'running', last_deploy: '2026-05-27 00:30:00' },
  { id: 2, name: 'test', display_name: '测试环境', deploy_type: 'docker', status: 'idle', last_deploy: '2026-05-26 18:00:00' },
  { id: 3, name: 'prod', display_name: '生产环境', deploy_type: 'k8s', status: 'running', last_deploy: '2026-05-25 03:00:00' },
])
const deployList = ref([
  { id: 1, env_name: '开发环境', image_tag: 'v1.0.0', status: 'success', deployed_by: 'admin', created_at: '2026-05-27 00:30:00' },
  { id: 2, env_name: '测试环境', image_tag: 'v0.9.0', status: 'success', deployed_by: 'admin', created_at: '2026-05-26 18:00:00' },
  { id: 3, env_name: '生产环境', image_tag: 'v0.8.0', status: 'failed', deployed_by: 'admin', created_at: '2026-05-25 03:00:00' },
])
const showEnv = ref(false)
const envForm = ref({ name: '', display_name: '', deploy_type: 'docker' })

const loadData = async () => {
  try {
    const [envRes, deployRes] = await Promise.all([getEnvironments(), getDeployHistory()])
    envList.value = envRes.data || envList.value
    deployList.value = deployRes.data || deployList.value
  } catch {}
}

const saveEnv = async () => {
  try {
    await createEnvironment(envForm.value)
    ElMessage.success('环境创建成功')
    showEnv.value = false
    loadData()
  } catch {}
}

const handleDeploy = (env: any) => {
  ElMessage.info(`开始部署到 ${env.display_name}...`)
}

onMounted(loadData)
</script>

<style scoped>
.deploy-page { display: flex; flex-direction: column; gap: 20px; }
.page-header { padding: 20px; }
.header-actions { display: flex; justify-content: space-between; align-items: center; }
.header-actions h2 { margin: 0; }
.env-card { border-radius: 8px; }
.env-header { display: flex; align-items: center; gap: 8px; margin-bottom: 12px; }
.env-name { font-size: 16px; font-weight: 600; }
.env-body p { color: #606266; margin: 4px 0; font-size: 13px; }
.env-footer { margin-top: 12px; display: flex; gap: 8px; }
.table-card { border-radius: 8px; }
</style>
