<template>
  <div class="projects-page">
    <el-card class="page-header">
      <div class="header-actions">
        <h2>项目管理</h2>
        <el-button type="primary" @click="showCreate = true">
          <el-icon><Plus /></el-icon>
          新建项目
        </el-button>
      </div>
    </el-card>

    <el-card class="table-card">
      <el-table :data="projectList" stripe style="width: 100%">
        <el-table-column prop="name" label="项目名称" />
        <el-table-column prop="display_name" label="显示名称" />
        <el-table-column prop="repo_type" label="仓库类型" width="120" />
        <el-table-column prop="created_at" label="创建时间" width="180" />
        <el-table-column label="操作" width="200" fixed="right">
          <template #default="{ row }">
            <el-button link type="primary" size="small" @click="viewProject(row)">查看</el-button>
            <el-button link type="primary" size="small" @click="editProject(row)">编辑</el-button>
            <el-button link type="danger" size="small" @click="deleteProject(row)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <!-- 新建/编辑对话框 -->
    <el-dialog v-model="showCreate" :title="isEdit ? '编辑项目' : '新建项目'" width="520px">
      <el-form :model="form" label-width="100px">
        <el-form-item label="项目名称">
          <el-input v-model="form.name" placeholder="请输入项目名称" />
        </el-form-item>
        <el-form-item label="显示名称">
          <el-input v-model="form.display_name" placeholder="请输入显示名称" />
        </el-form-item>
        <el-form-item label="仓库类型">
          <el-select v-model="form.repo_type" style="width: 100%">
            <el-option label="Gitea" value="gitea" />
            <el-option label="GitHub" value="github" />
            <el-option label="GitLab" value="gitlab" />
            <el-option label="Gitee" value="gitee" />
          </el-select>
        </el-form-item>
        <el-form-item label="仓库地址">
          <el-input v-model="form.repo_url" placeholder="请输入仓库地址" />
        </el-form-item>
        <el-form-item label="项目描述">
          <el-input v-model="form.description" type="textarea" :rows="3" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showCreate = false">取消</el-button>
        <el-button type="primary" @click="saveProject">保存</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from "vue";
import { ElMessage, ElMessageBox } from "element-plus";
import { getProjects, createProject, updateProject, deleteProject as apiDelete } from "@/api/project";

const projectList = ref([
  { id: 1, name: "devops-star", display_name: "DevOpsStar 平台", repo_type: "gitea", created_at: "2026-05-27 00:00:00" },
  { id: 2, name: "52cv-top", display_name: "52CV 博客", repo_type: "gitea", created_at: "2026-05-20 10:30:00" },
]);
const showCreate = ref(false);
const isEdit = ref(false);
const form = ref({
  id: 0,
  name: "",
  display_name: "",
  repo_type: "gitea",
  repo_url: "",
  description: "",
});

const loadProjects = async () => {
  try {
    const res = await getProjects();
    projectList.value = res.data || projectList.value;
  } catch {}
};

const saveProject = async () => {
  try {
    if (isEdit.value) {
      await updateProject(String(form.value.id), form.value);
      ElMessage.success("更新成功");
    } else {
      await createProject(form.value);
      ElMessage.success("创建成功");
    }
    showCreate.value = false;
    loadProjects();
  } catch {}
};

const viewProject = (row: any) => {
  // 跳转到项目详情
};

const editProject = (row: any) => {
  isEdit.value = true;
  form.value = { ...row };
  showCreate.value = true;
};

const deleteProject = async (row: any) => {
  await ElMessageBox.confirm("确定删除该项目？", "提示", { type: "warning" });
  try {
    await apiDelete(String(row.id));
    ElMessage.success("删除成功");
    loadProjects();
  } catch {}
};

onMounted(() => {
  loadProjects();
});
</script>

<style scoped>
.projects-page {
  display: flex;
  flex-direction: column;
  gap: 20px;
}
.page-header {
  padding: 20px;
}
.header-actions {
  display: flex;
  justify-content: space-between;
  align-items: center;
}
.header-actions h2 {
  margin: 0;
}
.table-card {
  padding: 20px;
}
</style>
