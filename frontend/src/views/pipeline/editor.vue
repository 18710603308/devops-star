<template>
  <div class="editor-page">
    <el-page-header @back="router.back()" style="margin-bottom: 20px">
      <template #content>
        <span>{{ isEdit ? '编辑流水线' : '新建流水线' }}</span>
      </template>
    </el-page-header>

    <el-row :gutter="20">
      <!-- 节点面板 -->
      <el-col :span="6">
        <el-card header="节点面板" class="node-panel">
          <div
            v-for="node in nodeTypes"
            :key="node.type"
            class="node-item"
            draggable="true"
            @dragstart="onDragStart($event, node)"
          >
            <el-icon :size="20"><component :is="node.icon" /></el-icon>
            <span>{{ node.label }}</span>
          </div>
        </el-card>
      </el-col>

      <!-- 画布区域 -->
      <el-col :span="12">
        <el-card header="流水线画布" class="canvas-card">
          <div class="vue-flow-wrapper">
            <VueFlow
              :nodes="nodes"
              :edges="edges"
              @dragover="onDragOver"
              @drop="onDrop"
              @node-click="onNodeClick"
              @edge-click="onEdgeClick"
            >
              <Background />
              <Controls />
              <MiniMap />
            </VueFlow>
          </div>
        </el-card>
      </el-col>

      <!-- 配置面板 -->
      <el-col :span="6">
        <el-card header="节点配置" class="config-panel">
          <el-form label-width="80px" v-if="selectedNode">
            <el-form-item label="节点名称">
              <el-input v-model="selectedNode.label" />
            </el-form-item>
            <el-form-item label="节点类型">
              <el-select v-model="selectedNode.type" style="width: 100%">
                <el-option label="拉取代码" value="checkout" />
                <el-option label="构建" value="build" />
                <el-option label="测试" value="test" />
                <el-option label="部署" value="deploy" />
              </el-select>
            </el-form-item>
            <el-form-item label="命令">
              <el-input v-model="selectedNode.data.command" type="textarea" :rows="4" />
            </el-form-item>
            <el-form-item label="镜像">
              <el-input v-model="selectedNode.data.image" placeholder="如：node:20-alpine" />
            </el-form-item>
            <el-form-item>
              <el-button type="danger" @click="deleteNode" style="width: 100%">删除节点</el-button>
            </el-form-item>
          </el-form>
          <el-empty v-else description="请点击画布中的节点进行配置" />
        </el-card>
      </el-col>
    </el-row>

    <!-- YAML 预览 -->
    <el-card style="margin-top: 20px">
      <template #header>
        <span>YAML 配置预览</span>
        <el-button style="float: right" @click="copyYaml">复制</el-button>
      </template>
      <el-input type="textarea" :model-value="yamlPreview" :rows="10" readonly />
    </el-card>

    <div style="margin-top: 20px; text-align: center">
      <el-button @click="router.back()">取消</el-button>
      <el-button type="primary" @click="savePipeline">保存流水线</el-button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { Plus, SetUp, Check, Promotion } from '@element-plus/icons-vue'

// 导入 Vue Flow
import { VueFlow, useVueFlow } from '@vue-flow/core'
import { Background } from '@vue-flow/background'
import { Controls } from '@vue-flow/controls'
import { MiniMap } from '@vue-flow/minimap'
import '@vue-flow/core/dist/style.css'
import '@vue-flow/controls/dist/style.css'
import '@vue-flow/minimap/dist/style.css'

import { getPipeline, createPipeline, updatePipeline } from '@/api/pipeline'

const route = useRoute()
const router = useRouter()
const isEdit = ref(!!route.params.id)
const selectedNode = ref<any>(null)

// 节点类型定义
const nodeTypes = ref([
  { type: 'checkout', label: '拉取代码', icon: 'SetUp' },
  { type: 'build', label: '构建', icon: 'SetUp' },
  { type: 'test', label: '测试', icon: 'Check' },
  { type: 'deploy', label: '部署', icon: 'Promotion' },
])

// Vue Flow 状态
const { nodes, edges, addNodes, addEdges, removeNodes, removeEdges, setNodes, setEdges } = useVueFlow()

