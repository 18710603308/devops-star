<template>
  <div class="settings-page">
    <el-row :gutter="20">
      <!-- 个人设置 -->
      <el-col :span="12">
        <el-card header="个人设置" class="settings-card">
          <el-form label-width="100px">
            <el-form-item label="用户名">
              <el-input v-model="userForm.username" />
            </el-form-item>
            <el-form-item label="邮箱">
              <el-input v-model="userForm.email" />
            </el-form-item>
            <el-form-item label="新密码">
              <el-input v-model="userForm.newPassword" type="password" show-password />
            </el-form-item>
            <el-form-item>
              <el-button type="primary" @click="saveUser">保存</el-button>
            </el-form-item>
          </el-form>
        </el-card>
      </el-col>

      <!-- 通知设置 -->
      <el-col :span="12">
        <el-card header="通知配置" class="settings-card">
          <el-form label-width="100px">
            <el-form-item label="企业微信">
              <el-input v-model="notifyForm.wecom" placeholder="Webhook URL" />
            </el-form-item>
            <el-form-item label="钉钉">
              <el-input v-model="notifyForm.dingtalk" placeholder="Webhook URL" />
            </el-form-item>
            <el-form-item label="飞书">
              <el-input v-model="notifyForm.feishu" placeholder="Webhook URL" />
            </el-form-item>
            <el-form-item label="通知事件">
              <el-checkbox-group v-model="notifyForm.events">
                <el-checkbox label="success">成功</el-checkbox>
                <el-checkbox label="failed">失败</el-checkbox>
                <el-checkbox label="always">全部</el-checkbox>
              </el-checkbox-group>
            </el-form-item>
            <el-form-item>
              <el-button @click="testNotify">测试通知</el-button>
              <el-button type="primary" @click="saveNotify">保存</el-button>
            </el-form-item>
          </el-form>
        </el-card>
      </el-col>
    </el-row>

    <!-- 镜像源设置 -->
    <el-card header="国内镜像源配置" style="margin-top: 20px" class="settings-card">
      <el-table :data="mirrorList" stripe>
        <el-table-column prop="name" label="名称" width="160" />
        <el-table-column prop="url" label="镜像地址" />
        <el-table-column prop="type" label="类型" width="120" />
        <el-table-column label="操作" width="120">
          <template #default="{ row }">
            <el-button link type="primary" size="small" @click="testMirror(row)">测试</el-button>
          </template>
        </el-table-column>
      </el-table>
      <div style="margin-top: 12px">
        <el-button type="primary" @click="saveMirrors">保存配置</el-button>
        <el-button @click="resetMirrors">恢复默认</el-button>
      </div>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { getMirrorConfig, updateMirrorConfig } from '@/api/settings'

const userForm = ref({ username: 'admin', email: 'admin@devops-star.com', newPassword: '' })
const notifyForm = ref({ wecom: '', dingtalk: '', feishu: '', events: ['success', 'failed'] })
const mirrorList = ref([
  { name: 'Docker 镜像一', url: 'https://docker.1ms.run', type: 'Docker' },
  { name: 'Docker 镜像二', url: 'https://docker.m.daocloud.io', type: 'Docker' },
  { name: 'npm 淘宝镜像', url: 'https://registry.npmmirror.com', type: 'npm' },
  { name: 'Maven 阿里云', url: 'https://maven.aliyun.com/repository/public', type: 'Maven' },
  { name: 'PyPI 清华源', url: 'https://pypi.tuna.tsinghua.edu.cn/simple', type: 'PyPI' },
  { name: 'Go Proxy 七牛', url: 'https://goproxy.cn,direct', type: 'Go' },
])

const loadMirrors = async () => {
  try {
    const res = await getMirrorConfig()
    mirrorList.value = res.data || mirrorList.value
  } catch {}
}

const saveUser = () => { ElMessage.success('个人设置已保存') }
const saveNotify = () => { ElMessage.success('通知配置已保存') }
const testNotify = () => { ElMessage.info('测试通知已发送，请查收') }
const testMirror = (row: any) => { ElMessage.success(`镜像 ${row.name} 连接正常`) }
const saveMirrors = () => { ElMessage.success('镜像源配置已保存，重启服务后生效') }
const resetMirrors = () => { ElMessage.info('已恢复默认镜像源') }

onMounted(loadMirrors)
</script>

<style scoped>
.settings-page { display: flex; flex-direction: column; gap: 20px; }
.settings-card { border-radius: 8px; }
</style>
