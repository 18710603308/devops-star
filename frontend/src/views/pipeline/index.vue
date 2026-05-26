<template>
  <div class="pipeline-page">
    <el-card class="page-header">
      <div class="header-actions">
        <h2>流水线管理</h2>
        <el-button type="primary" @click="showCreate = true">
          <el-icon><Plus /></el-icon>
          新建流水线
        </el-button>
      </div>
    </el-card>

    <el-card class="table-card">
      <el-table :data="pipelineList" stripe style="width: 100%">
        <el-table-column prop="name" label="流水线名称" />
        <el-table-column prop="project_name" label="所属项目" width="160" />
        <el-table-column prop="status" label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="statusType(row.status)">{{ statusText(row.status) }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="last_run_id" label="最后运行" width="160" />
        <el-table-column prop="updated_at" label="更新时间" width="180" />
        <el-table-column label="操作" width="240" fixed="right">
          <template #default="{ row }">
            <el-button link type="primary" size="small" @click="editPipeline(row)">编辑</el-button>
            <el-button link type="primary" size="small" @click="triggerPipeline(row)">触发</el-button>
            <el-button link type="primary" size="small" @click="viewLogs(row)">日志</el-button>
            <el-button link type="danger" size="small" @click="deletePipeline(row)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <!-- 新建/编辑对话框 -->
    <el-dialog v-model="showCreate" :title="isEdit ? '编辑流水线' : '新建流水线'" width="600px">
      <el-form :model="form" label-width="100px">
        <el-form-item label="流水线名称">
          <el-input v-model="form.name" placeholder="请输入流水线名称" />
        </el-form-item>
        <el-form-item label="所属项目">
          <el-select v-model="form.project_id" style="width: 100%">
            <el-option v-for="p in projectList" :key="p.id" :label="p.name" :value="p.id" />
          </el-select>
        </el-form-item>
        <el-form-item label="描述">
          <el-input v-model="form.description" type="textarea" :rows="3" />
        </el-form-item>
        <el-form-item label="YAML 配置">
          <el-input v-model="form.config_yaml" type="textarea" :rows="6" placeholder="请输入流水线 YAML 配置" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showCreate = false">取消</el-button>
        <el-button type="primary" @click="savePipeline">保存</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from "vue";
import { ElMessage, ElMessageBox } from "element-plus";
import { Plus } from "@element-plus/icons-vue";
import { getPipelines, createPipeline, updatePipeline, triggerPipeline as apiTrigger, deletePipeline as apiDelete } from "@/api/pipeline";
import { getProjects } from "@/api/project";

const pipelineList = ref([
  { id: 1, name: "前端构建部署", project_name: "devops-star", status: "success", last_run_id: "run-001", updated_at: "2026-05-27 00:30:00" },
  { id: 2, name: "后端镜像构建", project_name: "devops-star", status: "running", last_run_id: "run-002", updated_at: "2026-05-27 00:35:00" },
]);
const projectList = ref<any[]>([]);
const showCreate = ref(false);
const isEdit = ref(false);
const form = ref({ id: 0, name: "", project_id: 0, description: "", config_yaml: "" });

const statusType = (s: string) => ({ success: "success", running: "primary", failed: "danger", idle: "info" }[s] || "info");
const statusText = (s: string) => ({ success: "成功", running: "运行中", failed: "失败", idle: "空闲" }[s] || s);

const loadData = async () => {
  try {
    const res = await getPipelines();
    pipelineList.value = res.data || pipelineList.value;
    const res2 = await getProjects();
    projectList.value = res2.data || [];
  } catch {}
};

const savePipeline = async () => {
  try {
    if (isEdit.value) {
      await updatePipeline(String(form.value.id), form.value);
      ElMessage.success("更新成功");
    } else {
      await createPipeline(form.value);
      ElMessage.success("创建成功");
    }
    showCreate.value = false;
    loadData();
  } catch {}
};

const editPipeline = (row: any) => {
  isEdit.value = true;
  form.value = { ...row };
  showCreate.value = true;
};

const triggerPipeline = async (row: any) => {
  try {
    await apiTrigger(String(row.id));
    ElMessage.success("流水线已触发");
    loadData();
  } catch {}
};

const viewLogs = (row: any) => {
  // 打开日志抽屉
};

const deletePipeline = async (row: any) => {
  await ElMessageBox.confirm("确定删除该流水线？", "提示", { type: "warning" });
  try {
    await apiDelete(String(row.id));
    ElMessage.success("删除成功");
    loadData();
  } catch {}
};

onMounted(loadData);
</script>

<style scoped>
.pipeline-page { display: flex; flex-direction: column; gap: 20px; }
.page-header { padding: 20px; }
.header-actions { display: flex; justify-content: space-between; align-items: center; }
.header-actions h2 { margin: 0; }
.table-card { padding: 20px; }
</style>