// 初始化示例节点（如果是新建）
const initNodes = () => {
  setNodes([
    {
      id: '1',
      type: 'default',
      position: { x: 100, y: 100 },
      label: '拉取代码',
      data: {
        type: 'checkout',
        command: 'git clone $REPO_URL .',
        image: 'alpine/git',
      },
    },
    {
      id: '2',
      type: 'default',
      position: { x: 100, y: 250 },
      label: '构建项目',
      data: {
        type: 'build',
        command: 'npm run build',
        image: 'node:20-alpine',
      },
    },
  ])

  setEdges([
    { id: 'e1-2', source: '1', target: '2' },
  ])
}

// 拖拽开始
const onDragStart = (event: DragEvent, node: any) => {
  event.dataTransfer?.setData('nodeType', JSON.stringify(node))
}

// 拖拽悬停
const onDragOver = (event: DragEvent) => {
  event.preventDefault()
}

// 放下节点
const onDrop = (event: DragEvent) => {
  event.preventDefault()

  const nodeTypeStr = event.dataTransfer?.getData('nodeType')
  if (!nodeTypeStr) return

  const nodeType = JSON.parse(nodeTypeStr)
  const position = {
    x: event.offsetX - 80,
    y: event.offsetY - 40,
  }

  const newNode = {
    id: String(Date.now()),
    type: 'default',
    position,
    label: nodeType.label,
    data: {
      type: nodeType.type,
      command: '',
      image: '',
    },
  }

  addNodes([newNode])
}

// 点击节点
const onNodeClick = (event: any) => {
  const node = nodes.value.find((n: any) => n.id === event.node.id)
  selectedNode.value = node
}

// 点击边
const onEdgeClick = () => {
  selectedNode.value = null
}

// 删除节点
const deleteNode = () => {
  if (!selectedNode.value) return
  removeNodes([selectedNode.value.id])
  selectedNode.value = null
}

// YAML 预览
const yamlPreview = computed(() => {
  let yaml = 'stages:\n'
  for (const node of nodes.value) {
    yaml += `  - name: ${node.label}\n`
    yaml += `    type: ${node.data.type}\n`
    yaml += `    image: ${node.data.image || 'alpine'}\n`
    yaml += `    script:\n      - ${node.data.command || '# 请配置命令'}\n`
  }
  return yaml
})

// 保存流水线
const savePipeline = async () => {
  try {
    // 将节点和边序列化为 JSON
    const config = {
      nodes: nodes.value,
      edges: edges.value,
      yaml: yamlPreview.value,
    }

    const data = {
      name: `流水线-${Date.now()}`,
      project_id: Number(route.query.project_id) || 1,
      config_yaml: JSON.stringify(config),
    }

    if (isEdit.value) {
      await updatePipeline(String(route.params.id), data)
      ElMessage.success('更新成功')
    } else {
      await createPipeline(data)
      ElMessage.success('创建成功')
    }
    router.push('/pipeline')
  } catch (err: any) {
    ElMessage.error(err.message || '保存失败')
  }
}

// 复制 YAML
const copyYaml = () => {
  navigator.clipboard.writeText(yamlPreview.value)
  ElMessage.success('已复制到剪贴板')
}

onMounted(async () => {
  if (isEdit.value) {
    try {
      const res: any = await getPipeline(String(route.params.id))
      // 解析 config_yaml 为节点
      if (res.config_yaml) {
        try {
          const config = JSON.parse(res.config_yaml)
          if (config.nodes) setNodes(config.nodes)
          if (config.edges) setEdges(config.edges)
        } catch {
          // 如果不是 JSON，可能是旧版 YAML，忽略
        }
      }
    } catch {}
  } else {
    // 初始化示例节点
    initNodes()
  }
})
</script>

<style scoped>
.editor-page { padding: 0; }
.node-panel { min-height: 500px; }
.node-item {
  display: flex; align-items: center; gap: 8px;
  padding: 10px; margin-bottom: 8px;
  border: 1px solid #dcdfe6; border-radius: 8px;
  cursor: grab; transition: all 0.2s;
}
.node-item:hover { background: #ecf5ff; border-color: #409EFF; }
.canvas-card { min-height: 500px; }
.vue-flow-wrapper {
  height: 450px;
  border: 1px solid #dcdfe6;
  border-radius: 8px;
}
.config-panel { min-height: 500px; }
</style>
