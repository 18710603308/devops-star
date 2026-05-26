<template>
  <div class="project-detail" v-if="project">
    <el-page-header @back="router.back()" style="margin-bottom: 20px">
      <template #content>
        <span class="text-large font-600 mr-3"> {{ project.display_name || project.name }} </span>
        <el-tag>{{ project.repo_type }}</el-tag>
      </template>
    </el-page-header>

    <el-row :gutter="20">
      <el-col :span="16">
        <el-card header="项目信息" class="info-card">
          <el-descriptions :column="2" border>
            <el-descriptions-item label="项目名称">{{ project.name }}</el-descriptions-item>
            <el-descriptions-item label="仓库类型">{{ project.repo_type }}</el-descriptions-item>
            <el-descriptions-item label="仓库地址">
              <el-link type="primary" :href="project.repo_url" target="_blank">{{ project.repo_url || '暂无' }}</el-link>
            </el-descriptions-item>
            <el-descriptions-item label="创建时间">{{ project.created_at }}</el-descriptions-item>
            <el-descriptions-item label="描述" :span="2">{{ project.description || '暂无描述' }}</el-descriptions-item>
          </el-descriptions>
        </el-card>

        <el-card header="流水线列表" style="margin-top: 20px">
          <el-table :data="pipelineList" stripe style="width: 100%">
            <el-table-column prop="name" label="流水线名称" />
            <el-table-column prop="status" label="状态" width="100">
              <template #default="{ row }">
                <el-tag :type="statusType(row.status)">{{ statusText(row.status) }}</el-tag>
              </template>
            </el-table-column>
            <el-table-column label="操作" width="160" fixed="right">
              <template #default="{ row }">
                <el-button link type="primary" size="small" @click="editPipeline(row)">编辑</el-button>
                <el-button link type="primary" size="small" @click="triggerPipeline(row)">触发</el-button>
              </template>
            </el-table-column>
          </el-table>
        </el-card>
      </el-col>

      <el-col :span="8">
        <el-card header="项目成员" class="member-card">
          <div style="margin-bottom: 12px">
            <el-button type="primary" size="small" @click="showAddMember = true">添加成员</el-button>
          </div>
          <el-table :data="memberList" stripe style="width: 100%">
            <el-table-column prop="username" label="用户名" />
            <el-table-column prop="role" label="角色" width="100" />
            <el-table-column label="操作" width="100">
              <template #default="{ row }">
                <el-button link type="danger" size="small">移除</el-button>
              </template>
            </el-table-column>
          </el-table>
        </el-card>

        <el-card header="快速操作" style="margin-top: 20px">
          <el-space direction="vertical" style="width: 100%">
            <el-button style="width: 100%" @click="router.push('/pipeline/editor/new?project_id=' + project.id)">新建流水线</el-button>
            <el-button style="width: 100%" @click="goToGitea">查看代码仓库</el-button>
            <el-button style="width: 100%" type="danger" @click="handleDeleteProject">删除项目</el-button>
          </el-space>
        </el-card>
      </el-col>
    </el-row>

    <el-dialog v-model="showAddMember" title="添加成员" width="420px">
      <el-form label-width="80px">
        <el-form-item label="用户名">
          <el-input v-model="newMember.username" placeholder="请输入用户名" />
        </el-form-item>
        <el-form-item label="角色">
          <el-select v-model="newMember.role" style="width: 100%">
            <el-option label="开发者" value="developer" />
            <el-option label="维护者" value="maintainer" />
            <el-option label="管理员" value="admin" />
          </el-select>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showAddMember = false">取消</el-button>
        <el-button type="primary" @click="addMember">确认</el-button>
      </template>
    </el-dialog>
  </div>
  <div v-else class="loading-container">
    <el-spinner size="large" />
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import { getProject, deleteProject as apiDelete, getProjectMembers, addProjectMember } from '@/api/project'

const route = useRoute()
const router = useRouter()

const project = ref<any>(null)
const memberList = ref<any[]>([])
const pipelineList = ref<any[]>([])
const showAddMember = ref(false)
const newMember = ref({ username: '', role: 'developer' })

const loadProject = async () => {
  try {
    const res = await getProject(route.params.id as string)
    project.value = res.data || res
    loadMembers()
    loadPipelines()
  } catch (err: any) {
    ElMessage.error(err.message || '加载失败')
    router.back()
  }
}

const loadMembers = async () => {
  try {
    const res = await getProjectMembers(route.params.id as string)
    memberList.value = res.data || []
  } catch {}
}

const loadPipelines = () => {
  // 模拟数据
  pipelineList.value = [
    { id: 1, name: '前端构建部署', status: 'success' },
    { id: 2, name: '后端镜像构建', status: 'running' },
  ]
}

const statusType = (s: string) => ({ success: 'success', running: 'primary', failed: 'danger', idle: 'info' }[s] || 'info')
const statusText = (s: string) => ({ success: '成功', running: '运行中', failed: '失败', idle: '空闲' }[s] || s)

const editPipeline = (row: any) => {
  router.push(`/pipeline/editor/${row.id}`)
}

const triggerPipeline = async (row: any) => {
  try {
    await import('@/api/pipeline').then(m => m.triggerPipeline(String(row.id)))
    ElMessage.success('流水线已触发')
    loadPipelines()
  } catch {}
}

const addMember = async () => {
  try {
    await addProjectMember(route.params.id as string, { user_id: 0, role: newMember.value.role })
    ElMessage.success('成员添加成功')
    showAddMember.value = false
    loadMembers()
  } catch {}
}

const goToGitea = () => {
  window.open(`http://localhost:3000/${project.value?.name}`, '_blank')
}

const handleDeleteProject = async () => {
  await ElMessageBox.confirm('确定删除该项目？此操作不可恢复！', '提示', { type: 'warning' })
  try {
    await apiDelete(route.params.id as string)
    ElMessage.success('删除成功')
    router.push('/projects')
  } catch {}
}

onMounted(() => {
  loadProject()
})
</script>

<style scoped>
.project-detail { padding: 0; }
.info-card { margin-bottom: 0; }
.member-card { margin-bottom: 0; }
.loading-container { display: flex; justify-content: center; align-items: center; height: 400px; }
</style>
